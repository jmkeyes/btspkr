package btspkr

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

type EEPROM struct {
	device *Device
	offset int64
	mutex  sync.RWMutex
}

var (
	_ io.Reader   = &EEPROM{}
	_ io.ReaderAt = &EEPROM{}

	_ io.Writer   = &EEPROM{}
	_ io.WriterAt = &EEPROM{}

	_ io.Seeker = &EEPROM{}
	_ io.Closer = &EEPROM{}
)

func (mem *EEPROM) Read(p []byte) (int, error) {
	n, err := mem.ReadAt(p, mem.offset)

	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	mem.offset += int64(n)

	return n, err
}

func (mem *EEPROM) ReadAt(p []byte, off int64) (int, error) {
	mem.mutex.RLock()
	defer mem.mutex.RUnlock()

	if off > 0xfff {
		return 0, io.EOF // consider wrapping errors here
	}

	// Read from device memory.
	buf, err := mem.device.ReadMemory(uint16(off))

	if err != nil {
		return 0, fmt.Errorf("cannot read from device: %w", err)
	}

	return copy(p, buf), nil
}

func (mem *EEPROM) Write(p []byte) (int, error) {
	n, err := mem.WriteAt(p, mem.offset)

	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	mem.offset += int64(n)

	return n, err
}

func (mem *EEPROM) WriteAt(p []byte, off int64) (int, error) {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	if off < 0x200 || off > 0xfff {
		return 0, fmt.Errorf("cannot write to offset %04x: %w", off, ErrMemoryAccess)
	}

	// Write to device memory.
	buf, err := mem.device.WriteMemory(uint16(off), p)

	if err != nil {
		return 0, fmt.Errorf("cannot write to device: %w", err)
	}

	return len(buf), nil
}

func (mem *EEPROM) Seek(off int64, whence int) (int64, error) {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	var offset int64

	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	switch whence {
	case io.SeekStart:
		offset = 0x000 + off
	case io.SeekCurrent:
		offset = mem.offset + off
	case io.SeekEnd:
		offset = 0xfff + off
	default:
		return 0, errors.New("invalid whence")
	}

	if offset < 0 || offset > 0xfff {
		return 0, errors.New("invalid offset")
	}

	// Update the stored offset.
	mem.offset = offset

	return mem.offset, nil
}

func (mem *EEPROM) Close() error {
	mem.mutex.Lock()
	defer mem.mutex.Unlock()

	mem.device = nil
	mem.offset = 0
	return nil
}
