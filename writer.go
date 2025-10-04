package btspkr

import (
	"bytes"
	"fmt"
	"io"
	"time"
)

type MessageWriter struct {
	w        io.Writer
	timeout  time.Duration
	cmdCount int
}

func NewWriter(w io.Writer, timeout time.Duration) (*MessageWriter, error) {
	return &MessageWriter{w, timeout, 0}, nil
}

func (w *MessageWriter) sendWakeCommand() error {
	// Write a series of 10x 0xFF bytes to "wake up" the device.
	wakeDevice := bytes.Repeat([]byte("\xff"), 10)
	wakeDevice = append(wakeDevice, CRLF...)

	if _, err := w.w.Write(wakeDevice); err != nil {
		return fmt.Errorf("cannot write wake command: %w", err)
	}

	return nil
}

const wakeDeviceInterval = 10

func (w *MessageWriter) WriteMessage(m Message) error {
	msg, ok := m.(Request)

	if !ok {
		return fmt.Errorf("cannot marshal message: %+v", m)
	}

	// Send a wake command before every 10th message.
	if (w.cmdCount % wakeDeviceInterval) == 0 {
		if err := w.sendWakeCommand(); err != nil {
			return err
		}
	}

	// Encode the message into a byte slice.
	data, err := msg.MarshalBinary()

	if err != nil {
		return fmt.Errorf("cannot encode message: %w", err)
	}

	// Append a delimiter to the end of the message.
	data = append(data, CRLF...)

	if _, err := w.w.Write(data); err != nil {
		return fmt.Errorf("cannot write message: %w", err)
	}

	w.cmdCount++

	return nil
}
