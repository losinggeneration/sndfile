package sndfile

/*
#cgo linux    LDFLAGS: -lsndfile

#include <stdlib.h>
#include <sndfile.h>
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"io"
	"runtime"
	"unsafe"
)

// SndFile is the main type for this library.
type SndFile struct {
	// Info provides needed information about an opened file
	Info Info

	sndfile *C.SNDFILE
}

// Mode is used to describe how files should be opened (read, write, or both)
type Mode C.int

const (
	// ReadMode specifies a file will be opened read-only
	ReadMode = Mode(C.SFM_READ)
	// WriteMode specifies a file will be opened write-only
	WriteMode = Mode(C.SFM_WRITE)
	// ReadWriteMode specifies a file will be opened and able to read & write to the file
	ReadWriteMode = Mode(C.SFM_RDWR)
)

func parseError() error {
	e := C.sf_strerror(nil)

	return fmt.Errorf(C.GoString(e))
}

// Open will open a file in ReadMode.
func Open(path string) (*SndFile, error) {
	var s SndFile

	err := s.Open(path, ReadMode)
	if err != nil {
		return nil, err
	}

	return &s, nil
}

// Open will open the specified path with a given mode. SndFile.Info should be filled
// in prior to making this call to open the file with the corerct format information.
func (s *SndFile) Open(path string, mode Mode) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	s.Info.fillCInfo()

	// ReadMode requires format to be 0, make sure it is
	if mode == ReadMode {
		s.Info.info.format = 0
	}

	s.sndfile = C.sf_open(cpath, C.int(mode), &s.Info.info)
	if s.sndfile == nil {
		return parseError()
	}
	s.Info.fillInfo()

	runtime.SetFinalizer(s, func(s *SndFile) {
		s.Close()
	})

	return nil
}

// Close will close the SndFile. It's a good idea to call this even though Open sets
// a finalizer on s to call Close when s is garbage collected.
func (s *SndFile) Close() {
	if s.sndfile != nil {
		C.sf_close(s.sndfile)
		s.sndfile = nil
	}
}

// ReadFrames will read an amount frames from a source and put it into an int16 slice
func (s *SndFile) ReadFrames(frames uint) ([]int16, error) {
	if frames > s.Info.Frames {
		frames = s.Info.Frames
	}

	p := make([]int16, frames*uint(s.Info.Channels))

	n := C.sf_readf_short(s.sndfile, (*C.short)(unsafe.Pointer(&p[0])), C.sf_count_t(frames))
	if n == 0 {
		return nil, io.EOF
	}

	return p, nil
}

// ReadFrames32 will read an amount frames from a source and put it into an int32 slice
func (s *SndFile) ReadFrames32(frames uint) ([]int32, error) {
	if frames > s.Info.Frames {
		frames = s.Info.Frames
	}

	p := make([]int32, frames*uint(s.Info.Channels))

	n := C.sf_readf_int(s.sndfile, (*C.int)(unsafe.Pointer(&p[0])), C.sf_count_t(frames))
	if n == 0 {
		return nil, io.EOF
	}

	return p, nil
}

// Int16ToByte is a little endian conversion of int16 slice into a byte slice
func Int16ToByte(p []int16) []byte {
	b := make([]byte, len(p)*2)
	for i := 0; i < len(p); i++ {
		o := i * 2
		binary.LittleEndian.PutUint16(b[o:o+2], uint16(p[i]))
	}

	return b
}

// Int16ToByteBe is a big endian conversion of int16 slice into a byte slice
func Int16ToByteBe(p []int16) []byte {
	b := make([]byte, len(p)*2)
	for i := 0; i < len(p); i++ {
		o := i * 2
		binary.BigEndian.PutUint16(b[o:o+2], uint16(p[i]))
	}

	return b
}
