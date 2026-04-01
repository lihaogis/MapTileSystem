// 数据源类型
export interface DataSource {
  id: string
  name: string
  type: 'xyz' | '3dtiles'
  format?: string
  path: string
  minZoom?: number
  maxZoom?: number
  centerLat?: number
  centerLng?: number
  defaultZoom?: number
  status: 'enabled' | 'disabled'
  createdAt: string
  updatedAt: string
}

// API Key 类型
export interface ApiKey {
  id: string
  key: string
  name: string
  dataSources: string[]
  ipWhitelist?: string[]
  refererWhitelist?: string[]
  status: 'enabled' | 'disabled'
  lastUsedAt?: string
  createdAt: string
  callCount: number
}

// 调用日志类型
export interface CallLog {
  id: number
  apiKeyId: string
  dataSourceId: string
  apiKeyName?: string
  dataSourceName?: string
  tileZ: number
  tileX: number
  tileY: number
  statusCode: number
  responseTime: number
  ipAddress: string
  userAgent: string
  createdAt: string
}

// 统计数据类型
export interface Statistics {
  todayTotal: number
  yesterdayTotal: number
  monthTotal: number
  activeKeys: number
  totalAllTime: number
}

// 用户类型
export interface User {
  id: string
  username: string
  role: string
}

