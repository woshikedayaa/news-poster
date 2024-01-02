// Package configs
// 这个包在初始化的时候上传config.yaml 中的配置
// 先查看etcd中的配置是否配置完成 如果配置完成了就不进行操作
// 如果没有配置完成就先获取 分布式锁 获取成功就将配置上传上去
// 上传完成就把配置完成的flag设置为true 这样其他微服务就不需要再更新
// 其他微服务就watch这个配置key 然后有反应自己更新就行
package configs

// Mysql
const (
	MysqlWrite         = "mysql.write"
	MysqlWriteAddr     = MysqlWrite + ".addr"
	MysqlWriteUser     = MysqlWrite + ".user"
	MysqlWritePassword = MysqlWrite + ".password"
	MysqlWriteArgs     = MysqlWrite + ".args"

	MysqlRead         = "mysql.read"
	MysqlReadeAddr    = MysqlRead + ".addr"
	MysqlReadUser     = MysqlRead + ".user"
	MysqlReadPassword = MysqlRead + ".password"
	MysqlReadArgs     = MysqlRead + ".args"
)

// Redis
const (
	RedisCluster     = "redis.cluster"
	RedisClusterAddr = RedisCluster + ".addr"
	// RedisClusterAuth 写 "" 为没有验证
	RedisClusterAuth = RedisCluster + ".auth"
)
