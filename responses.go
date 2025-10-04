package btspkr

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
)

type RawMessage []byte

func (msg RawMessage) Bytes() []byte {
	return msg
}

func (msg RawMessage) MarshalBinary() ([]byte, error) {
	return msg, nil
}

func (msg *RawMessage) UnmarshalBinary(b []byte) error {
	*msg = make([]byte, len(b))
	copy(*msg, b)
	return nil
}

type DeviceName struct {
	name string
}

func (msg *DeviceName) String() string {
	return msg.name
}

func (msg *DeviceName) UnmarshalBinary(b []byte) error {
	if len(b) <= 2 || !bytes.HasPrefix(b, []byte("MM")) {
		return fmt.Errorf("not a valid device name: %+q", b)
	}

	msg.name = string(b[2:])

	return nil
}

type PinCode struct {
	code string
}

func (msg *PinCode) String() string {
	return msg.code
}

func (msg *PinCode) UnmarshalBinary(b []byte) error {
	if len(b) != 6 || !bytes.HasPrefix(b, []byte("MN")) {
		return fmt.Errorf("not a valid PIN code: %+q: %w", b, ErrMessageInvalid)
	}

	msg.code = string(b[2:])

	return nil
}

type DeviceStatus struct {
	connectionState ConnectionState
	playbackState   PlaybackState
	pairingState    PairingState
	muteState       MuteState
}

func (msg *DeviceStatus) ConnectionState() ConnectionState {
	return msg.connectionState
}

func (msg *DeviceStatus) PlaybackState() PlaybackState {
	return msg.playbackState
}

func (msg *DeviceStatus) PairingState() PairingState {
	return msg.pairingState
}

func (msg *DeviceStatus) MuteState() MuteState {
	return msg.muteState
}

var statusResponseRegexp = regexp.MustCompile(`^STC([0-7])N([0-1])S([0-3])M([0-1])$`)

func (msg *DeviceStatus) UnmarshalBinary(b []byte) error {
	match := statusResponseRegexp.FindStringSubmatch(string(b))

	if match == nil {
		return fmt.Errorf("not a valid status message: %+q: %w", b, ErrMessageInvalid)
	}

	connectionState, _ := strconv.ParseUint(match[1], 10, 8)
	playbackState, _ := strconv.ParseUint(match[2], 10, 8)
	pairingState, _ := strconv.ParseUint(match[3], 10, 8)
	muteState, _ := strconv.ParseUint(match[4], 10, 8)

	msg.connectionState = ConnectionState(connectionState)
	msg.playbackState = PlaybackState(playbackState)
	msg.pairingState = PairingState(pairingState)
	msg.muteState = MuteState(muteState)

	return nil
}

type MemorySegment struct {
	addr uint16
	data [8]byte
}

func (msg *MemorySegment) Addr() uint16 {
	return msg.addr
}

func (msg *MemorySegment) Bytes() []byte {
	return msg.data[:]
}

func (msg *MemorySegment) UnmarshalBinary(b []byte) error {
	if len(b) != 10 || !bytes.HasPrefix(b, []byte("E")) {
		return fmt.Errorf("not a valid memory segment: [% X]: %w", b, ErrMessageInvalid)
	}

	var (
		addr = uint8(b[1])
		data = []byte(b[2:10])
	)

	// Store the returned lower byte of the address.
	msg.addr = uint16(addr)

	// Copy what the device gave us.
	copy(msg.data[:], data)

	return nil
}
