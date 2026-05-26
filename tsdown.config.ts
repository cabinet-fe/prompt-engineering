import { defineConfig } from 'tsdown'

export default defineConfig({
  entry: ['src/index.ts', 'src/cli/index.ts'],
  format: ['esm'],
  dts: true,
  clean: true,
  platform: 'node',
  unbundle: true,
  fixedExtension: false,
})
