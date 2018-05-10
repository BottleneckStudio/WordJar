package cache_test

import (
	"testing"

	"github.com/BottleneckStudio/WordJar/cache"
)

func TestSetCache(t *testing.T) {
	got, err := cache.Set("testKey", "test value", 360)
	want := true

	if err != nil {
		t.Errorf("Memcached is not running. Please check: %v", err)
	}

	if want != got {
		t.Errorf("Want: %t, but got %t", want, got)
	}

	t.Logf("Successfully saved to cache")
}
