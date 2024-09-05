package core

type Obj struct {
	TypeEncoding   uint8
	Value          interface{}
	LastAccessedAt uint32
}

/**

OBJ_ENCODING_EMBSTR is an internal Redis encoding used for storing small strings (up to 44 bytes) in a single, contiguous memory allocation.
•	This encoding optimizes memory use and improves performance for short strings by reducing the number of memory allocations and keeping the object metadata and string data together.
•	Strings with lengths exceeding 44 bytes will be stored using the OBJ_ENCODING_RAW format, which separates the metadata from the string data.
*/

var OBJ_TYPE_STRING uint8 = 0 << 4

var OBJ_ENCODING_INT uint8 = 1
var OBJ_ENCODING_RAW uint8 = 0
var OBJ_ENCODING_EMBSTR uint8 = 8
