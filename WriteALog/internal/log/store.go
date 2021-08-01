package log

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
	"sync"
)

// encoding for storing record sizes
var enc = binary.BigEndian
// number of bytes used to store a records length
const lenWidth = 8


type store struct {
	// why don't we just make this explicitly define Read() and Write()
	// since ideally we could use memory for this
	// like how Pete Bourgon injests inMemory and also with files
	// well I do realise *os.File houses lots of code
	// left to see what we actually use
	*os.File
	buf *bufio.Writer
	size uint64
	mu sync.Mutex
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	return &store{
		File: f,
		size: uint64(fi.Size()),
		buf: bufio.NewWriter(f),
	}, nil
}

func (s *store) Append(p []byte) (numWritten uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size

	// write length of record, so we know how to read
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, pos, err
	}

	n, err := s.buf.Write(p)

	if err != nil {
		return 0, pos, err
	}

	// here, no error but we could no write all the bytes
	// just golang things
	if n != len(p) {
		return uint64(n), pos, err
	}

	n  += lenWidth
	s.size += uint64(n)

	return uint64(n), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	if pos > s.size {
		return nil, fmt.Errorf("out of bounds, requested pos %v where limit is %v", pos, s.size)
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	size := make([]byte, lenWidth)
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	b := make([]byte, enc.Uint64(size))

	if _, err := s.File.ReadAt(b, int64(pos + lenWidth)); err != nil {
		// TODO check for _ < len(b) which indicated EOF
		return nil, err
	}
	return b, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	return s.File.ReadAt(p, off)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}

	return s.File.Close()
}
