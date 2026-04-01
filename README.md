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
- 数据源权限精细控制
- IP 白名单限制
- Referer 白名单限制
- 调用次数统计

### 3. 访问日志
- 完整的瓦片请求日志记录
- 按 API Key、数据源、IP、时间筛选
- 显示响应时间和状态码
- 180天明细数据保留
- 永久聚合统计数据

### 4. 仪表盘
- 今日/昨日/本月调用统计
- 全部累计调用次数
- 近7天请求趋势图表
- Top 5 API Key 排行
- 实时数据可视化

### 5. 安全特性
- 登录密码 SHA256 加密传输
- 连续登录失败5次锁定30分钟
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
│   │   ├── components/       # 组件
│   │   ├── router/           # 路由
│   │   ├── stores/           # 状态管理
│   │   ├── types/            # 类型定义
│   │   ├── utils/            # 工具函数
│   │   └── views/            # 页面视图
│   ├── package.json
│   └── README.md
│
├── backend/                  # 后端项目（Go + Gin）
│   ├── cmd/
│   │   └── server/           # 主程序入口
│   ├── internal/
│   │   ├── handler/          # HTTP 处理器
│   │   ├── middleware/       # 中间件
│   │   ├── model/            # 数据模型
│   │   ├── service/          # 业务逻辑
│   │   └── repository/       # 数据访问
│   ├── pkg/
│   │   ├── config/           # 配置管理
│   │   ├── database/         # 数据库
│   │   ├── logger/           # 日志
│   │   └── utils/            # 工具
│   ├── configs/              # 配置文件
│   ├── go.mod
│   └── README.md
│
├── PRD.md                    # 产品需求文档
├── AGENTS.md                 # 开发规范
└── README.md                 # 项目说明
```

## 技术栈

### 前端
- Vue 3 + TypeScript + Composition API
- Element Plus (UI 组件库)
- Tailwind CSS (样式框架)
- ECharts (数据可视化)
- Leaflet (2D 地图预览)
- Axios (HTTP 客户端)
- Crypto-js (密码加密)

### 后端
- Go 1.21+
- Gin (Web 框架)
- GORM (ORM)
- PostgreSQL (数据库)
- Redis (缓存/会话)
- 定时任务调度器

## 快速开始

### 环境要求
- Node.js 16+
- Go 1.21+
- PostgreSQL 12+
- Redis 6+

### 前端开发

```bash
# 进入前端目录
cd map-tile-system

# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 访问 http://localhost:5173
```

### 后端开发

```bash
# 进入后端目录
cd backend

# 安装 Go 依赖
go mod download

# 配置数据库
# 1. 创建数据库: CREATE DATABASE map_tile_system;
# 2. 创建 Redis 实例

# 修改配置文件 configs/config.yaml
# 配置数据库连接、Redis 连接等

# 启动后端服务
go run cmd/server/main.go

# 服务运行在 http://localhost:8080
```

### 默认账户
- 用户名: `admin`
- 密码: `admin123`

## 功能模块

### 1. 用户认证
- 管理员登录
- JWT Token 认证

### 2. 数据源管理
- XYZ 栅格瓦片管理
- 3D Tiles 数据源管理
- 数据源启用/禁用

### 3. 授权管理
- API Key 生成和管理
- 数据源权限绑定
- IP 白名单
- Referer 白名单

### 4. 瓦片服务
- XYZ 瓦片服务 (`/tiles/:dataset/:z/:x/:y`)
- 3D Tiles 服务 (`/tiles/:dataset/tileset.json`)
- API Key 鉴权

### 5. 统计报表
- 今日/昨日/本月调用统计
- 调用趋势图表
- Top Key 排行
- 调用明细导出

## API 接口

### 认证
- `POST /api/auth/login` - 登录

### 数据源
- `GET /api/datasources` - 获取列表
- `POST /api/datasources` - 创建
- `PUT /api/datasources/:id` - 更新
- `DELETE /api/datasources/:id` - 删除

### API Key
- `GET /api/apikeys` - 获取列表
- `POST /api/apikeys` - 创建
- `PUT /api/apikeys/:id` - 更新
- `DELETE /api/apikeys/:id` - 删除

### 统计
- `GET /api/statistics/overview` - 概览
- `GET /api/statistics/trend` - 趋势
- `GET /api/statistics/details` - 明细

### 瓦片服务
- `GET /tiles/:dataset/:z/:x/:y` - XYZ 瓦片
- `GET /tiles/:dataset/tileset.json` - 3D Tiles

## 部署

### Docker 部署

```bash
# 构建前端
cd map-tile-system
npm run build

# 构建后端
cd ../backend
docker build -t map-tile-system-backend .

# 运行
docker run -p 8080:8080 map-tile-system-backend
```

### 生产部署

1. 前端构建并部署到 Nginx
2. 后端编译并使用 systemd 管理
3. 配置 PostgreSQL 和 Redis
4. 配置 Nginx 反向代理

## 开发规范

- 使用 TypeScript 确保类型安全
- 组件使用函数式组件 + Hooks
- 使用 Tailwind CSS 编写样式
- 遵循 ESLint 和 Prettier 规范
- Go 代码遵循标准格式化规范

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request。
