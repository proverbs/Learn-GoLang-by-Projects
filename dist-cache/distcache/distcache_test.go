package distcache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

var loadCounts = make(map[string]int, len(db))

var (
	defaultGetter Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	dbGetter Getter = GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		})
)

func TestGetGroup(t *testing.T) {
	groupName := "scores"
	NewGroup(groupName, 2<<10, defaultGetter)
	
	if group := GetGroup(groupName); group == nil || group.name != groupName {
		t.Fatalf("group %s not exist", groupName)
	}

	if group := GetGroup(groupName + "111"); group != nil {
		t.Fatalf("expected nil, but %s got", group.name)
	}
}

func TestGetter(t *testing.T) {
	var f Getter = defaultGetter

	expected := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expected) {
		t.Fatal("test callback failed")
	}
}

func TestGet(t *testing.T) {
	g := NewGroup("scores", 2<<10, dbGetter)

	for k, v := range db {
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatal("failed to get value of Tom")
		}
		if _, err := g.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}

	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("the value of unknown should be empty, but %s got", view)
	}
}
