package memcachier

import "github.com/memcachier/mc"

// Memcachier ...
type Memcachier struct {
	memcachierClient *mc.Client
}

// Config configures the memcachier
// implementation
type Config struct {
	Server   string
	Username string
	Password string
}

// NewMemcachier returns instance of memcachier
func NewMemcachier(config Config) *Memcachier {
	return &Memcachier{
		memcachierClient: mc.NewMC(config.Server, config.Username, config.Password),
	}
}

// Set returns a boolean value
// after setting a value using
// the specified `key`, returns
// error otherwise.
func (memcachier *Memcachier) Set(key string, data interface{}) (bool, error) {
	_, err := memcachier.memcachierClient.Set(key, data.(string), uint32(0), uint32(0), uint64(0))

	if err != nil {
		return false, err
	}

	return true, nil
}

// Get returns the `data` saved in cache
// using the specified `key`.
func (memcachier *Memcachier) Get(key string) (interface{}, error) {
	val, _, _, err := memcachier.memcachierClient.Get(key)

	if err != nil {
		return nil, err
	}
	return val, nil
}

// Delete returns a boolean value
// if there is a successful deletion
// using the specified `key`,
// returns error otherwise.
func (memcachier *Memcachier) Delete(key string) (bool, error) {
	err := memcachier.memcachierClient.Del(key)

	if err != nil {
		return false, err
	}

	return true, nil
}
