package btspkr

import (
	"fmt"
	"io"
	"regexp"
	"time"
)

const DefaultTimeout = 250 * time.Millisecond

// Device combines a MessageReader, MessageWriter and an EEPROM helper.
type Device struct {
	r   *MessageReader
	w   *MessageWriter
	mem *EEPROM
}

func NewDevice(rwc io.ReadWriter, timeout time.Duration) (*Device, error) {
	reader, err := NewReader(rwc, timeout)

	if err != nil {
		return nil, fmt.Errorf("cannot init reader: %w", err)
	}

	writer, err := NewWriter(rwc, timeout)

	if err != nil {
		return nil, fmt.Errorf("cannot init writer: %w", err)
	}

	device := &Device{
		r: reader,
		w: writer,
	}

	// Allow io operations on the EEPROM.
	device.mem = &EEPROM{device: device}

	return device, nil
}

func (device *Device) Close() error {
	if err := device.mem.Close(); err != nil {
		return fmt.Errorf("unable to close EEPROM: %w", err)
	}

	return nil
}

func (device *Device) xmit(req Request, resp Response) error {
	if err := device.w.WriteMessage(req); err != nil {
		return fmt.Errorf("cannot write command: %w", err)
	}

	if err := device.r.ReadMessage(resp); err != nil {
		return fmt.Errorf("cannot read response: %w", err)
	}

	return nil
}

func (device *Device) Status() (*DeviceStatus, error) {
	var (
		req  = GetDeviceStatus()
		resp DeviceStatus
	)

	if err := device.xmit(req, &resp); err != nil {
		return nil, fmt.Errorf("could not get device status: %w", err)
	}

	return &resp, nil
}

func (device *Device) DeviceName() (string, error) {
	var (
		req  = GetDeviceName()
		resp DeviceName
	)

	if err := device.xmit(req, &resp); err != nil {
		return "", fmt.Errorf("cannot get device name: %w", err)
	}

	return resp.String(), nil
}

func (device *Device) SetDeviceName(name string) error {
	var (
		req  = SetDeviceName(name)
		resp DeviceName
	)

	if len(name) > 30 {
		return fmt.Errorf("device name length exceeded")
	}

	if err := device.xmit(req, &resp); err != nil {
		return fmt.Errorf("cannot set device name: %w", err)
	}

	// The returned name should match what we set.
	if name != resp.String() {
		return fmt.Errorf("device name inconsistency")
	}

	return nil
}

func (device *Device) PinCode() (string, error) {
	var (
		req  = GetPinCode()
		resp PinCode
	)

	if err := device.xmit(req, &resp); err != nil {
		return "", fmt.Errorf("cannot get pin code: %w", err)
	}

	return resp.String(), nil
}

func (device *Device) SetPinCode(code string) error {
	var (
		req  = SetPinCode(code)
		resp PinCode
	)

	// A PIN code should be exactly four integer digits.
	if ok, _ := regexp.MatchString(`^\d{4}$`, code); !ok {
		return fmt.Errorf("not a pin code: %+v", code)
	}

	if err := device.xmit(req, &resp); err != nil {
		return fmt.Errorf("cannot set pin code: %w", err)
	}

	// The returned pin code should match what we set.
	if code != resp.String() {
		return fmt.Errorf("pin code inconsistency")
	}

	return nil
}

func (device *Device) ReadMemory(addr uint16) ([]byte, error) {
	var (
		req  = ReadMemory(addr)
		resp MemorySegment
	)

	if err := device.xmit(req, &resp); err != nil {
		return nil, fmt.Errorf("cannot read memory: %w", err)
	}

	if (addr & 0xff) != resp.Addr() {
		return nil, fmt.Errorf("did not match %04x != %04x", addr, resp.Addr())
	}

	return resp.Bytes(), nil
}

func (device *Device) WriteMemory(addr uint16, b []byte) ([]byte, error) {
	var (
		req  = WriteMemory(addr, b)
		resp MemorySegment
	)

	if err := device.xmit(req, &resp); err != nil {
		return nil, fmt.Errorf("cannot write memory: %w", err)
	}

	if (addr & 0xff) != resp.Addr() {
		return nil, fmt.Errorf("address did not match %04x != %04x", addr, resp.Addr())
	}

	return resp.Bytes(), nil
}

func (device *Device) Memory() *EEPROM {
	return device.mem
}
