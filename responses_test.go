package btspkr_test

import (
	"bytes"
	"errors"
	"testing"

	"jmkeyes.ca/btspkr"
)

func TestRawMessage(t *testing.T) {
	var (
		data = []byte("MESSAGE")
		msg  btspkr.RawMessage
	)

	// Unmarshal the message from the byte slice.
	if err := msg.UnmarshalBinary(data); err != nil {
		t.Errorf("want nil got %+v", err)
	}

	// Marshal the message back into a byte slice.
	if buf, err := msg.MarshalBinary(); err != nil {
		t.Errorf("want nil got %+v", err)
	} else {
		if !bytes.Equal(buf, data) {
			t.Errorf("want %+v got %+v", data, buf)
		}
	}
}

func TestDeviceName(t *testing.T) {
	tests := []struct {
		test  string
		given []byte
		name  string
		err   error
	}{
		{
			"Normal",
			[]byte("MMLANDO"),
			"LANDO",
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			var m btspkr.DeviceName

			if err := m.UnmarshalBinary(tt.given); !errors.Is(err, tt.err) {
				t.Errorf("want %+v got %+v", tt.err, err)
			} else {
				if tt.name != m.String() {
					t.Errorf("want %+v got %+v", tt.name, m.String())
				}
			}
		})
	}
}

func TestPinCode(t *testing.T) {
	tests := []struct {
		test  string
		given []byte
		code  string
		err   error
	}{
		{
			"Normal",
			[]byte("MN0000"),
			"0000",
			nil,
		},
		{
			"Truncated",
			[]byte("MN0"),
			"",
			btspkr.ErrMessageInvalid,
		},
		{
			"Extended",
			[]byte("MN00000"),
			"",
			btspkr.ErrMessageInvalid,
		},
		{
			"NonNumeric",
			[]byte("MNDEADBEEF"),
			"",
			btspkr.ErrMessageInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			var m btspkr.PinCode

			if err := m.UnmarshalBinary(tt.given); !errors.Is(err, tt.err) {
				t.Errorf("want %+v got %+v", tt.err, err)
			} else {
				if tt.code != m.String() {
					t.Errorf("want %+v got %+v", tt.code, m.String())
				}
			}
		})
	}
}

func TestDeviceStatus(t *testing.T) {
	tests := []struct {
		test       string
		given      []byte
		connection btspkr.ConnectionState
		playback   btspkr.PlaybackState
		pairing    btspkr.PairingState
		mute       btspkr.MuteState
		err        error
	}{
		{
			"Normal",
			[]byte("STC0N0S0M0"),
			btspkr.NotConnected,
			btspkr.Playing,
			btspkr.Bound,
			btspkr.Muted,
			nil,
		},
		{
			"Invalid",
			[]byte("\xff\xff\xff\xff\xff"),
			0, 0, 0, 0,
			btspkr.ErrMessageInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			var m btspkr.DeviceStatus

			if err := m.UnmarshalBinary(tt.given); !errors.Is(err, tt.err) {
				t.Errorf("want %+v got %+v", tt.err, err)
			} else {
				if tt.connection != m.ConnectionState() {
					t.Errorf("want %+v got %+v", tt.connection, m.ConnectionState())
				}
				if tt.playback != m.PlaybackState() {
					t.Errorf("want %+v got %+v", tt.playback, m.PlaybackState())
				}
				if tt.pairing != m.PairingState() {
					t.Errorf("want %+v got %+v", tt.pairing, m.PairingState())
				}
				if tt.mute != m.MuteState() {
					t.Errorf("want %+v got %+v", tt.mute, m.MuteState())
				}
			}
		})
	}
}

func TestMemorySegment(t *testing.T) {
	tests := []struct {
		test  string
		given []byte
		addr  uint16
		bytes [8]byte
		err   error
	}{
		{
			"Normal",
			[]byte("E\xfe\xff\xff\xff\xff\xff\xff\xff\xff"),
			0xfe,
			[8]byte{
				0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff,
			},
			nil,
		},
		{
			"Empty",
			[]byte{},
			0x00,
			[8]byte{},
			btspkr.ErrMessageInvalid,
		},
		{
			"Truncated",
			[]byte("E\xfe\xff"),
			0x00,
			[8]byte{},
			btspkr.ErrMessageInvalid,
		},
		{
			"Extended",
			[]byte("E\xfe\xff\xff\xff\xff\xff\xff\xff\xff\xff"),
			0x00,
			[8]byte{},
			btspkr.ErrMessageInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.test, func(t *testing.T) {
			var m btspkr.MemorySegment

			if err := m.UnmarshalBinary(tt.given); !errors.Is(err, tt.err) {
				t.Errorf("want %+v got %+v", tt.err, err)
			} else {
				if tt.addr != m.Addr() {
					t.Errorf("want %+v got %+v", tt.addr, m.Addr())
				}
				if !bytes.Equal(tt.bytes[:], m.Bytes()) {
					t.Errorf("want [% X] got [% X]", tt.bytes, m.Bytes())
				}
			}
		})
	}
}
