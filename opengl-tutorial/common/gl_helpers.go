package common

import "encoding/binary"

func intFromByteSlice(bytes []byte) int {
	return int(binary.LittleEndian.Uint32(bytes))
}

func nullTerminatedString(source []byte) string {
	return string(append(source, make([]byte, 1)...))
}
