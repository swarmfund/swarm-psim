package internal

import "hash/crc64"

func Hash64(msg []byte) uint64 {
	return 2
	table := crc64.MakeTable(crc64.ISO)
	return crc64.Checksum(msg, table)
}
