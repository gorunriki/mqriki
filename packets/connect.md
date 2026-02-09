# CONNECT Packet Structure

Fixed Header:
  Byte 1: 0001 0000 (CONNECT = 1, Flags = 0)
  Byte 2+: Remaining Length (variable)

Variable Header:
  Protocol Name: "MQTT" (length + string)
  Protocol Level: 4 (0x04 for MQTT 3.1.1)
  Connect Flags: 1 byte
    Bit 0: Username flag
    Bit 1: Password flag
    Bit 2: Will retain
    Bit 3-4: Will QoS (2 bits)
    Bit 5: Will flag
    Bit 6: Clean Session
    Bit 7: Reserved (0)
  Keep Alive: 2 bytes (MSB, LSB)

Payload:
  Client Identifier (string)
  Will Topic (if Will flag = 1)
  Will Message (if Will flag = 1)
  Username (if Username flag = 1)
  Password (if Password flag = 1)
