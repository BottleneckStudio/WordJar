package cache

// Cacher defines the different
// methods of caching
type Cacher interface {
	Set(string, interface{}) (bool, error)
	Get(string) (interface{}, error)
	Delete(string) (bool, error)
}

// Set ...
func Set(key string, data interface{}, cacher Cacher) (bool, error) {
	return cacher.Set(key, data)
}

// Get ...
func Get(key string, cacher Cacher) (interface{}, error) {
	return cacher.Get(key)
}

// Delete ...
func Delete(key string, cacher Cacher) (bool, error) {
	return cacher.Delete(key)
}
