package structutil

import "reflect"

// Clone 可以把 raw 的 公有 属性克隆到 dst
// 使用反射 如果涉及大量字段 请自行编写 Clone
func Clone[T any](raw *T, dst *T) {
	rawValue := reflect.ValueOf(raw).Elem()
	dstValue := reflect.ValueOf(dst).Elem()
	for i := 0; i < rawValue.NumField(); i++ {
		rawField := rawValue.Field(i)
		// 判断是不是私有变量 私有变量就不赋值
		if !rawField.CanInterface() {
			continue
		}
		dstValue.Field(i).Set(rawField)
	}
}
