package buff

import (
	"testing"
)

const xmaxInt = int(^uint(0) >> 1)

func testWrite(testNo int, initCap int, p []byte, t *testing.T) {

	buf := NewBuffer(initCap)
	if buf.Len() != 0 || buf.Cap() != initCap {
		t.Fatalf("[test%d]invalid cap %d or len %d", testNo, buf.Cap(), buf.Len())
	}

	var n int
	var err error

	// capacityをまたがないwrite
	loop := initCap / len(p)
	for i := 0; i < loop; i++ {
		n, err = buf.Write(p)
		if err != nil {
			t.Fatal(err)
		}
		if n != len(p) {
			t.Fatalf("[test%d]invalid return size %d", testNo, len(p))
		}
		if buf.Len() != len(p)*(i+1) || buf.Cap() != initCap {
			t.Fatalf("[test%d]invalid cap %d or len %d", testNo, buf.Cap(), buf.Len())
		}
	}

	// capacityをまたぐwrite
	n, err = buf.Write(p)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(p) {
		t.Fatalf("[test%d]invalid return size %d", testNo, len(p))
	}
	if buf.Len() != len(p)*(loop+1) || buf.Cap()%initCap != 0 {
		t.Fatalf("[test%d]invalid cap %d or len %d", testNo, buf.Cap(), buf.Len())
	}

	// bufの中身が正しいか確認
	out := make([]byte, buf.Len())
	buf.CopyTo(out)
	for i, v := range out {
		if v != p[i%len(p)] {
			t.Fatalf("[test%d]invalid buf val[%d] %d", testNo, i, v)
		}
	}

}

func TestWriter(t *testing.T) {

	testWrite(1, 4, []byte{1, 2}, t)
	testWrite(2, 5, []byte{1, 2}, t)
	testWrite(3, 1, []byte{1, 2, 3}, t)

	{
		initCap := 3
		buf := NewBuffer(initCap)
		if buf.Len() != 0 || buf.Cap() != initCap {
			t.Fatalf("invalid cap %d or len %d", buf.Cap(), buf.Len())
		}
		buf.Grow(3)
		if buf.Len() != 3 || buf.Cap() != initCap {
			t.Fatalf("invalid cap %d or len %d", buf.Cap(), buf.Len())
		}
		buf.Grow(3)
		if buf.Len() != 6 || buf.Cap() != initCap*2 {
			t.Fatalf("invalid cap %d or len %d", buf.Cap(), buf.Len())
		}
	}
	{
		initCap := 3
		buf := NewBuffer(initCap)
		buf.Write([]byte{1, 2})
		buf.Write([]byte{1, 2})
		buf.Write([]byte{1, 2})
		buf.Reset()
		if buf.Len() != 0 || buf.Tell() != 0 {
			t.Fatalf("invalid len %d off %d", buf.Len(), buf.Tell())
		}
	}
	{
		initCap := 3
		buf := NewBuffer(initCap)
		buf.WriteAt([]byte{1, 2}, 2)
		if buf.Len() != 4 {
			t.Fatalf("invalid len %d", buf.Len())
		}
		if buf.Cap() != initCap*2 {
			t.Fatalf("invalid cap %d", buf.Cap())
		}
		if buf.Tell() != 4 {
			t.Fatalf("invalid Offset %d", buf.Tell())
		}
	}

}
