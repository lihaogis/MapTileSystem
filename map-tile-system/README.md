# 地图瓦片服务管理系统

这是一个地图瓦片服务的管理系统，支持 XYZ 格式栅格瓦片和 3D Tiles 三维格式瓦片，提供授权 key 方式调用，可查看每天的调用统计。

## 技术栈
- Vue 3 + TypeScript
- Element Plus (UI 组件库)
- Tailwind CSS (样式框架)
- ECharts (数据可视化)
- Vue Router (路由管理)
- Pinia (状态管理)
- Axios (HTTP 客户端)

## 技术栈
- Vue 3 + TypeScript
- Element Plus (UI 组件库)
- Tailwind CSS (样式框架)
- ECharts (数据可视化)
- Vue Router (路由管理)
- Pinia (状态管理)
- Axios (HTTP 客户端)

## 项目结构
```
src/
├── api/              # API 接口
│   └── request.ts    # Axios 封装
├── components/       # 组件
│   ├── common/       # 通用组件
│   └── layout/       # 布局组件
│       └── MainLayout.vue
├── router/           # 路由配置
├── stores/           # 状态管理
├── types/            # TypeScript 类型定义
├── utils/            # 工具函数
└── views/            # 页面视图
    ├── auth/         # 授权管理
    ├── dashboard/    # 仪表盘
    ├── datasource/   # 数据源管理
    ├── login/        # 登录页
    └── statistics/   # 统计报表
```

## 功能模块
1. **登录认证** - 管理员登录，支持"记住我"功能
2. **仪表盘** - 展示今日/昨日/本月调用统计，近7天趋势图，Top 5 Key 排行
3. **数据源管理** - 管理 XYZ 栅格瓦片和 3D Tiles 数据源
4. **授权管理** - 管理 API Key，配置访问权限
5. **统计报表** - 查看调用趋势和明细数据，支持导出

## 访问地址
- 开发环境: http://localhost:5173
- 默认登录: 任意用户名和密码（当前为 Mock 数据）

## Recommended IDE Setup

[VS Code](https://code.visualstudio.com/) + [Vue (Official)](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).

## Recommended Browser Setup

- Chromium-based browsers (Chrome, Edge, Brave, etc.):
  - [Vue.js devtools](https://chromewebstore.google.com/detail/vuejs-devtools/nhdogjmejiglipccpnnnanhbledajbpd)
  - [Turn on Custom Object Formatter in Chrome DevTools](http://bit.ly/object-formatters)
- Firefox:
  - [Vue.js devtools](https://addons.mozilla.org/en-US/firefox/addon/vue-js-devtools/)
  - [Turn on Custom Object Formatter in Firefox DevTools](https://fxdx.dev/firefox-devtools-custom-object-formatters/)

## Type Support for `.vue` Imports in TS

TypeScript cannot handle type information for `.vue` imports by default, so we replace the `tsc` CLI with `vue-tsc` for type checking. In editors, we need [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) to make the TypeScript language service aware of `.vue` types.

## Customize configuration

See [Vite Configuration Reference](https://vite.dev/config/).

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```
