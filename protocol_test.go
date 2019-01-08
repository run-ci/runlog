package runlog

import (
	"encoding/binary"
	"strings"
	"testing"
)

func TestPacketDecode(t *testing.T) {
	type result struct {
		packet Packet
		err    error
	}

	tests := []struct {
		label string
		// Using a string here to make it easier to read what
		// log message is actually being passed in as input.
		taskid   uint32
		input    string
		expected result
		actual   result
	}{
		{
			label:  "golden",
			input:  strings.Repeat("a", 255),
			taskid: uint32(0),
			expected: result{
				packet: Packet{
					TaskID:     uint32(0),
					ByteLength: uint8(255),
					Payload:    []byte(strings.Repeat("a", 255)),
				},
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			input := getTestPayload(test.taskid, test.input)

			err := test.actual.packet.Decode(input)
			if err != test.expected.err {
				t.Fatalf("expected err %v\n\ngot: %v\n\n", test.expected.err, err)
			}

			if !payloadEquals(test.expected.packet.Payload, test.actual.packet.Payload) {
				t.Fatalf("expected payload: %v (%v)\n\ngot: %v (%v)\n\n",
					test.expected.packet.Payload,
					len(test.expected.packet.Payload),
					test.actual.packet.Payload,
					len(test.actual.packet.Payload),
				)
			}
		})
	}
}

func getTestPayload(taskid uint32, str string) []byte {
	if len(str) > 255 {
		panic("cannot use string larger than 255 for test payload")
	}

	data := make([]byte, 4+1+len(str))

	binary.BigEndian.PutUint32(data, taskid)
	data[4] = uint8(len(str))
	copy(data[5:], []byte(str))

	return data
}

func payloadEquals(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
