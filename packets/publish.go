package packets

import "fmt"

type PublishPacket struct {
	Dup      bool
	QoS      byte
	Retain   bool
	Topic    string
	PacketID uint16
	Payload  []byte
}

func EncodePublish(packet *PublishPacket) []byte {
	result := []byte{0x30}

	// set flags
	flags := byte(0)
	if packet.Dup {
		flags |= 0x08 // OR to switch on the DUP flag (bit 3)
	}
	flags |= (packet.QoS & 0x03) << 1 // AND to get first 2 bits, then shift left by 1, then OR to set QoS bits (bits 1 and 2)
	if packet.Retain {
		flags |= 0x01 // OR to switch on the retain flag (bit 0)
	}
	result[0] |= flags // OR the flags into the first byte to set first byte of fixed header

	// set variable header
	variableHeader := []byte{}

	// add topic
	variableHeader = append(variableHeader, byte(len(packet.Topic)>>8), byte(len(packet.Topic)&0xFF)) // topic length (2 bytes, big endian)
	variableHeader = append(variableHeader, []byte(packet.Topic)...)                                  // append topic to variable header

	// add packet ID if QoS > 0
	if packet.QoS > 0 {
		variableHeader = append(variableHeader, byte(packet.PacketID>>8), byte(packet.PacketID&0xFF)) // packet ID (2 bytes, big endian)
	}

	// add payload
	payload := packet.Payload

	// calculate remaining length
	remainingLength := len(variableHeader) + len(payload)
	result = append(result, encodeLength(remainingLength)...)
	result = append(result, variableHeader...)
	result = append(result, payload...)

	return result
}

// DecodePublish parses a full PUBLISH packet (starting from fixed header)
// `data` must contain the fixed header byte(s) and the remaining length and remaining bytes.
func DecodePublish(data []byte) (*PublishPacket, error) {
	fmt.Println(string(data)) // view data content for debugging

	// need at least the first fixed header byte + one remaining-length byte
	// check the packet length
	if len(data) < 2 {
		return nil, fmt.Errorf("packet too short")
	}

	// ambil byte pertama dan flagsnya
	// Fixed header first byte contains packet type (high nibble) and flags (low nibble).
	// We only need the low nibble for PUBLISH: DUP (bit3), QoS (bits1-2), RETAIN (bit0).
	b0 := data[0]
	flags := b0 & 0x0F
	dup := (flags & 0x08) != 0
	qos := (flags >> 1) & 0x03
	retain := (flags & 0x01) != 0

	// Decode Remaining Length (MQTT variable-length encoding). This starts at data[1].
	// Each byte contributes 7 bits; MSB set indicates continuation.
	idx := 1
	multiplier := 1
	value := 0
	for {
		if idx >= len(data) {
			return nil, fmt.Errorf("malformed remaining length")
		}
		digit := int(data[idx])
		idx++
		value += (digit & 0x7F) * multiplier
		// if MSB==0 this is the last byte
		if (digit & 0x80) == 0 {
			break
		}
		multiplier *= 128
		// safety: remaining length uses at most 4 bytes
		if multiplier > 128*128*128 {
			return nil, fmt.Errorf("malformed remaining length")
		}
	}

	// `value` is the Remaining Length: number of bytes for variable header + payload
	remaining := value
	if len(data)-idx < remaining {
		return nil, fmt.Errorf("incomplete packet: need %d bytes, have %d", remaining, len(data)-idx)
	}

	// Slice out the variable header + payload for easier parsing
	buf := data[idx : idx+remaining]
	// variable header must contain at least topic length (2 bytes)
	if len(buf) < 2 {
		return nil, fmt.Errorf("malformed variable header")
	}

	// Topic is encoded as 2-byte length (big-endian) followed by that many bytes
	topicLen := int(buf[0])<<8 | int(buf[1])
	if 2+topicLen > len(buf) {
		return nil, fmt.Errorf("malformed topic length")
	}
	topic := string(buf[2 : 2+topicLen])
	pos := 2 + topicLen

	// If QoS > 0, a 2-byte Packet Identifier follows the topic
	var packetID uint16
	if qos > 0 {
		if pos+2 > len(buf) {
			return nil, fmt.Errorf("missing packet identifier")
		}
		packetID = uint16(buf[pos])<<8 | uint16(buf[pos+1])
		pos += 2
	}

	// The remainder of buf is the application payload
	payload := make([]byte, len(buf)-pos)
	copy(payload, buf[pos:])

	// Build and return the parsed PublishPacket
	pkt := &PublishPacket{
		Dup:      dup,
		QoS:      qos,
		Retain:   retain,
		Topic:    topic,
		PacketID: packetID,
		Payload:  payload,
	}
	return pkt, nil
}
