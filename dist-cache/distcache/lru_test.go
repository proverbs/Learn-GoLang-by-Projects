package distcache

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestAdd(t *testing.T) {
	c := New(int64(0), nil)
	c.Add("name", String("proverbs"))
	c.Add("name", String("pro"))

	if c.nbytes != int64(len("name") + len("pro")) {
		t.Fatal("expected 7 but got", c.nbytes)
	}
}

func TestGet(t *testing.T) {
	c := New(int64(0), nil)
	c.Add("name", String("pro"))

	if v, ok := c.Get("name"); !ok || string(v.(String)) != "pro" {
		t.Fatalf("test cache hit name=pro failed")
	}
	if _, ok := c.Get("key"); ok {
		t.Fatalf("test cache miss key failed")
	}
}

func TestEvict(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	c := New(int64(cap), nil)
	c.Add(k1, String(v1))
	c.Add(k2, String(v2))
	c.Add(k3, String(v3))

	if _, ok := c.Get("key1"); ok || c.Len() != 2 {
		t.Fatalf("test evict key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, val Val) {
		keys = append(keys, key)
	}
	c := New(int64(10), callback)
	c.Add("key1", String("123456"))
	c.Add("k2", String("k2"))
	c.Add("k3", String("k3"))
	c.Add("k4", String("k4"))

	expected := []string{"key1", "k2"}

	if !reflect.DeepEqual(expected, keys) {
		t.Fatalf("Call OnEvicted failed, expected keys are %s", expected)
	}
}
