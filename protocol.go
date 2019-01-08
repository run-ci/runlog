package runlog

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	// ErrShortPacket is returned whenever a short packet is being
	// encoded or decoded.
	ErrShortPacket = errors.New("got payload smaller than 255 bytes")
	// ErrCorruptData is returned whenever there's a mismatch between
	// what the packet thinks the length should be and what the payload
	// length actually is.
	ErrCorruptData = errors.New("payload length and packet length don't match")
)

// Packet represents a small chunk of log data.
type Packet struct {
	TaskID uint32
	// ByteLength can't be too big, because if it is the
	// streaming becomes too choppy. Messages are always
	// assumed to take up the maximum byte length. If
	// not it is assumed that the stream is closed.
	ByteLength uint8
	Payload    []byte
}

// Decode deserializes a byte buffer into a Packet. If the
// Packet's ByteLength is less than the maximum size allowed
// ErrShortPacket is returned.
func (p *Packet) Decode(buf []byte) error {
	p.TaskID = binary.BigEndian.Uint32(buf[0:4])
	p.ByteLength = buf[4]
	// Adding 5 to the biggest possible uint8 results in an overflow so a conversion
	// to int is necessary first.
	p.Payload = buf[5 : int(p.ByteLength)+5]

	if p.ByteLength < 0xFF {
		return ErrShortPacket
	}

	// If the packet's real payload length isn't actually
	// what the packet says it should be, something must
	// have gone wrong.
	if int(p.ByteLength) != len(p.Payload) {
		fmt.Printf("bytelength: %v payload length: %v\n", int(p.ByteLength), len(p.Payload))
		return ErrCorruptData
	}

	return nil
}
