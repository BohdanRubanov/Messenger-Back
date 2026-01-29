package utils

const (
	timeCost = uint32(2)         // number of iterations
	memoryKB = uint32(64 * 1024) // memory cost in KiB
	threads  = uint8(2)          // number of parallel threads
	keyLen   = uint32(32)        // length of the derived key in bytes
	saltLen  = 16                // length of salt in bytes
)