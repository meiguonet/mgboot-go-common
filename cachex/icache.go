package cachex

type ICache interface {
	Get(key string, defaultValue ...interface{}) interface{}
	Set(key string, value interface{}, ttl ...interface{}) bool
	Delete(key string) bool
	Clear() bool
	GetMultiple(keys []string, defaultValue ...interface{}) []interface{}
	SetMultiple(entries []map[string]interface{}, ttl ...interface{}) bool
	DeleteMultiple(keys []string) bool
	Has(key string) bool
}
