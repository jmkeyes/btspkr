package btspkr_test

import (
	"bytes"
	"testing"

	"jmkeyes.ca/btspkr"
)

var WakeCommandBytes = bytes.Repeat([]byte("\xff"), 10)

func TestWriter(t *testing.T) {
	t.Run("EncodesBasicMessage", func(t *testing.T) {
		data := "MESSAGE"
		dataLength := len(data)

		port := bytes.Buffer{}
		w, err := btspkr.NewWriter(&port, DefaultTimeout)

		if err != nil {
			t.Fatalf("cannot init writer: %+v", err)
		}

		msg := btspkr.RawMessage(data)

		if err := w.WriteMessage(&msg); err != nil {
			t.Fatalf("cannot write message: %+v", err)
		}

		buf := port.Bytes()

		if got := buf[0:len(WakeCommandBytes)]; !bytes.Equal(got, WakeCommandBytes) {
			t.Errorf("expected: %+q got %+q", WakeCommandBytes, got)
		}

		if got := buf[12:(12 + dataLength)]; !bytes.Equal(got, []byte(data)) {
			t.Errorf("expected: %+q got: %+q", []byte(data), got)
		}

		if got := bytes.Count(buf, btspkr.CRLF); got != 2 {
			t.Errorf("expected: 2 got: %d", got)
		}
	})

	t.Run("WakeMessageCount", func(t *testing.T) {
		port := bytes.Buffer{}
		w, err := btspkr.NewWriter(&port, DefaultTimeout)

		if err != nil {
			t.Fatalf("cannot init writer: %+v", err)
		}

		msg := btspkr.RawMessage("TEST")

		for i := 0; i < 50; i++ {
			if err := w.WriteMessage(&msg); err != nil {
				t.Fatalf("cannot write message: %+v", err)
				break
			}
		}

		want := 5
		buf := port.Bytes()
		if got := bytes.Count(buf, WakeCommandBytes); got != want {
			t.Errorf("Expected %d wake messages; got %d", want, got)
		}
	})
}
