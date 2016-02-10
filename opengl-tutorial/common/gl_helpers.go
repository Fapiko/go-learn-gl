package common

import "encoding/binary"

func intFromByteSlice(bytes []byte) int {
	return int(binary.LittleEndian.Uint32(bytes))
}

func int32FromByteSlice(bytes []byte) int32 {
	return int32(binary.LittleEndian.Uint32(bytes))
}

func nullTerminatedString(source []byte) string {
	return string(append(source, make([]byte, 1)...))
}
