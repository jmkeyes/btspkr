package btspkr

import (
	"encoding"
	"fmt"
	"log"
)

// Command is a wrapper around a byte slice.
type Command []byte

type commandFunc func(args ...any) Command

var _ encoding.BinaryMarshaler = Command{}

// NewCommand creates a Command object from a format string.
func NewCommand(format string) commandFunc {
	return func(args ...any) Command {
		// Assemble the command into a byte slice.
		data := fmt.Appendf(nil, format, args...)

		// Make it useful.
		return Command(data)
	}
}

// String returns a string representation of the Command.
func (cmd Command) String() string {
	return string(cmd)
}

// Bytes returns a byte slice representation of the Command.
func (cmd Command) Bytes() []byte {
	return cmd
}

// MarshalBinary returns the Command as ready for transmission.
func (cmd Command) MarshalBinary() ([]byte, error) {
	return cmd.Bytes(), nil
}

// The vast majority of commands work with simple format strings.
const (
	cmdGetDeviceStatus       = "AT#CY"
	cmdDecreaseVolume        = "AT#VD"
	cmdIncreaseVolume        = "AT#VU"
	cmdGetVolumeLevel        = "AT#VV"
	cmdSetVolumeLevel        = "AT#VV%02d"
	cmdPlayPauseMusic        = "AT#MA"
	cmdStopMusic             = "AT#MC"
	cmdForwardMusic          = "AT#MD"
	cmdBackwardMusic         = "AT#ME"
	cmdReset                 = "AT#CC"
	cmdDisconnect            = "AT#MI"
	cmdGetDeviceName         = "AT#MM"
	cmdSetDeviceName         = "AT#MM%30s"
	cmdGetPinCode            = "AT#MN"
	cmdSetPinCode            = "AT#MN%4s"
	cmdEnablePairing         = "AT#CA"
	cmdDisablePairing        = "AT#CB"
	cmdSetMutePinInputState  = "AT#CZ"
	cmdSetMutePinOutputState = "AT#VF"
	cmdMemoryOperationPrefix = "AT#MS"
)

var (
	// Volume control Commands.
	DecreaseVolume = NewCommand(cmdDecreaseVolume)
	IncreaseVolume = NewCommand(cmdIncreaseVolume)
	GetVolumeLevel = NewCommand(cmdGetVolumeLevel)
	SetVolumeLevel = NewCommand(cmdSetVolumeLevel)

	// Music control Commands.
	PlayPauseMusic = NewCommand(cmdPlayPauseMusic)
	StopMusic      = NewCommand(cmdStopMusic)
	ForwardMusic   = NewCommand(cmdForwardMusic)
	BackwardMusic  = NewCommand(cmdBackwardMusic)

	// Bluetooth control Commands.
	Reset            = NewCommand(cmdReset)
	Disconnect       = NewCommand(cmdDisconnect)
	GetDeviceName    = NewCommand(cmdGetDeviceName)
	SetDeviceName    = NewCommand(cmdSetDeviceName)
	GetPinCode       = NewCommand(cmdGetPinCode)
	SetPinCode       = NewCommand(cmdSetPinCode)
	EnablePairing    = NewCommand(cmdEnablePairing)
	DisablePairing   = NewCommand(cmdDisablePairing)
	SetMutePinInput  = NewCommand(cmdSetMutePinInputState)
	SetMutePinOutput = NewCommand(cmdSetMutePinOutputState)
	GetDeviceStatus  = NewCommand(cmdGetDeviceStatus)

	// EEPROM read/write Commands.
	ReadMemory = func(addr uint16) Command {
		// A read Command always sets the top bit.
		var value uint16 = 1 << 15 // 0x8000

		// Combine with the 12 bit address.
		value |= uint16(addr & 0xfff)

		// Extract the high and low bytes.
		var (
			high = byte(value >> 8)
			low  = byte(value & 0xff)
		)

		// Assemble the "read memory" Command.
		data := []byte(cmdMemoryOperationPrefix)
		data = append(data, high, low)

		return data
	}

	WriteMemory = func(addr uint16, buf []byte) Command {
		// A write Command sets the top two bits.
		var value uint16 = 3 << 13 // 0x6000

		// Combine with the 12 bit address.
		value |= uint16(addr & 0xfff)

		// Extract the high and low bytes.
		var (
			high = byte(value >> 8)
			low  = byte(value & 0xff)
		)

		// Assemble the "write memory" Command.
		data := []byte(cmdMemoryOperationPrefix)
		data = append(data, high, low)

		// Append the data to write to the EEPROM.
		data = append(data, buf...)

		log.Printf("DEBUG: write memory Command: len=%d cmd=[% X]", len(data), data)

		return []byte("")
	}
)
