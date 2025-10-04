package main

import (
	"log"
	"os"
	"time"

	"go.bug.st/serial"
	btspkr "jmkeyes.ca/btspkr"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("%s /dev/cu.usbserial-XXXXXXXX output.bin", os.Args[0])
	}

	// First argument is the serial device character device.
	serialDevice := os.Args[1]

	// Open the serial port at 9600 baud 8N1.
	port, err := serial.Open(serialDevice, &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	})

	if err != nil {
		log.Fatalf("could not open serial port: %+v", err)
	}

	defer func() {
		if err := port.Close(); err != nil {
			log.Fatalf("Unable to close serial port: %+v", err)
		}
	}()

	// Communicate over the serial device port.
	device, err := btspkr.NewDevice(port, 1*time.Second)

	if err != nil {
		log.Fatalf("could not communicate with device: %+v", err)
	}

	defer func() {
		if err := device.Close(); err != nil {
			log.Fatalf("Unable to close device: %+v", err)
		}
	}()

	// Save the output here.
	outputFileName := os.Args[2]
	outputFile, err := os.Create(outputFileName)

	if err != nil {
		log.Fatalf("Unable to dump eeprom: %+v", err)
	}

	defer func() {
		if err := outputFile.Close(); err != nil {
			log.Fatalf("could not close output file!")
		}
	}()

	// Copy the device EEPROM to the output file.
	n, err := io.Copy(outputFile, device.Memory())

	if err != nil {
		log.Printf("Got error when copying EEPROM: %+v", err)
	} else {
		log.Printf("Wrote %d bytes to output file.", n)
	}

	// Flush the file contents to disk.
	outputFile.Sync()
}
