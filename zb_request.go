package main

/*
func SendRemoteATCommandRequest(destination64, destination16, ATCommand string) (err error) {
	var frame xbee.RemoteATCommandRequestFrame
	frame.StartDelimiter = 0x7E
	frame.Length = 0x000F
	frame.Type = 0x17
	frame.ID = 0x01

	// handle destination address
	if destination64 == "" && destination16 != "" {
		destination64 = "FFFFFFFFFFFFFFFF"
	} else if destination16 == "" && destination64 != "" {
		destination16 = "FFFE"
	} else if destination16 != "" && destination64 != "" {
		destination16 = "FFFE"
	} else {
		return fmt.Errorf("no valid destination address received")
	}

	b64, err := hex.DecodeString(destination64)
	if err != nil {
		return fmt.Errorf("failed to encode 64 bit destination address to hex: %w", err)
	}
	frame.DestinationAddress64 = binary.BigEndian.Uint64(b64)

	b16, err := hex.DecodeString(destination16)
	if err != nil {
		return fmt.Errorf("failed to encode 16 bit destination address to hex: %w", err)
	}
	frame.DestinationAddress16 = binary.BigEndian.Uint16(b16)

	frame.RemoteCommandOptions = 0x00

	// handle AT command
	atb, err := hex.DecodeString(ATCommand)
	if err != nil {
		return fmt.Errorf("failed to encode 16 bit AT command to hex: %w", err)
	}
	frame.ATCommand = binary.BigEndian.Uint16(atb)

	// handle parameter value
	// frame.Length += parameter value length

	buf := new(bytes.Buffer)
	//binary.Write(buf, binary.BigEndian, frame.StartDelimiter)
	//binary.Write(buf, binary.BigEndian, frame.Length)
	binary.Write(buf, binary.BigEndian, frame.Type)
	binary.Write(buf, binary.BigEndian, frame.ID)
	binary.Write(buf, binary.BigEndian, frame.DestinationAddress64)
	binary.Write(buf, binary.BigEndian, frame.DestinationAddress16)
	binary.Write(buf, binary.BigEndian, frame.RemoteCommandOptions)
	binary.Write(buf, binary.BigEndian, frame.ATCommand)
	binary.Write(buf, binary.BigEndian, frame.ParameterValue)
	//binary.Write(buf, binary.BigEndian, frame.Checksum)

	// compute checksum
	var sum uint64
	ba := buf.Bytes()
	for _, b := range ba {
		sum += uint64(b)
	}
	lastByte := ba[len(ba)-1]
	frame.Checksum = 0xff - lastByte

	log.Debugf("frame: %+v", frame)

	return nil
}
*/
