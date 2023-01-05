package config

import (
	"strconv"
	"testing"
)

func BenchmarkConfig_ExpShortURL(b *testing.B) {
	c := NewConfig(WithParseEnv())
	for i := 0; i < b.N; i++ {
		c.ExpShortURL(strconv.Itoa(i))
	}
}
