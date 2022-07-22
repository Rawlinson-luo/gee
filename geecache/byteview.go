package geecache

type ByteView struct {
	b []byte
}

func (b ByteView) Len() int {
	return len(b.b)
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}

func (b ByteView) String() string {
	return string(b.b)
}

func cloneBytes(srcBytes []byte) []byte {
	bytes := make([]byte, len(srcBytes))
	copy(bytes, srcBytes)
	return bytes
}
