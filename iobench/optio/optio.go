// Line Scanner with minimized copying.
// We wish we'd do a circular buffer and readv() into it but that's not available in Go (outside of linux for some reason
// despite posix).
package optio

import (
	"bytes"
	"errors"
	"io"
)

type SlidingBuffer struct {
	buf           []byte
	start         int
	length        int
	readsize      int
	capacity      int
	maxLineLength int
	f             io.Reader
	eof           bool
}

const newLine = '\n'

var ErrNoEOL = errors.New("no EOL found within maximum line length")

func LineScanner(f io.Reader, bufferSize, readSize, maxLineLength int) *SlidingBuffer {
	if readSize > bufferSize/2 {
		panic("readsize must be less than half of the size of the buffer (1/4 or more is better)")
	}
	if maxLineLength > readSize {
		panic("maxLineLength must be smaller than readSize")
	}
	sb := &SlidingBuffer{
		buf:           make([]byte, bufferSize),
		start:         0,
		length:        0,
		readsize:      readSize,
		capacity:      bufferSize,
		maxLineLength: maxLineLength,
		f:             f,
	}
	return sb
}

func (sb *SlidingBuffer) restOfBuffer() []byte {
	return sb.buf[sb.start : sb.start+sb.length]
}

func (sb *SlidingBuffer) EOF() bool {
	return sb.eof && (sb.length == 0) // make sure last data is read through Line() first
}

func (sb *SlidingBuffer) Line() ([]byte, error) {
	if sb.length < sb.maxLineLength {
		err := sb.read() // EOF isn't returned here
		if err != nil {
			return sb.restOfBuffer(), err
		}
	}
	idx := bytes.IndexByte(sb.buf[sb.start:sb.start+sb.length], newLine)
	if idx < 0 {
		err := ErrNoEOL
		if sb.eof {
			err = nil
		}
		res := sb.restOfBuffer()
		sb.start = 0
		sb.length = 0
		return res, err
	}
	res := sb.buf[sb.start : sb.start+idx]
	idx++ // consume the newline
	sb.start += idx
	sb.length -= idx
	return res, nil
}

func (sb *SlidingBuffer) read() error {
	end := sb.start + sb.length
	if end+sb.readsize >= sb.capacity {
		sb.compact()
		end = sb.length
	}
	// readMult := (sb.capacity - end) / sb.readsize // read un-fragmented units.
	n, err := sb.f.Read(sb.buf[end : end+ /* readMult* */ sb.readsize])
	sb.length += n
	if err == io.EOF {
		sb.eof = true
		return nil
	}
	return err
}

func (sb *SlidingBuffer) compact() {
	if sb.length == 0 {
		sb.start = 0
		return
	}
	if sb.start == 0 {
		return
	}
	copy(sb.buf, sb.buf[sb.start:sb.start+sb.length])
	sb.start = 0
}
