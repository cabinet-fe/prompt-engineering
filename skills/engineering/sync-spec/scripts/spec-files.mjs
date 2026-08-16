#!/usr/bin/env node
// 规格-文件索引工具：维护和查询「规格 -> 影响文件/glob」。
// 用法：
//   node spec-files.mjs init <index-file>
//   node spec-files.mjs set <index-file> <spec-path> --module <name> --files <pattern...>
//   node spec-files.mjs remove <index-file> <spec-path>
//   node spec-files.mjs query <index-file> [--json] [--stdin] <file...>
//   node spec-files.mjs list <index-file>

import fs from 'node:fs';
import path from 'node:path';

const VERSION = 1;

function usage() {
  process.stderr.write(`usage:
  node spec-files.mjs init <index-file>
  node spec-files.mjs set <index-file> <spec-path> --module <name> --files <pattern...>
  node spec-files.mjs remove <index-file> <spec-path>
  node spec-files.mjs query <index-file> [--json] [--stdin] <file...>
  node spec-files.mjs list <index-file>
`);
}

function resolveIndex(file) {
  return path.resolve(process.cwd(), file);
}

function ensureIndex(file) {
  const abs = resolveIndex(file);
  if (fs.existsSync(abs)) return { abs, data: readIndex(abs) };
  fs.mkdirSync(path.dirname(abs), { recursive: true });
  const data = { version: VERSION, specs: {} };
  fs.writeFileSync(abs, `${JSON.stringify(data, null, 2)}\n`);
  return { abs, data };
}

function readIndex(file) {
  const abs = resolveIndex(file);
  if (!fs.existsSync(abs)) return { version: VERSION, specs: {} };
  const raw = fs.readFileSync(abs, 'utf8');
  try {
    const data = JSON.parse(raw);
    return {
      version: data.version ?? VERSION,
      specs: data.specs && typeof data.specs === 'object' ? data.specs : {},
    };
  } catch (err) {
    throw new Error(`无法解析索引 ${abs}: ${err.message}`);
  }
}

function writeIndex(file, data) {
  const abs = resolveIndex(file);
  fs.mkdirSync(path.dirname(abs), { recursive: true });
  fs.writeFileSync(abs, `${JSON.stringify({ version: VERSION, specs: data.specs }, null, 2)}\n`);
  return abs;
}

