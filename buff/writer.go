package buff

import (
	"errors"
)

const smallBufferSize = 64
const maxInt = int(^uint(0) >> 1)

// Buffer ...
type Buffer struct {
	buf []byte
	off int // next write at buf[off]
}

// NewBuffer ...
func NewBuffer(n int) *Buffer {
	buf := makeSlice(n)
	buffer := &Buffer{buf: buf}
	buffer.Reset()
	return buffer
}

// Reset ...
func (b *Buffer) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
}

// Bytes ...
func (b *Buffer) Bytes() []byte {
	return b.buf[:]
}

// CopyTo ...
func (b *Buffer) CopyTo(dst []byte) {
	copy(dst, b.buf)
}

// Len ...
func (b *Buffer) Len() int {
	return len(b.buf)
}

// Cap ...
func (b *Buffer) Cap() int {
	return cap(b.buf)
}

// Seek ...
func (b *Buffer) Seek(n int) {
	b.off = n
}

// Tell ...
func (b *Buffer) Tell() int {
	return b.off
}

// Grow ...
func (b *Buffer) Grow(n int) {
	if n < 0 {
		panic("buff.Buffer.Grow: negative count")
	}
	b.grow(n)
}

// Write ...
func (b *Buffer) Write(p []byte) (n int, err error) {
	if growSize := b.neededGrowSize(b.off, len(p)); 0 < growSize {
		if ok := b.tryGrowByReslice(growSize); !ok {
			b.grow(growSize)
		}
	}

	n = b.writeToBuff(p)

	return n, nil
}

// WriteAt ...
func (b *Buffer) WriteAt(p []byte, off int64) (n int, err error) {
	b.off = int(off)
	return b.Write(p)
}

func ceil(v, unit int) int {
	return (v + (unit - 1)) / unit * unit
}

func (b *Buffer) grow(n int) {
	/*
	 * ex.)
	 * len(b.buf) is 2, cap(b.buf) is 6, n is 3.
	 * before : xxoooo
	 * after  : xxxxxo
	 * x is used index, o is free index.
	 */
	m := len(b.buf)
	newBufLen := m + n

	// 空き領域が足りていればlengthだけを伸ばす
	if ok := b.tryGrowByReslice(n); ok {
		return
	}

	// 新しいスライスのcapサイズを計算する
	c := cap(b.buf)
	var unit int
	if c == 0 {
		unit = smallBufferSize
	} else {
		unit = c
	}
	newCapSize := ceil(newBufLen, unit)

	// capacityが足りていないから新しくスライスをつくる
	newBuf := makeSlice(newCapSize)
	copy(newBuf, b.buf)

	// lengthを必要な分だけ伸ばす
	b.buf = newBuf[:newBufLen]
	return
}

func (b *Buffer) tryGrowByReslice(n int) bool {
	bufLen := len(b.buf)
	free := cap(b.buf) - bufLen
	if n <= free {
		b.buf = b.buf[:bufLen+n]
		return true
	}
	return false
}

func (b *Buffer) neededGrowSize(off, writeSize int) int {
	bufLen := len(b.buf)

	if off+writeSize <= bufLen {
		return 0
	}

	n := writeSize
	if bufLen <= off {
		// ライトする位置がlenの範囲外
		n += int(off - bufLen)
	} else {
		// ライトする位置がlenの範囲内
		// length内で空いているサイズ分は確保するサイズから引く
		n -= int(bufLen - off)
	}

	if n < 0 {
		n = 0
	}

	return n
}

func (b *Buffer) writeToBuff(p []byte) (n int) {
	n = copy(b.buf[b.off:], p)
	b.off += n
	return n
}

func makeSlice(n int) []byte {
	// If the make fails, give a known error.
	defer func() {
		if recover() != nil {
			panic(errors.New("buff.Buffer: too large"))
		}
	}()
	return make([]byte, n)
}
