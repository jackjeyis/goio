package queue

import "strconv"

type IOBuffer struct {
	buf  []byte
	wt   uint64
	rt   uint64
	size uint64
}

func NewIOBuffer() *IOBuffer {
	return &IOBuffer{}
}

func (b *IOBuffer) EnsureWrite(size uint64) (uint64, error) {
	remain_bytes := b.size - b.wt
	if remain_bytes >= size {
		return b.wt, nil
	}

	data_bytes := b.wt - b.rt
	if b.rt >= size {
		if data_bytes == 0 {
			b.wt = 0
			b.rt = 0
			return b.wt, nil
		} else if b.rt > 1024*1024 {
			copy(b.buf, b.buf[b.rt:b.wt])
			b.wt = data_bytes
			b.rt = 0
			return b.wt, nil
		}
	}

	for b.size < b.wt+size {
		if b.size == 0 {
			b.size = 1024
		} else {
			b.size *= 2
		}
	}

	if b.size > 0xFFFFFFFF {
		return 0, nil
	}

	buffer := make([]byte, b.size)
	if b.buf != nil {
		copy(buffer, b.buf[:])
	}
	b.buf = buffer
	return b.wt, nil
}

func (b *IOBuffer) Write(bs []byte) {
	size := uint64(len(bs))
	wstart, err := b.EnsureWrite(size)
	if err != nil {
		return
	}
	copy(b.buf[wstart:], bs)
	b.Produce(size)
}

func (b *IOBuffer) Read(n uint64) []byte {
	buffer := b.buf[b.rt : b.rt+n]
	b.Consume(n)
	return buffer
}

func (b *IOBuffer) Produce(size uint64) {
	b.wt += size
}

func (b *IOBuffer) Consume(size uint64) {
	b.rt += size
}

func (b *IOBuffer) GetReadSize() uint64 {
	return b.wt - b.rt
}

func (b *IOBuffer) GetRead() uint64 {
	return b.rt
}

func (b *IOBuffer) GetWrite() uint64 {
	return b.wt
}

func (b *IOBuffer) Buffer() []byte {
	return b.buf
}

func (b *IOBuffer) Byte(cnt uint64) []byte {
	return b.buf[b.wt+cnt : b.wt+cnt+1]
}

func (b *IOBuffer) Reset() {
	b.wt = 0
	b.rt = 0
}

func (b *IOBuffer) Summary() string {
	return "[Size:" + strconv.FormatUint(b.size, 10) + ",Wt:" + strconv.FormatUint(b.wt, 10) + ",Rt:" + strconv.FormatUint(b.rt, 10) + "]"
}

/*func main() {
	buf := NewIOBuffer()
	buf.Write([]byte(strings.Repeat("hello", 1024)))
	//fmt.Println(string(buf.buf[:buf.GetReadSize()]))
	fmt.Println(buf.size, buf.wt)
}*/
