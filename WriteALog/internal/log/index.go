package log

import (
	"io"
	"os"

	"github.com/tysontate/gommap"
)

const (
	offWidth   uint64 = 4
	posWidth   uint64 = 8
	entryWidth        = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}

	fi, err := os.Stat(f.Name())

	if err != nil {
		return nil, err
	}

	idx.size = uint64(fi.Size())

	if err = os.Truncate(f.Name(), int64(c.Segment.MaxIndexBytes)); err != nil {
		return nil, err
	}

	if idx.mmap, err = gommap.Map(idx.file.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED); err != nil {
		return nil, err
	}

	return idx, nil
}

func (i *index) Read(offset int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}

	// i really do not understand this
	// -1 is some special value?
	// the offset itself is uint i guess
	if offset == -1 {
		// how many items actually in index
		// then first one?
		// so some special way to get the first value?
		// its zero indexed
		// so -1 makes sense
		out = uint32((i.size / entryWidth) - 1)
	} else {
		out = uint32(offset)
	}

	// get actual byte based position
	pos = uint64(out) * entryWidth
	// if what we'd read is outside the range, then EOF
	if i.size < (pos + entryWidth) {
		return 0, 0, io.EOF
	}

	// wtf, he should have done a better job explaining this
	// ok, index works by storing two reference values
	// <uint32RecordOffset> <uint64StoreFilePosition>
	// fetch both
	out = enc.Uint32(i.mmap[pos : pos+offWidth])
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entryWidth])

	return out, pos, nil
}

func (i *index) Write(off uint32, pos uint64) error {
	if uint64(len(i.mmap)) < (i.size + entryWidth) {
		return io.EOF
	}

	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)
	enc.PutUint64(i.mmap[i.size+offWidth:i.size+entryWidth], pos)

	i.size += entryWidth
	return nil
}

func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}

	if err := i.file.Sync(); err != nil {
		return err
	}

	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}

	return i.file.Close()
}

func (i *index) Name() string {
	return i.file.Name()
}
