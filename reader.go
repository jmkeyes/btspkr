package btspkr

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"time"
)

type MessageReader struct {
	r       io.Reader
	timeout time.Duration

	messageChan chan []byte
	statusChan  chan []byte
	errorChan   chan error
}

const maxPendingStatusUpdates = 10

func NewReader(r io.Reader, timeout time.Duration) (*MessageReader, error) {
	reader := &MessageReader{
		r:       r,
		timeout: timeout,

		messageChan: make(chan []byte),
		statusChan:  make(chan []byte, maxPendingStatusUpdates),
		errorChan:   make(chan error),
	}

	if err := reader.Reset(); err != nil {
		return nil, err
	}

	return reader, nil
}

func (r MessageReader) Reset() error {
	go r.pollForMessages()
	return nil
}

// A message full of NUL (ie: \x00) bytes is a device reset.
func isResetMessage(b []byte) bool {
	return len(b) == bytes.Count(b, []byte("\x00"))
}

// A status message is exactly 10 bytes and begins with "ST".
func isStatusMessage(b []byte) bool {
	return len(b) == 10 && bytes.HasPrefix(b, []byte("ST"))
}

// An error message is exactly the string "error".
func isErrorMessage(b []byte) bool {
	return len(b) == 5 && bytes.Equal(b, []byte("error"))
}

// An "unknown" command is exactly the string "Unknown".
func isUnknownCommand(b []byte) bool {
	return len(b) == 7 && bytes.Equal(b, []byte("Unknown"))
}

func (r MessageReader) pollForMessages() {
	defer func() {
		close(r.messageChan)
		close(r.statusChan)
		close(r.errorChan)
	}()

	scanner := bufio.NewScanner(r.r)
	scanner.Split(MessageSplitter())

	for scanner.Scan() {
		msg := []byte(scanner.Text())

		switch {
		// If 'msg' was just NUL bytes then the device was reset.
		case isResetMessage(msg):
			r.errorChan <- ErrDeviceReset
		case isErrorMessage(msg):
			r.errorChan <- ErrCommandError
		case isUnknownCommand(msg):
			r.errorChan <- ErrCommandUnknown
		// Device status updates are spurious; handle separately.
		case isStatusMessage(msg):
			r.statusChan <- msg
		// Otherwise this was a standard message.
		default:
			r.messageChan <- msg
		}
	}

	// Handle scanning errors.
	if err := scanner.Err(); err != nil {
		r.errorChan <- err
	}
}

func (r *MessageReader) ReadMessage(m Message) error {
	msg, ok := m.(Response)

	if !ok {
		return fmt.Errorf("reader: cannot unmarshal message: %+v", msg)
	}

	// Default to the generic response channel.
	source := r.messageChan

	// Status updates are on a specific channel.
	if _, ok := msg.(*DeviceStatus); ok {
		source = r.statusChan
	}

	select {
	// Try to read a message from the source.
	case data, ok := <-source:
		if ok {
			return msg.UnmarshalBinary(data)
		} else {
			r.messageChan = nil
			r.statusChan = nil
			return ErrReaderClosed
		}
	// Try to read any errors that occur.
	case err, ok := <-r.errorChan:
		if ok {
			return err
		} else {
			r.errorChan = nil
			return ErrReaderClosed
		}
	// Otherwise timeout the call.
	case <-time.After(r.timeout):
		return ErrTimeout
	}
}
