package configs

import "sync"

// Mysql
const (
	MysqlWriteKey = "mysql-write"
	MysqlReadKey  = "mysql-read"
)

var onceF sync.Once
var completed bool

func ConfigCompleted() bool {
	return completed
}
