package packets

type SubscribePacket struct {
	PacketID uint16
	Topics   []Subscription
}

type Subscription struct {
	Topic string
	QoS   byte
}

func EncodeSubscribe(packet *SubscribePacket) []byte {
	result := []byte{0x82} // 1000 0010 (SUBSCRIBE packet)

	// variable header
	variableHeader := []byte{byte(packet.PacketID >> 8), byte(packet.PacketID & 0xFF)} //packet ID (2bytes, big endian)

	// payload
	payload := []byte{}
	for _, sub := range packet.Topics {
		payload = append(payload, byte(len(sub.Topic)>>8), byte(len(sub.Topic)&0xFF)) // topic length (2 bytes, big endian)
		payload = append(payload, []byte(sub.Topic)...)                               // append topic
		payload = append(payload, sub.QoS)                                            // append QoS byte
	}

	// calculate remaining length
	remainingLength := len(variableHeader) + len(payload)
	result = append(result, encodeLength(remainingLength)...)
	result = append(result, variableHeader...)
	result = append(result, payload...)

	return result
}
