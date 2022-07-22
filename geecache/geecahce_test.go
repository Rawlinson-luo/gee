package geecache

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

func TestGetter(t *testing.T) {
	f := GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	except := []byte("key")

	if v, _ := f.Get("key"); !reflect.DeepEqual(v, except) {
		t.Fatalf("call back error")
	}

}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))

	gee := NewGroup("scores", 2<<10, GetterFunc(func(key string) ([]byte, error) {
		//log.Println("[slow db] search key: ", key)
		//if v, ok := db[key]; ok {
		//	if _, ok := loadCounts[key]; !ok {
		//		loadCounts[key] = 0
		//	}
		//	loadCounts[key] += 1
		//	return []byte(v), nil
		//}
		//return nil, fmt.Errorf("%s not exist\n", key)
		log.Println("[SlowDB] search key", key)
		if v, ok := db[key]; ok {
			if _, ok := loadCounts[key]; !ok {
				loadCounts[key] = 0
			}
			loadCounts[key] += 1
			return []byte(v), nil
		}
		return nil, fmt.Errorf("%s not exist", key)
	}))

	for k, v := range db {
		if view, err := gee.get(k); err != nil || view.String() != v {
			t.Fatalf("can't get %s\n", k)
		}

		if _, err := gee.get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss\n", k)
		}
	}

	if v, err := gee.get("unknown"); err == nil {
		t.Fatalf("the key unknown should be empty, but get value %s\n", v)
	}

}
