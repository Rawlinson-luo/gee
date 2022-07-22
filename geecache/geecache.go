package geecache

import "sync"

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

type Group struct {
	name      string
	getter    GetterFunc
	mainCache cache
}

var (
	mu     = sync.Mutex{}
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter GetterFunc) *Group {
	if getter == nil {
		panic("getter is nil")
	}

	mu.Lock()
	defer mu.Unlock()
	group := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}

	groups[name] = group

	return group

}

func GetGroup(name string) (*Group, bool) {
	mu.Lock()
	defer mu.Unlock()
	if group, ok := groups[name]; ok {
		return group, ok
	}

	return nil, false
}

func (g *Group) get(key string) (ByteView, error) {
	if v, ok := g.mainCache.get(key); ok {
		return v, nil
	}

	localValue, err := g.getter(key)

	if err != nil {
		return ByteView{}, err
	}

	localByteView := ByteView{cloneBytes(localValue)}

	g.mainCache.add(key, localByteView)
	return localByteView, nil
}
