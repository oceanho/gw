package gwfunc

import (
	assert2 "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestExec_Normal(t *testing.T) {
	var f = func() {
		time.Sleep(1 * time.Second)
	}
	ok := Timeout(f, 5*time.Second)
	assert2.False(t, ok)
}

func TestExec_Timeout(t *testing.T) {
	var f = func() {
		time.Sleep(5 * time.Second)
	}
	ok := Timeout(f, 1*time.Second)
	assert2.True(t, ok)
}

func BenchmarkExec_Normal(b *testing.B) {
	var f = func() {
		time.Sleep(1 * time.Nanosecond)
	}
	for i := 0; i < b.N; i++ {
		_ = Timeout(f, 1*time.Nanosecond)
	}
}

func BenchmarkExec_Timeout(b *testing.B) {
	var f = func() {
		time.Sleep(5 * time.Nanosecond)
	}
	for i := 0; i < b.N; i++ {
		_ = Timeout(f, 1*time.Nanosecond)
	}
}
