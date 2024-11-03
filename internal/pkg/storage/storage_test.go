package storage

import (
	"github.com/stretchr/testify/require"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

type testCase struct {
	name  string
	key   string
	value string
	kind  Kind
}

type benchCase struct {
	cnt int
}

var cases = []benchCase{
	{1},
	{10},
	{100},
	{1000},
	{10000},
}

func BenchmarkSet(b *testing.B) {
	for i, tCase := range cases {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			s := NewStorage()
			for j := 0; j < tCase.cnt; j++ {
				s.Set(strconv.Itoa(j), strconv.Itoa(j), 10*time.Second) 
			}

			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				s.Set(strconv.Itoa(rand.Intn(tCase.cnt)), "fkjashdf", 10*time.Second) 
			}
		})
	}
}


// goos: darwin
// goarch: arm64
// pkg: go-course-2024/internal/pkg/storage
// cpu: Apple M1
// BenchmarkSet/0-8        23786077                50.37 ns/op
// BenchmarkSet/1-8        17505363                67.57 ns/op
// BenchmarkSet/2-8        16961250                70.31 ns/op
// BenchmarkSet/3-8        14643208                81.18 ns/op
// BenchmarkSet/4-8        12635286                95.27 ns/op
// PASS
// ok      go-course-2024/internal/pkg/storage     7.587s

func BenchmarkGet(b *testing.B) {
	for i, tCase := range cases {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			s := NewStorage()
			for j := 0; j < tCase.cnt; j++ {
				s.Set(strconv.Itoa(j), strconv.Itoa(j), 10*time.Second)
			}
			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				s.Get(strconv.Itoa(rand.Intn(tCase.cnt)))
			}
		})
	}
}


// goos: darwin
// goarch: arm64
// pkg: go-course-2024/internal/pkg/storage
// cpu: Apple M1
// BenchmarkGet/0-8        18357508                65.42 ns/op          144 B/op          2 allocs/op
// BenchmarkGet/1-8        14298465                82.42 ns/op          144 B/op          2 allocs/op
// BenchmarkGet/2-8        13514629                87.08 ns/op          144 B/op          2 allocs/op
// BenchmarkGet/3-8        11473394               102.2 ns/op           146 B/op          2 allocs/op
// BenchmarkGet/4-8        10437296               113.8 ns/op           147 B/op          2 allocs/op
// PASS
// ok      go-course-2024/internal/pkg/storage     7.509s

func BenchmarkGetSet(b *testing.B) {
	for i, tCase := range cases {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			s := NewStorage()
			for j := 0; j < tCase.cnt; j++ {
				s.Set(strconv.Itoa(j), strconv.Itoa(j), 10*time.Second) 
			}

			var result string

			b.ResetTimer()
			for j := 0; j < b.N; j++ {
				s.Set(strconv.Itoa(j), strconv.Itoa(j), 10*time.Second) 
				if val, err := s.Get(strconv.Itoa(rand.Intn(tCase.cnt))); err == nil {
					if strVal, ok := (*val).(string); ok {
						result = strVal
					} else {
						b.Error("Expected value to be of type string")
					}
				}
			}

			if result == "" {
				b.Error("Expected result to be non-empty")
			}
		})
	}
}



// goos: darwin
// goarch: arm64
// pkg: go-course-2024/internal/pkg/storage
// cpu: Apple M1
// BenchmarkGetSet/0-8      2741875               378.7 ns/op
// BenchmarkGetSet/1-8      3131571               384.2 ns/op
// BenchmarkGetSet/2-8      2968576               392.1 ns/op
// BenchmarkGetSet/3-8      3017692               409.2 ns/op
// BenchmarkGetSet/4-8      2647934               463.0 ns/op
// PASS
// ok      go-course-2024/internal/pkg/storage     8.799s

func TestSetGetBasic(t *testing.T) {
	cases := []testCase{
		{"1", "testKey1", "testValue1", KindString},
		{"2", "testKey2", "123", KindInt}, 
		{"3", "testKey3", "hello", KindString},
		{"4", "testKey4", "456", KindInt}, 
	}

	s := NewStorage()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := s.Set(c.key, c.value, 10*time.Second)
			if err != nil {
				t.Errorf("Unexpected error on Set: %v", err)
				return
			}

			sValue, err := s.Get(c.key)
			if err != nil {
				t.Errorf("Unexpected error on Get: %v", err)
				return
			}

			if sValue == nil || *sValue != c.value {
				t.Errorf("expected %v, got %v", c.value, *sValue)
			}
		})
	}
}


// PASS
// ok      go-course-2024/internal/pkg/storage     0.594s

func TestSetTypeDetermination(t *testing.T) {
	cases := []testCase{
		{"1", "intKey", "10", KindInt},   
		{"2", "stringKey", "keem", KindString},  
		{"3", "anotherIntKey", "42", KindInt},   
		{"4", "anotherStringKey", "kendrick", KindString}, 
	}

	s := NewStorage()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := s.Set(c.key, c.value, 10*time.Second)
			if err != nil {
				t.Errorf("Unexpected error on Set: %v", err)
				return
			}

			sKind, err := s.GetKind(c.key)
			if err != nil {
				t.Errorf("Unexpected error on GetKind: %v", err)
				return
			}

			require.Equal(t, string(c.kind), sKind, "expected kind %v, got %v", c.kind, sKind)
		})
	}
}

// PASS
// ok      go-course-2024/internal/pkg/storage     0.514s