function normalizeRel(file) {
  let rel = path.relative(process.cwd(), path.resolve(process.cwd(), file));
  if (rel === '') return '';
  return rel.split(path.sep).join('/').replace(/^\.\//, '');
}

function normalizeSpec(spec) {
  return spec.split(path.sep).join('/').replace(/^\.\//, '').replace(/^\/+/, '');
}

function escapeRegex(ch) {
  return /[\\^$.*+?()[\]{}|]/.test(ch) ? `\\${ch}` : ch;
}

function globToRegex(pattern) {
  let p = String(pattern).split(path.sep).join('/').replace(/^\.\//, '').replace(/\/+$/, '');
  const hasWildcard = /[*?]/.test(p);
  let out = '^';
  for (let i = 0; i < p.length; i += 1) {
    const c = p[i];
    if (c === '*' && p[i + 1] === '*') {
      if (p[i + 2] === '/') {
        out += '(?:.*/)?';
        i += 2;
      } else {
        out += '.*';
        i += 1;
      }
    } else if (c === '*') {
      out += '[^/]*';
    } else if (c === '?') {
      out += '[^/]';
    } else {
      out += escapeRegex(c);
    }
  }
  out += '$';
  const re = new RegExp(out);
  if (hasWildcard) return re;
  // 无通配符时同时支持精确文件和目录前缀。
  const prefix = new RegExp(`^${p.replace(/[\\^$.*+?()[\]{}|]/g, '\\$&')}/`);
  return { test: (file) => re.test(file) || prefix.test(file) };
}

function matchesPattern(file, pattern) {
  const matcher = globToRegex(pattern);
  return matcher.test(file);
}

function matchesInput(file, pattern) {
  // 输入可能是文件，也可能是目录/模块路径；目录输入允许命中 `dir/**` 这类 glob。
  return matchesPattern(file, pattern) || matchesPattern(`${file}/__spec_probe__`, pattern);
}

function readSpecTitle(indexAbs, specRel) {
  try {
    const specAbs = path.resolve(path.dirname(indexAbs), specRel);
    const head = fs.readFileSync(specAbs, 'utf8').slice(0, 1000);
    const match = head.match(/^#\s+(.+)$/m);
    return match ? match[1].trim() : '';
  } catch {
    return '';
  }
}

function queryIndex(file, files, { json = false } = {}) {
  const indexAbs = resolveIndex(file);
  const data = readIndex(indexAbs);
  const normalized = [...new Set(files.map(normalizeRel).filter(Boolean))];
  const matches = [];
  for (const [spec, entry] of Object.entries(data.specs ?? {})) {
    const patterns = Array.isArray(entry?.files) ? entry.files : [];
    const hit = [];
    for (const file of normalized) {
      for (const pattern of patterns) {
        if (matchesInput(file, pattern)) {
          hit.push({ file, pattern });
        }
      }
    }
    if (hit.length > 0) {
      matches.push({
        spec,
        module: entry?.module ?? '',
        title: readSpecTitle(indexAbs, spec),
        hit,
      });
    }
  }
  if (json) {
    process.stdout.write(`${JSON.stringify({ files: normalized, matches }, null, 2)}\n`);
    return;
  }
  if (matches.length === 0) {
    process.stdout.write('NO_MATCH\n');
    return;
  }
  for (const match of matches) {
    process.stdout.write(`${match.spec}${match.title ? ` :: ${match.title}` : ''}\n`);
    for (const h of match.hit) {
      process.stdout.write(`  ${h.file} <= ${h.pattern}\n`);
    }
  }
}

function setEntry(file, spec, module, patterns) {
  const { data } = ensureIndex(file);
  const key = normalizeSpec(spec);
  data.specs[key] = {
    module: module || data.specs[key]?.module || '',
    files: patterns.map(normalizeRel),
  };
  return writeIndex(file, data);
}

function removeEntry(file, spec) {
  const { data } = ensureIndex(file);
  const key = normalizeSpec(spec);
  delete data.specs[key];
  return writeIndex(file, data);
}

function listIndex(file) {
  const data = readIndex(file);
  const entries = Object.entries(data.specs ?? {});
  if (entries.length === 0) {
    process.stdout.write('EMPTY\n');
    return;
  }
  for (const [spec, entry] of entries.sort(([a], [b]) => a.localeCompare(b))) {
    const files = Array.isArray(entry?.files) ? entry.files.join(', ') : '';
    process.stdout.write(`${spec}\t${entry?.module ?? ''}\t${files}\n`);
  }
}

function main() {
  const args = process.argv.slice(2);
  const [cmd, indexFile, ...rest] = args;
  if (!cmd || !indexFile) {
    usage();
    process.exit(2);
  }

  try {
    if (cmd === 'init') {
      const { abs } = ensureIndex(indexFile);
      process.stdout.write(`${abs}\n`);
      return;
    }

    if (cmd === 'list') {
      listIndex(indexFile);
      return;
    }

    if (cmd === 'remove') {
      const spec = rest[0];
      if (!spec) throw new Error('remove 需要 spec 路径');
      const abs = removeEntry(indexFile, spec);
      process.stdout.write(`${abs}\n`);
      return;
    }

    if (cmd === 'set') {
      const spec = rest[0];
      if (!spec) throw new Error('set 需要 spec 路径');
      let module = '';
      const patterns = [];
      let i = 1;
      while (i < rest.length) {
        if (rest[i] === '--module') {
          module = rest[i + 1] ?? '';
          i += 2;
        } else if (rest[i] === '--files') {
          patterns.push(...rest.slice(i + 1));
          i = rest.length;
        } else {
          i += 1;
        }
      }
      if (patterns.length === 0) throw new Error('set 需要 --files 和至少一个路径/glob');
      const abs = setEntry(indexFile, spec, module, patterns);
      process.stdout.write(`${abs}\n`);
      return;
    }

    if (cmd === 'query') {
      const json = rest.includes('--json');
      const stdin = rest.includes('--stdin');
      const files = rest.filter((x) => x !== '--json' && x !== '--stdin');
      if (stdin) {
        const input = fs.readFileSync(0, 'utf8');
        files.push(...input.split(/\r?\n/).map((x) => x.trim()).filter(Boolean));
      }
      if (files.length === 0) throw new Error('query 需要至少一个文件路径，或用 --stdin 传入');
      queryIndex(indexFile, files, { json });
      return;
    }

    usage();
    process.exit(2);
  } catch (err) {
    process.stderr.write(`ERROR: ${err.message}\n`);
    process.exit(1);
  }
}

main();
