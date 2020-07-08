package distcache

type ByteView struct {
	data []byte
}

func (b ByteView) Len() int {
	return len(b.data)
}

func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.data)
}

func (b ByteView) String() string {
	return string(b.data)
}

func cloneBytes(data []byte) []byte {
	newData := make([]byte, len(data))
	copy(newData, data)
	return newData
}
