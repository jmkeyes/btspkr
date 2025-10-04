package btspkr_test

import (
	"bytes"
	"strings"
	"testing"

	"jmkeyes.ca/btspkr"
)

func TestReader(t *testing.T) {
	t.Run("ReadSimpleMessage", func(t *testing.T) {
		data := "\r\nMESSAGE\r\n"
		port := strings.NewReader(data)
		r, err := btspkr.NewReader(port, DefaultTimeout)

		if err != nil {
			t.Fatalf("cannot init reader: %+v", err)
		}

		msg := btspkr.RawMessage{}

		if err := r.ReadMessage(&msg); err != nil {
			t.Fatalf("cannot read message: %+v", err)
		}

		if want, got := []byte(data[2:9]), msg.Bytes(); !bytes.Equal(want, got) {
			t.Errorf("Expected raw message bytes to be [% X] but was [% X]", want, got)
		}
	})

	t.Run("ReadStatusMessage", func(t *testing.T) {
		data := "\r\nSTC0N0S0M0\r\n" // this
		port := strings.NewReader(data)
		r, err := btspkr.NewReader(port, DefaultTimeout)

		if err != nil {
			t.Fatalf("cannot init reader: %+v", err)
		}

		msg := btspkr.DeviceStatus{}

		if err := r.ReadMessage(&msg); err != nil {
			t.Fatalf("cannot read message: %+v", err)
		}

		if want, got := btspkr.ConnectionState(0), msg.ConnectionState(); want != got {
			t.Errorf("Expected connection state to be %+v but was %+v", want, got)
		}

		if want, got := btspkr.PlaybackState(0), msg.PlaybackState(); want != got {
			t.Errorf("Expected playback state to be %+v but was %+v", want, got)
		}

		if want, got := btspkr.PairingState(0), msg.PairingState(); want != got {
			t.Errorf("Expected pairing state to be %+v but was %+v", want, got)
		}
		if want, got := btspkr.MuteState(0), msg.MuteState(); want != got {
			t.Errorf("Expected mute state to be %+v but was %+v", want, got)
		}
	})

	t.Run("ReadMemorySegment", func(t *testing.T) {
		data := "\r\nEx\xff\xff\xff\xff\xff\xff\xff\xff\r\n"
		port := strings.NewReader(data)
		r, err := btspkr.NewReader(port, DefaultTimeout)

		if err != nil {
			t.Fatalf("cannot init reader: %+v", err)
		}

		msg := btspkr.MemorySegment{}

		if err := r.ReadMessage(&msg); err != nil {
			t.Fatalf("cannot read message: %+v", err)
		}

		if want, got := uint16(0x78), msg.Addr(); want != got {
			t.Errorf("Expected memory address to be %#02x but was %#02x", want, got)
		}

		if want, got := []byte(data[4:12]), msg.Bytes(); !bytes.Equal(want, got) {
			t.Errorf("Expected memory bytes to be [% X] but was [% X]", want, got)
		}
	})

	t.Run("ReadMemorySegmentWithCRLF", func(t *testing.T) {
		data := "\r\nEx\xff\xff\xff\r\n\xff\xff\xff\r\n"
		port := strings.NewReader(data)
		r, err := btspkr.NewReader(port, DefaultTimeout)

		if err != nil {
			t.Fatalf("cannot init reader: %+v", err)
		}

		msg := btspkr.MemorySegment{}

		if err := r.ReadMessage(&msg); err != nil {
			t.Fatalf("cannot read message: %+v", err)
		}

		if want, got := uint16(0x78), msg.Addr(); want != got {
			t.Errorf("Expected memory address to be %#02x but was %#02x", want, got)
		}

		if want, got := []byte(data[4:12]), msg.Bytes(); !bytes.Equal(want, got) {
			t.Errorf("Expected memory bytes to be [% X] but was [% X]", want, got)
		}
	})

	t.Run("ReadDeviceInitializationGarbage", func(t *testing.T) {
		data := bytes.Repeat([]byte("\x00"), 20)
		port := strings.NewReader(string(data))
		r, err := btspkr.NewReader(port, DefaultTimeout)

		if err != nil {
			t.Fatalf("cannot init reader: %+v", err)
		}

		msg := btspkr.RawMessage{}
		err = r.ReadMessage(&msg)

		switch {
		case err == nil:
			t.Errorf("should not read a message full of NUL bytes")
		case err != btspkr.ErrDeviceReset:
			t.Errorf("expected ErrDeviceReset but got: %+v", err)
		default:
			// Otherwise it was a successful test.
		}
	})
}
