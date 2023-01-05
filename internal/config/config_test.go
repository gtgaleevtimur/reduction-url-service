package config

import "testing"

func BenchmarkNewConfig(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewConfig(WithParseEnv())
	}
}
