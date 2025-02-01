package btree

type BTree interface {
	Get(key string) (string, bool)
	Append(key, existingVal, val string)
	IterateKey(prefixKey string, fn func(key string))
	Count() int
	Clear()
}
