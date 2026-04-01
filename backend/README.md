# 地图瓦片服务管理系统 - 后端

## 技术栈
- Go 1.21+
- Gin (Web 框架)
- GORM (ORM)
- PostgreSQL (数据库)
- Redis (缓存)
- Viper (配置管理)
- Zap (日志)

## 项目结构
```
backend/
├── cmd/
│   └── server/
│       └── main.go           # 主入口
├── internal/
│   ├── handler/              # HTTP 处理器
│   │   ├── handler.go        # 基础处理器
│   │   ├── datasource.go     # 数据源管理
│   │   ├── apikey.go         # API Key 管理
│   │   ├── statistics.go     # 统计数据
│   │   └── tile.go           # 瓦片服务
│   ├── middleware/           # 中间件
│   │   ├── auth.go           # 认证中间件
│   │   └── logger.go         # 日志中间件
│   ├── model/                # 数据模型
│   │   └── model.go
│   ├── service/              # 业务逻辑层
│   └── repository/           # 数据访问层
├── pkg/
│   ├── config/               # 配置管理
│   │   └── config.go
│   ├── database/             # 数据库连接
│   │   ├── postgres.go
│   │   └── redis.go
│   ├── logger/               # 日志工具
│   │   └── logger.go
│   └── utils/                # 工具函数
├── configs/
│   └── config.yaml           # 配置文件
├── data/
│   └── tiles/                # 瓦片数据存储目录
├── go.mod
└── go.sum
```

## 快速开始

### 1. 安装依赖
```bash
cd backend
go mod download
```

### 2. 配置数据库
创建 PostgreSQL 数据库：
```sql
CREATE DATABASE map_tile_system;
```

### 3. 修改配置
编辑 `configs/config.yaml`，配置数据库连接信息。

### 4. 运行服务
```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

## API 接口

### 认证接口
- `POST /api/auth/login` - 用户登录

### 数据源管理
- `GET /api/datasources` - 获取数据源列表
- `POST /api/datasources` - 创建数据源
- `GET /api/datasources/:id` - 获取单个数据源
- `PUT /api/datasources/:id` - 更新数据源
- `DELETE /api/datasources/:id` - 删除数据源

### API Key 管理
- `GET /api/apikeys` - 获取 API Key 列表
- `POST /api/apikeys` - 创建 API Key
- `GET /api/apikeys/:id` - 获取单个 API Key
- `PUT /api/apikeys/:id` - 更新 API Key
- `DELETE /api/apikeys/:id` - 删除 API Key

### 统计数据
- `GET /api/statistics/overview` - 获取统计概览
- `GET /api/statistics/trend` - 获取统计趋势
- `GET /api/statistics/details` - 获取统计明细

### 瓦片服务
- `GET /tiles/:dataset/:z/:x/:y` - XYZ 瓦片服务
- `GET /tiles/:dataset/tileset.json` - 3D Tiles 入口
- `GET /tiles/:dataset/*filepath` - 3D Tiles 其他文件

## 开发说明

### 数据库迁移
应用启动时会自动执行数据库迁移，创建所需的表结构。

### 日志
日志使用 Zap 库，支持结构化日志输出。

### 缓存
使用 Redis 缓存 API Key 权限信息和热点瓦片数据。

## 部署

### 使用 Docker
```bash
# 构建镜像
docker build -t map-tile-system-backend .

# 运行容器
docker run -p 8080:8080 map-tile-system-backend
```

### 直接部署
```bash
# 编译
go build -o server cmd/server/main.go

# 运行
./server
```

## 环境变量
可以通过环境变量覆盖配置文件中的设置：
- `SERVER_PORT` - 服务端口
- `DB_HOST` - 数据库主机
- `DB_PORT` - 数据库端口
- `DB_USER` - 数据库用户
- `DB_PASSWORD` - 数据库密码
- `REDIS_HOST` - Redis 主机
- `REDIS_PORT` - Redis 端口
