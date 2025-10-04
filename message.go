package btspkr

import (
	"bufio"
	"bytes"
	"encoding"
)

type Message interface {
	// Nothing to see here, yet.
}

type Response interface {
	encoding.BinaryUnmarshaler
}

type Request interface {
	encoding.BinaryMarshaler
}

var CRLF = []byte("\r\n")

func byteRunLength(b []byte, c byte) (i int) {
	if len(b) == 0 || b[0] != c {
		return -1
	}

	for i = 1; i < len(b) && b[i] == c; i++ {
		// Advance until the buffer is consumed.
	}

	return i
}

func MessageSplitter() bufio.SplitFunc {
	return func(frame []byte, atEOF bool) (int, []byte, error) {
		// At minimum we need a CRLF sequence.
		if len(frame) <= 2 {
			return 0, nil, nil
		}

		// Return what we have if there's data in the buffer at EOF.
		if atEOF && len(frame) >= 0 {
			return len(frame), frame, bufio.ErrFinalToken
		}

		var (
			frameStart  = 0
			frameLength = len(frame)
			frameEnd    = frameLength
		)

		// Some messages are initialization bytes indicating the device has reset.
		if resetIndex := byteRunLength(frame[frameStart:frameLength], '\x00'); resetIndex > 0 {
			return resetIndex, frame[frameStart:resetIndex], nil
		}

		// Look for the opening CRLF to establish the message start boundary.
		if frameStart = bytes.Index(frame[frameStart:frameEnd], CRLF); frameStart == -1 {
			return 0, nil, nil
		}

		// Define the message content as everything within the framing delimiters.
		var (
			messageStart  = frameStart + 2
			messageEnd    = frameLength - 2
			messageLength = messageEnd - messageStart
		)

		// Some messages are memory segments that contain arbitrary bytes.
		// This occurs when:
		//   - The frame is exactly 14 bytes in length.
		//   - The frame is surrounded by CRLF byte sequences.
		//   - The message within the frame is exactly 10 bytes in length.
		//   - The message begins with the ASCII character 'E'.
		if frameLength == 14 && messageLength == 10 && frame[messageStart] == 'E' && bytes.HasSuffix(frame, CRLF) {
			return frameLength, frame[messageStart:messageEnd], nil
		}

		// Look for the closing CRLF to update the (potential) message end boundary.
		if messageEnd = bytes.Index(frame[messageStart:frameLength], CRLF); messageEnd == -1 {
			return 0, nil, nil
		} else {
			// Compensate for the non-zero messageStart offset.
			messageEnd += 2
		}

		return frameLength, frame[messageStart:messageEnd], nil
	}
}
