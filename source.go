package bimg

// #include "source.h"
import "C"
import (
	"fmt"
	"github.com/h2non/bimg/pointer"
	"io"
	"sync"
	"unsafe"
)

// Source contains a libvips VipsSourceCustom and manages its lifecycle.
type Source struct {
	reader io.ReadCloser
	seeker io.Seeker
	src    *C.VipsSourceCustom
	ptr    unsafe.Pointer
	lock   sync.Mutex
}

// NewSource creates Source from reader
func NewSource(reader io.ReadCloser) *Source {
	s := &Source{reader: reader}
	seeker, ok := reader.(io.ReadSeeker)
	if ok {
		s.seeker = seeker
		s.ptr = pointer.Save(s)
		s.src = C.create_go_custom_source_with_seek(s.ptr)
	} else {
		s.ptr = pointer.Save(s)
		s.src = C.create_go_custom_source(s.ptr)
	}
	return s
}

// Close source
func (s *Source) Close() {
	s.lock.Lock()
	if s.ptr != nil {
		C.clear_source(&s.src)
		pointer.Unref(s.ptr)
		s.ptr = nil
		_ = s.reader.Close()
	}
	s.lock.Unlock()
}

// LoadImage creates a new Image from source
func (s *Source) LoadImage() (*Image, error) {
	fmt.Println("panggil")
	_, err := vipsImageFromSource(s.src)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
