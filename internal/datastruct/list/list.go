package list

type RedisList interface {
	Len() int

	// LPush(values ...string) int
	RPush(values ...string) int
	//
	// LPop() (string, bool)
	// RPop() (string, bool)
	//
	// LIndex(index int) (string, bool)
	// LSet(index int, value string) error
	//
	// LRange(start, stop int) []string
	//
	// LInsertBefore(pivot, value string) (int, bool)
	// LInsertAfter(pivot, value string) (int, bool)
	//
	// LRem(count int, value string) int
	// LTrim(start, stop int)
}
