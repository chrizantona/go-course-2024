package storage

import (
  "math/rand"
  "strconv"
  "testing"
)

type testCase struct {
  name  string
  key   string
  value string
  kind  Kind
}

type benchCase struct {
  name string
  cnt  int
  // key string
  // value string
  // kind string
}

var cases = []benchCase{
  {"1", 1},
  {"10", 10},
  {"100", 100},
  {"1000", 1000},
  {"10000", 10000},
}

func BenchmarkSet(b *testing.B) {
  for _, tCase := range cases {
    b.Run(tCase.name, func(b *testing.B) {
      s := NewStorage()

      for i := 0; i < tCase.cnt; i++ {
        s.Set(strconv.Itoa(i), strconv.Itoa(i))
      }

      b.ResetTimer()
      for i := 0; i < b.N; i++ {
        s.Set(strconv.Itoa(rand.Intn(tCase.cnt)), "fkjashdf")
      }
    })
  }
}

func BenchmarkGet(b *testing.B) {
  for _, tCase := range cases {
    b.Run(tCase.name, func(b *testing.B) {
      s := NewStorage()

      for i := 0; i < tCase.cnt; i++ {
        s.Set(strconv.Itoa(i), strconv.Itoa(i))
      }

      b.ResetTimer()
      for i := 0; i < b.N; i++ {
        s.Get(strconv.Itoa(rand.Intn(tCase.cnt)))
      }
    })
  }
}

func BenchmarkGetSet(b *testing.B) {
  for _, tCase := range cases {
    b.Run(tCase.name, func(b *testing.B) {
      s := NewStorage()

      for i := 0; i < tCase.cnt; i++ {
        s.Set(strconv.Itoa(i), strconv.Itoa(i))
      }

      b.ResetTimer()

      for i := 0; i < b.N; i++ {
        s.Set(strconv.Itoa(i), strconv.Itoa(i))
        s.Get(strconv.Itoa(rand.Intn(tCase.cnt)))
      }
    })
  }
}

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
		s.Set(c.key, c.value)
		sValue := s.Get(c.key)
  
		if *sValue != c.value {
		  t.Errorf("expected %v, got %v", c.value, *sValue)
		}
	  })
	}
  }

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
		s.Set(c.key, c.value)
		sKind := s.GetKind(c.key)
  
		if sKind != string(c.kind) {
		  t.Errorf("expected kind %v, got %v", c.kind, sKind)
		}
	  })
	}
  }