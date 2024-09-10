// https://gist.githubusercontent.com/lojic/de7c3726c91f8160fbeb44b4ff12fef4/raw/901c70d7e5d2f024e2c063d3a4b8566267596dd6/recordsource.go

package blockio

import (
	"bytes"
	"io"
)

const newline = 10

type RecordSource struct {
	begin        int
	numAvailable int
	isEof        bool
	capacity     int
	is           io.Reader
	buffer       []byte
}

func BuildRecordSource(istr io.Reader, capacity int) *RecordSource {
	rs := new(RecordSource)
	rs.is = istr
	rs.capacity = capacity
	rs.buffer = make([]byte, capacity)

	return rs
}

func (rs *RecordSource) FillBuffer() {
	if !rs.isEof {

		if rs.begin > 0 {
			copy(rs.buffer, rs.buffer[rs.begin:rs.begin+rs.numAvailable])
			rs.begin = 0
		}

		num_read, err := rs.is.Read(rs.buffer[rs.numAvailable:rs.capacity])

		if err != nil {
			if err == io.EOF {
				rs.isEof = true
			} else {
				panic(err)
			}
		}

		rs.numAvailable += num_read
	}
}

func (rs *RecordSource) NextLine() []byte {
	for loop := 1; loop <= 2; loop++ {
		idx := bytes.IndexByte(rs.buffer[rs.begin:rs.begin+rs.numAvailable], newline)

		if idx >= 0 {
			len := idx + 1
			line := rs.buffer[rs.begin : rs.begin+len]
			rs.begin += len
			rs.numAvailable -= len

			return line
		} else {
			rs.FillBuffer()

			if rs.isEof {
				line := rs.buffer[rs.begin : rs.begin+rs.numAvailable]
				rs.begin += rs.numAvailable
				rs.numAvailable = 0

				return line
			}
		}
	}

	panic("blockio.RecordSource.NextLine(): Could not find line")
}
