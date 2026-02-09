package packets

type ConnectPacket struct {
	ProtocolName    string
	ProtocolVersion byte
	CleanSession    bool
	KeepAlive       uint16
	ClientID        string
}

func EncodeConnect(packet *ConnectPacket) []byte {
	result := []byte{0x10}

	// init variableHeader
	variableHeader := []byte{}

	// add protocol name
	variableHeader = append(variableHeader, 0, 4)
	variableHeader = append(variableHeader, []byte("MQTT")...)

	// add protocol version
	variableHeader = append(variableHeader, 4)

	// init connect flags, need to change if starts using username, password, etc..
	connectFlags := byte(0)
	if packet.CleanSession {
		connectFlags |= 0x02
	}

	// add connect flag
	variableHeader = append(variableHeader, connectFlags)

	// add keep alive (2 bytes, big endian)
	variableHeader = append(variableHeader, byte(packet.KeepAlive>>8), byte(packet.KeepAlive&0xFF))

	// add pyaload
	clientID := packet.ClientID
	variableHeader = append(variableHeader, byte(len(clientID)>>8), byte(len(clientID)&0xFF))
	variableHeader = append(variableHeader, []byte(clientID)...)

	// add remail length to fixed header
	remainingLength := len(variableHeader)
	result = append(result, encodeLength(remainingLength)...)

	// combine fixed header with variable header
	result = append(result, variableHeader...)

	return result

}

// untuk encode variable lenth dari integer
func encodeLength(length int) []byte {
	var result []byte
	for {
		digit := byte(length % 128)
		length = length / 128
		if length > 0 {
			digit = digit | 0x80
		}
		result = append(result, digit)
		if length <= 0 {
			break
		}
	}
	return result
}
