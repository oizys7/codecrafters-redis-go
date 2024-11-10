package _type

import "time"

const (
	// StringType is redis string
	StringType = "string"
	// ListType is redis list
	ListType = "list"
	// SetType is redis set
	SetType = "set"
	// HashType is redis hash
	HashType = "hash"
	// ZSetType is redis sorted set
	ZSetType = "zset"
	// AuxType is redis metadata key-value pair
	AuxType = "aux"
	// DBSizeType is for RDB_OPCODE_RESIZEDB
	DBSizeType = "dbsize"
	// StreamType is a redis stream
	StreamType = "stream"
)

const (
	// StringEncoding for string
	StringEncoding = "string"
	// ListEncoding is formed by a length encoding and some string
	ListEncoding = "list"
	// SetEncoding is formed by a length encoding and some string
	SetEncoding = "set"
	// ZSetEncoding is formed by a length encoding and some string
	ZSetEncoding = "zset"
	// HashEncoding is formed by a length encoding and some string
	HashEncoding = "hash"
	// ZSet2Encoding is zset version2 which stores doubles in binary format
	ZSet2Encoding = "zset2"
	// ZipMapEncoding has been deprecated
	ZipMapEncoding = "zipmap"
	// ZipListEncoding  stores data in contiguous memory
	ZipListEncoding = "ziplist"
	// IntSetEncoding is a ordered list of integers
	IntSetEncoding = "intset"
	// QuickListEncoding is a list of ziplist
	QuickListEncoding = "quicklist"
	// ListPackEncoding is a new replacement for ziplist
	ListPackEncoding = "listpack"
	// QuickList2Encoding is a list of listpack
	QuickList2Encoding = "quicklist2"
)

type RedisObject interface {
	// GetType returns redis type of object: string/list/set/hash/zset
	GetType() string
	// GetKey returns key of object
	GetKey() string
	// GetDBIndex returns db index of object
	GetDBIndex() int
	// GetExpiration returns expiration time, expiration of persistent object is nil
	GetExpiration() *time.Time
	// GetSize returns rdb value size in Byte
	GetSize() int
	// GetElemCount returns number of elements in list/set/hash/zset
	GetElemCount() int
	// GetEncoding returns encoding of object
	GetEncoding() string
}

// BaseObject is basement of redis object
type BaseObject struct {
	DB         int         `json:"db"`                   // DB is db index of redis object
	Key        string      `json:"key"`                  // Key is key of redis object
	Expiration *time.Time  `json:"expiration,omitempty"` // Expiration is expiration time, expiration of persistent object is nil
	Size       int         `json:"size"`                 // Size is rdb value size in Byte
	Type       string      `json:"type"`                 // Type is one of string/list/set/hash/zset
	Encoding   string      `json:"encoding"`             // Encoding is the exact encoding method
	Extra      interface{} `json:"-"`                    // Extra stores more detail of encoding for memory profiler and other usages
}

// GetKey returns key of object
func (o *BaseObject) GetKey() string {
	return o.Key
}

// GetDBIndex returns db index of object
func (o *BaseObject) GetDBIndex() int {
	return o.DB
}

// GetEncoding returns encoding of object
func (o *BaseObject) GetEncoding() string {
	return o.Encoding
}

// GetExpiration returns expiration time, expiration of persistent object is nil
func (o *BaseObject) GetExpiration() *time.Time {
	return o.Expiration
}

// GetSize  returns rdb value size in Byte
func (o *BaseObject) GetSize() int {
	return o.Size
}

// GetElemCount returns number of elements in list/set/hash/zset
func (o *BaseObject) GetElemCount() int {
	return 0
}
