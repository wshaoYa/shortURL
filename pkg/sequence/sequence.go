package sequence

// Sequence 发号器
type Sequence interface {
	Next() (uint64, error)
}
