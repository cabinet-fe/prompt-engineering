# HTML 嵌套与样式膨胀

后台/中台优先 UI 库与布局组件；结构宜浅，CSS 宜少。

## 反例：用嵌套堆间距

```vue
<div class="page">
  <div class="page-inner">
    <div class="page-body">
      <div class="card-wrap">
        <div class="card">
          <div class="card-body">
            <el-form>...</el-form>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
```

```scss
.page { padding: 16px; }
.page-inner { padding: 8px; }
.page-body { margin: 8px; }
.card-wrap { background: #fff; }
.card { border-radius: 4px; }
.card-body { padding: 16px; }
```

## 正例：布局组件 + 表单

```vue
<page-container>
  <el-card>
    <el-form>...</el-form>
  </el-card>
</page-container>
```

（`page-container` / `el-card` 以项目实际布局与 UI 库为准。）

## 反例：重复、无用 CSS

```scss
.btn-primary { color: #fff; background: #409eff; }
// 与 UI 库 Primary Button 重复
.toolbar .btn-primary { color: #fff; background: #409eff; }
```

模板已改用 `<el-button type="primary">`，自定义 `.btn-primary` 仍留着。

## 正例

用 `<el-button type="primary">`；删除无用自定义按钮样式。

## 反例：样式累加覆盖（`!important` 补丁）

为了调整外观，不断在文件末尾追加新规则或 `!important` 覆盖旧样式，留下一堆互相打架、实际无效的死样式：

```scss
.card-title { font-size: 16px; color: #333; }
// 后续调整又在下面追加覆盖
.card-title { font-size: 14px !important; color: #111 !important; }
```

## 正例：直接替换并收敛

直接修改或替换原样式规则，同步删掉被覆盖且不再起作用的选择器和属性，不留 `!important` 补丁。

## 反例：为对齐多包三层 flex

```vue
<div class="flex">
  <div class="flex-1">
    <div class="flex items-center justify-between">
      <div class="flex items-center gap-2">...</div>
    </div>
  </div>
</div>
```

## 正例

```vue
<el-space wrap>
  <el-button>导出</el-button>
  <el-button type="primary">新建</el-button>
</el-space>
```

或项目 `a-space` / `ProTable` toolbar 插槽——优先设计系统已有能力。

## 自检

- 去掉某一层 wrapper 后布局是否仍成立？成立则删。
- 这段 CSS 是否只是重复 UI 库 token/组件皮肤？是则删。
- 是否靠后写规则或 `!important` 覆盖了前面的旧样式？是则合并回原规则并删掉被覆盖的死属性。
