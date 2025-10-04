package btspkr

import "fmt"

type (
	ConnectionState int
	PairingState    int
	PlaybackState   int
	MuteState       int
)

const (
	NotConnected ConnectionState = 0
	Connected    ConnectionState = 7

	Bound        PairingState = 0
	Discoverable PairingState = 2
	Paired       PairingState = 3

	Playing PlaybackState = 0
	Paused  PlaybackState = 1

	Muted   MuteState = 0
	Unmuted MuteState = 1
)

func (state ConnectionState) String() string {
	switch state {
	case -1:
		return "Unknown"
	case NotConnected:
		return "NotConnected"
	case Connected:
		return "Connected"
	default:
		return fmt.Sprintf("ConnectionState(%d)", int(state))
	}
}

func (state PairingState) String() string {
	switch state {
	case -1:
		return "Unknown"
	case Bound:
		return "Bound"
	case Discoverable:
		return "Discoverable"
	case Paired:
		return "Paired"
	default:
		return fmt.Sprintf("PairingState(%d)", int(state))
	}
}

func (state PlaybackState) String() string {
	switch state {
	case -1:
		return "Unknown"
	case Playing:
		return "Playing"
	case Paused:
		return "Paused"
	default:
		return fmt.Sprintf("PlaybackState(%d)", int(state))
	}
}

func (state MuteState) String() string {
	switch state {
	case -1:
		return "Unknown"
	case Muted:
		return "Muted"
	case Unmuted:
		return "Unmuted"
	default:
		return fmt.Sprintf("MuteState(%d)", int(state))
	}
}
