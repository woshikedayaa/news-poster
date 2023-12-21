package sliceutil

// DeleteElement 只是删除一个 然后迁移
// 无内存开销
func DeleteElement[T any](a []T, index int) []T {
	for i := index; i < len(a)-1; i++ {
		a[i] = a[i+1]
	}

	return a[:len(a)-1]
}
