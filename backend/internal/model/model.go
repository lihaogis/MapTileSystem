package model

import (
	"time"

	"github.com/lib/pq"
)

// DataSource 数据源模型
// 用于存储地图瓦片数据源的配置信息，支持 XYZ 栅格瓦片和 3D Tiles 两种类型
type DataSource struct {
	ID          string    `json:"id" gorm:"primaryKey;comment:数据源唯一标识"`                    // 数据源唯一标识，UUID 格式
	Name        string    `json:"name" gorm:"not null;comment:数据源名称"`                       // 数据源名称，用于展示和识别
	Type        string    `json:"type" gorm:"not null;comment:数据源类型"`                       // 数据源类型：xyz（栅格瓦片）、3dtiles（三维瓦片）
	Format      string    `json:"format" gorm:"comment:瓦片格式"`                               // 瓦片格式：png、jpeg、webp（仅 XYZ 类型需要）
	Path        string    `json:"path" gorm:"not null;comment:存储路径"`                       // 瓦片数据存储路径，可以是本地路径或云存储路径
	MinZoom     int       `json:"minZoom" gorm:"comment:最小缩放级别"`                            // 最小缩放级别（仅 XYZ 类型需要）
	MaxZoom     int       `json:"maxZoom" gorm:"comment:最大缩放级别"`                            // 最大缩放级别（仅 XYZ 类型需要）
	CenterLat   *float64  `json:"centerLat" gorm:"comment:地图中心纬度"`                         // 地图预览中心纬度（可选）
	CenterLng   *float64  `json:"centerLng" gorm:"comment:地图中心经度"`                         // 地图预览中心经度（可选）
	DefaultZoom *int      `json:"defaultZoom" gorm:"comment:默认缩放级别"`                       // 地图预览默认缩放级别（可选）
	Status      string    `json:"status" gorm:"default:'enabled';comment:状态"`             // 状态：enabled（启用）、disabled（禁用）
	CreatedAt   time.Time `json:"createdAt" gorm:"comment:创建时间"`                           // 创建时间
	UpdatedAt   time.Time `json:"updatedAt" gorm:"comment:更新时间"`                           // 更新时间
}

// TableName 指定表名
func (DataSource) TableName() string {
	return "data_sources"
}

// ApiKey API 密钥模型
// 用于管理访问瓦片服务的 API Key，支持权限控制和访问限制
type ApiKey struct {
	ID               string     `json:"id" gorm:"primaryKey;comment:API Key 唯一标识"`                      // API Key 唯一标识，UUID 格式
	Key              string     `json:"key" gorm:"uniqueIndex;not null;comment:API 密钥"`                 // API 密钥，32位随机字符串，用于鉴权
	Name             string     `json:"name" gorm:"not null;comment:密钥名称"`                              // 密钥名称，用于标识用途（如：生产环境、测试环境）
	DataSources      pq.StringArray `json:"dataSources" gorm:"type:text[];comment:授权的数据源列表"`               // 授权的数据源 ID 列表
	IPWhitelist      string     `json:"ipWhitelist" gorm:"type:text;comment:IP 白名单"`                    // IP 白名单，JSON 数组格式，支持 CIDR 格式
	RefererWhitelist string     `json:"refererWhitelist" gorm:"type:text;comment:Referer 白名单"`          // Referer 白名单，JSON 数组格式，用于防盗链
	Status           string     `json:"status" gorm:"default:'enabled';comment:状态"`                    // 状态：enabled（启用）、disabled（禁用）
	LastUsedAt       *time.Time `json:"lastUsedAt" gorm:"comment:最后使用时间"`                               // 最后使用时间，记录最近一次调用的时间
	CallCount        int64      `json:"callCount" gorm:"default:0;comment:调用次数"`                       // 累计调用次数
	CreatedAt        time.Time  `json:"createdAt" gorm:"comment:创建时间"`                                 // 创建时间
	UpdatedAt        time.Time  `json:"updatedAt" gorm:"comment:更新时间"`                                 // 更新时间
}

// TableName 指定表名
func (ApiKey) TableName() string {
	return "api_keys"
}

// CallLog 调用日志模型
// 用于记录每次瓦片请求的详细信息，用于统计分析和监控
type CallLog struct {
	ID           uint      `json:"id" gorm:"primaryKey;comment:日志唯一标识"`                  // 日志唯一标识，自增主键
	ApiKeyID     string    `json:"apiKeyId" gorm:"index;comment:API Key ID"`            // 调用的 API Key ID
	DataSourceID string    `json:"dataSourceId" gorm:"index;comment:数据源 ID"`            // 访问的数据源 ID
	TileZ        int       `json:"tileZ" gorm:"comment:瓦片 Z 坐标（缩放级别）"`                   // 瓦片 Z 坐标（缩放级别）
	TileX        int       `json:"tileX" gorm:"comment:瓦片 X 坐标"`                         // 瓦片 X 坐标
	TileY        int       `json:"tileY" gorm:"comment:瓦片 Y 坐标"`                         // 瓦片 Y 坐标
	StatusCode   int       `json:"statusCode" gorm:"comment:HTTP 状态码"`                  // HTTP 响应状态码（200、404、403 等）
	ResponseTime int       `json:"responseTime" gorm:"comment:响应时间（毫秒）"`                // 响应时间，单位：毫秒
	IPAddress    string    `json:"ipAddress" gorm:"comment:客户端 IP 地址"`                  // 客户端 IP 地址
	UserAgent    string    `json:"userAgent" gorm:"comment:客户端 User-Agent"`             // 客户端 User-Agent 信息
	CreatedAt    time.Time `json:"createdAt" gorm:"index;comment:请求时间"`                 // 请求时间，用于统计和查询
}

// TableName 指定表名
func (CallLog) TableName() string {
	return "call_logs"
}

// User 用户模型
// 用于管理系统管理员账户
type User struct {
	ID        string    `json:"id" gorm:"primaryKey;comment:用户唯一标识"`                    // 用户唯一标识，UUID 格式
	Username  string    `json:"username" gorm:"uniqueIndex;not null;comment:用户名"`       // 用户名，用于登录，唯一
	Password  string    `json:"-" gorm:"not null;comment:密码哈希"`                         // 密码哈希值，不返回给前端
	Role      string    `json:"role" gorm:"default:'admin';comment:角色"`                 // 用户角色：admin（管理员）
	CreatedAt time.Time `json:"createdAt" gorm:"comment:创建时间"`                          // 创建时间
	UpdatedAt time.Time `json:"updatedAt" gorm:"comment:更新时间"`                          // 更新时间
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// StatisticsSummary 统计汇总模型
// 按天汇总调用统计，用于长期数据保留
type StatisticsSummary struct {
	ID           uint      `json:"id" gorm:"primaryKey;comment:汇总记录唯一标识"`           // 汇总记录唯一标识，自增主键
	Date         time.Time `json:"date" gorm:"uniqueIndex;comment:统计日期"`            // 统计日期（按天）
	TotalCalls   int64     `json:"totalCalls" gorm:"comment:总调用次数"`                // 当天总调用次数
	SuccessCalls int64     `json:"successCalls" gorm:"comment:成功调用次数"`             // 当天成功调用次数（状态码 200）
	FailedCalls  int64     `json:"failedCalls" gorm:"comment:失败调用次数"`              // 当天失败调用次数（状态码非 200）
	CreatedAt    time.Time `json:"createdAt" gorm:"comment:创建时间"`                  // 创建时间
}

// TableName 指定表名
func (StatisticsSummary) TableName() string {
	return "statistics_summary"
}
