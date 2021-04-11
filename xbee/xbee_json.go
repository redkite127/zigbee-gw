package xbee

// type remoteATCommandRequestFrame struct {
// 	StartDelimiter       string
// 	Length               string
// 	Type                 string
// 	ID                   string
// 	DestinationAddress64 string
// 	DestinationAddress16 string
// 	RemoteCommandOptions string
// 	ATCommand            string
// 	ParameterValue       string
// 	Checksum             string
// }

// func (f RemoteATCommandRequestFrame) MarshalJSON() ([]byte, error) {
// 	var fj remoteATCommandRequestFrame
// 	log.Println("test)")
// 	fj.StartDelimiter = hex.EncodeToString([]byte{f.StartDelimiter})
// 	fj.Length = hex.EncodeToString(f.Length[:])
// 	fj.Type = hex.EncodeToString([]byte{f.Type})
// 	fj.ID = hex.EncodeToString([]byte{f.ID})
// 	fj.DestinationAddress64 = hex.EncodeToString(f.DestinationAddress64[:])
// 	fj.DestinationAddress16 = hex.EncodeToString(f.DestinationAddress16[:])
// 	fj.RemoteCommandOptions = hex.EncodeToString([]byte{f.RemoteCommandOptions})
// 	fj.ATCommand = hex.EncodeToString(f.ATCommand[:])
// 	fj.ParameterValue = hex.EncodeToString(f.ParameterValue[:])
// 	fj.Checksum = hex.EncodeToString([]byte{f.Checksum})

// 	jsonValue, err := json.Marshal(fj)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return jsonValue, nil
// }
