package packets

import "fmt"

type SubackPacket struct {
	PacketID    uint16
	ReturnCodes []byte
}

func DecodeSuback(data []byte) (*SubackPacket, error) {
	if len(data) < 4 || data[0] != 0x90 {
		return nil, fmt.Errorf("SUBACK too short")
	}
	packetID := uint16(data[2])<<8 | uint16(data[3])
	returnCodes := data[4:]

	return &SubackPacket{
		PacketID:    packetID,
		ReturnCodes: returnCodes,
	}, nil
}
