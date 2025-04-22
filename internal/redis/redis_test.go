package redis

import (
	"testing"
)

func TestRedis(t *testing.T) {
	e := Set("ss")
	if e != nil {
		t.Errorf("Erro ao definir valor no Redis: %v", e)
	}
	// val, err := Get("key")
}
