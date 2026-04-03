# 地图瓦片管理系统

一个功能完整的地图瓦片服务管理系统，支持 XYZ 栅格瓦片和 3D Tiles 三维瓦片的统一管理、授权访问和数据统计。

## 核心功能

### 1. 数据源管理
- 支持 XYZ 栅格瓦片（PNG/JPG/WebP）
- 支持 3D Tiles 三维瓦片
- 可视化地图预览（Leaflet 2D / Cesium 3D）
- 自定义地图中心点和默认缩放级别
- 数据源启用/禁用控制

### 2. 授权管理
- API Key 生成和管理
- 数据源权限精细控制（可绑定多个数据源）
- IP 白名单限制
- Referer 白名单限制
- 调用次数统计
- 自动生成对应格式的访问 URL（XYZ/3D Tiles）

### 3. 访问日志
- 完整的瓦片请求日志记录（含 XYZ 和 3D Tiles）
- 请求类型标识（XYZ / 3D Tiles）
- 按 API Key、数据源、IP、时间筛选
- 显示响应时间和状态码
- 分页浏览，支持自定义每页条数
- 180 天明细数据保留，永久聚合统计

### 4. 仪表盘
- 今日 / 昨日 / 本月调用统计
- 全部累计调用次数
- 近 7 天请求趋势图表
- Top 5 API Key 排行
- 实时数据可视化

### 5. 安全特性
- 登录密码 SHA256 加密传输
- 连续登录失败 5 次锁定 30 分钟
- JWT Token 认证
- Redis 会话管理
- 预览功能独立认证（不计入统计）

## 项目结构

```
02MapTileSystem/
├── map-tile-system/          # 前端项目（Vue3 + TypeScript）
│   ├── src/
│   │   ├── api/              # API 接口
│   │   ├── assets/           # 静态资源
│   │   ├── components/       # 组件（MapPreview2D / MapPreview3D）
│   │   ├── router/           # 路由
│   │   ├── stores/           # 状态管理
│   │   ├── types/            # 类型定义
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   │       ├── dashboard/    # 仪表盘
│   │       ├── datasource/   # 数据源管理
│   │       ├── auth/         # 授权管理
│   │       ├── access-log/   # 访问日志
│   │       └── statistics/   # 统计报表
│   ├── vite.config.ts
│   └── package.json
│
├── backend/                  # 后端项目（Go + Gin）
│   ├── cmd/
│   │   └── server/           # 主程序入口
│   ├── internal/
│   │   ├── handler/          # HTTP 处理器
│   │   ├── middleware/       # 中间件（JWT / API Key 鉴权）
│   │   └── model/            # 数据模型
│   ├── pkg/
│   │   ├── config/           # 配置管理
│   │   ├── database/         # 数据库初始化
│   │   ├── logger/           # 日志
│   │   ├── scheduler/        # 定时任务
│   │   └── utils/            # 工具
│   ├── configs/              # 配置文件
│   └── go.mod
│
├── PRD.md                    # 产品需求文档
└── README.md                 # 项目说明
```

## 技术栈

### 前端
- Vue 3 + TypeScript + Composition API
- Element Plus（UI 组件库）
- Tailwind CSS（样式框架）
- ECharts（数据可视化）
- Leaflet（2D 地图预览）
- Cesium（3D 地图预览）
- vite-plugin-cesium（Cesium Vite 集成）
- Axios（HTTP 客户端）
- Crypto-js（密码加密）

### 后端
- Go 1.21+
- Gin（Web 框架）
- GORM（ORM）
- PostgreSQL（数据库）
- Redis（缓存 / 会话 / 登录限流）
- 定时任务调度器（每日日志汇总与清理）

## 快速开始

### 环境要求
- Node.js 20+
- Go 1.21+
- PostgreSQL 12+
- Redis 6+

### 前端

```bash
cd map-tile-system

npm install

npm run dev
# 访问 http://localhost:5173
```

### 后端

```bash
cd backend

go mod download

# 修改 configs/config.yaml 配置数据库和 Redis 连接

go run cmd/server/main.go
# 服务运行在 http://localhost:8080
```

### 默认账户
- 用户名：`admin`
- 密码：`admin123`

## API 接口

### 认证
- `POST /api/auth/login` — 登录

### 数据源
- `GET    /api/datasources` — 获取列表
- `POST   /api/datasources` — 创建
- `PUT    /api/datasources/:id` — 更新
- `DELETE /api/datasources/:id` — 删除

### API Key
- `GET    /api/apikeys` — 获取列表
- `POST   /api/apikeys` — 创建
- `PUT    /api/apikeys/:id` — 更新
- `DELETE /api/apikeys/:id` — 删除

### 统计
- `GET /api/statistics/overview` — 概览
- `GET /api/statistics/trend` — 近 7 天趋势
- `GET /api/statistics/details` — 访问日志明细（分页）
- `GET /api/statistics/top-keys` — Top 5 API Key

### 瓦片服务（需携带 `?key=` 参数）
- `GET /tiles/:dataset/:z/:x/:y` — XYZ 瓦片
- `GET /tiles/:dataset/tileset.json` — 3D Tiles 入口
- `GET /tiles/:dataset/3dtiles/*filepath` — 3D Tiles 子文件

### 内部预览（JWT 认证，不计入统计）
- `GET /api/preview/xyz/:dataset/:z/:x/:y` — XYZ 预览
- `GET /api/preview/3dtiles/:dataset/*filepath` — 3D Tiles 预览

## 部署

### 生产部署

1. 前端构建并部署到 Nginx：
   ```bash
   cd map-tile-system && npm run build
   ```
2. 后端编译：
   ```bash
   cd backend && go build -o map-tile-system cmd/server/main.go
   ```
3. 配置 PostgreSQL 和 Redis
4. 配置 Nginx 反向代理（将 `/api` 和 `/tiles` 代理到后端 8080 端口）

## 许可证

MIT License
