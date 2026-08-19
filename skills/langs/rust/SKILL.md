---
name: rust
description: Rust 语言实践。按已安装/锁定的 rustc 1.y（优先 1.95 / 1.96 / 1.97）选用该版本 API。在编辑 .rs、Cargo.toml / Cargo.lock、rust-toolchain.toml，或编写 Rust 测试与宏时使用。
---

# Rust

禁止凭训练数据写 API。先定 **rustc 1.y**，再读对应 reference。Axum / Actix / Rocket 等 web 框架不在本技能里指定。

`edition` 是第二轴：只用来避开 edition 门控语法，**不能**代替 rustc 档。

## 定版本

1. 读已锁定工具链，不要只看 `Cargo.toml` 的 `rust-version`（那是 MSRV 下限）。
   - `rust-toolchain.toml` 或 `rust-toolchain` 的 `channel`（两文件都在时，按 rustup：**无后缀的 `rust-toolchain` 优先**）。
   - `channel` 是 `stable` / `beta` / `nightly` 时，再跑 `rustc --version` 取实际 `1.y`。
2. 没有工具链文件：用当前 `rustc --version`。
3. 仍没有：才用 `Cargo.toml` 的 `package.rust-version`（按该下限写，不要用更高档 API）。

取 **1.y**（`1.97.1` → `1.97`），打开下表，按该文件写代码。

| 1.y | 文件 |
| --- | --- |
| 1.95 | [references/1.95.md](references/1.95.md) |
| 1.96 | [references/1.96.md](references/1.96.md) |
| 1.97 | [references/1.97.md](references/1.97.md) |

未列出的更新 1.y：以已覆盖的最高档为底，再查 [RELEASES.md](https://github.com/rust-lang/rust/blob/master/RELEASES.md) / [blog.rust-lang.org](https://blog.rust-lang.org/)。更旧：不要使用本技能里的新 API。

写代码只开当前档。1.97 含仍该用的 1.95/1.96 索引；需要例子再打开旧档。

## edition

读 `Cargo.toml` 的 `edition`（workspace 可能在 `[workspace.package]`）。rustc 新、edition 仍是 `2021` 时，不要写 2024 才启用的语法。需要某 edition 特性时先改 `edition`，不要靠 `#![feature]`。

## 硬规则（1.95+ 都成立）

- **let chains**：`if let Some(x) = a && let Ok(y) = b { }`。不要嵌套 `if let`。
- 错误用 `?` / `let else` / `if let`。不要把 `unwrap()` / `expect()` 当控制流（测试里可以）。
- 优先 std 与 `Cargo.toml` 已有 crate。不要为已有功能再拉重复实现。
- 模块、错误类型、测试框架跟项目走。不要在本技能里安利 web 框架或新异步运行时。
- `unsafe` 保持最小，紧邻写清 SAFETY。
