package utils

import "encoding/json"

func Decode(data []byte, value any) {
	if data != nil {
		json.Unmarshal(data, value)
	}
}
func Encode(data any) []byte {
	b := make([]byte, 0)
	if data != nil {
		b, _ = json.Marshal(data)
	}
	return b
}

// func DecodeUint64(data []byte) (uint64, error) {
// 	m := &pb.Uint64{}
// 	if data != nil {
// 		err := proto.Unmarshal(data, m)
// 		if err != nil {
// 			return 0, err
// 		}
// 		return m.Name, nil
// 	}
// 	return 0, nil
// }
// func EncodeUint64(data uint64) ([]byte, error) {
// 	m := &pb.Uint64{
// 		Name: data,
// 	}
// 	binaryData, err := proto.Marshal(m)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return binaryData, nil
// }
