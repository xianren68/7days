package cache

// ByteView a read only secquence of bytes.
type ByteView struct {
	bs []byte
}

// Len return the number of bytes occupied by the byte secquence.
func (b ByteView) Len() int {
	return len(b.bs)
}

// ByteSlice return a clone data of byte secquence.
func (b ByteView) ByteSlice() []byte {
	return b.cloneBytes()
}

// cloneBytes return clone data to a new slice.
func (b ByteView) cloneBytes() []byte {
	bytes := make([]byte, 0, len(b.bs))
	return append(bytes, b.bs...)
}

// String return the string of bytes.
func (b ByteView) String() string {
	return string(b.bs)
}
