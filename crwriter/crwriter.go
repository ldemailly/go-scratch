package crwriter

import "io"

// CRWriter is a writer that adds \r before each \n.
// Needed for raw mode terminals (and I guess also if you want something DOS or http headers like).
type CRWriter struct {
	buf []byte
	Out io.Writer
}

// In case you want to ensure the memory used by the buffer is released.
func (c *CRWriter) Reset() {
	c.buf = nil
}

// Optimized to avoid many small writes by buffering and only writing \r when needed.
// No extra syscall, relies on append() to efficiently reallocate the buffer.
func (c *CRWriter) Write(orig []byte) (n int, err error) {
	l := len(orig)
	if l == 0 {
		return 0, nil
	}
	if l == 1 {
		if orig[0] != '\n' {
			return c.Out.Write(orig)
		}
		_, err = c.Out.Write([]byte("\r\n"))
		return 1, err
	}
	lastEmitted := 0
	for i, b := range orig {
		if b != '\n' { // IndexByte is probably faster than this.
			continue
		}
		// leave the \n for next append. I wish I could write
		//   c.buf = append(c.buf, orig[lastEmitted:i]..., `\r`)
		// instead of the 2 lines.
		c.buf = append(c.buf, orig[lastEmitted:i]...)
		c.buf = append(c.buf, '\r')
		lastEmitted = i
	}
	if lastEmitted == 0 {
		return c.Out.Write(orig)
	}
	c.buf = append(c.buf, orig[lastEmitted:]...)
	_, err = c.Out.Write(c.buf)
	return len(orig), err // in case caller checks... but we might have written "more".
}
