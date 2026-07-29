package ssg

import "hash"

type BytesObject interface {
	ToBytes() ([]byte, error)
	Length() int
}

type IntegerRepresent interface {
	ToInt64() int64
	ToUInt64() uint64
	ToInt32() int32
	ToUInt32() uint32
}

type BitsBlocks interface {
	GetBitsSize() int
}

type Serializer interface {
	Serialize() ([]byte, error)
	StrSerialize() string
}

type SignatureContainer interface {
	GetSignature() string
	SetSignature(signature string) bool
	SetSignatureByBytes(data []byte) bool
	SetSignatureByFunc(h func() hash.Hash) bool
}

// MetaDataProvider interface provides useful methods for getting/setting
// metadata for a data structure.
type MetaDataProvider interface {
	Get(key string) (string, error)
	GetInt(key string) (int, error)
	GetInt8(key string) (int8, error)
	GetInt16(key string) (int16, error)
	GetInt32(key string) (int32, error)
	GetInt64(key string) (int64, error)
	GetUInt(key string) (uint, error)
	GetUInt8(key string) (uint8, error)
	GetUInt16(key string) (uint16, error)
	GetUInt32(key string) (uint32, error)
	GetUInt64(key string) (uint64, error)
	GetBool(key string) (bool, error)

	GetNoErr(key string) string
	GetIntNoErr(key string) int
	GetInt8NoErr(key string) int8
	GetInt16NoErr(key string) int16
	GetInt32NoErr(key string) int32
	GetInt64NoErr(key string) int64
	GetUIntNoErr(key string) uint
	GetUInt8NoErr(key string) uint8
	GetUInt16NoErr(key string) uint16
	GetUInt32NoErr(key string) uint32
	GetUInt64NoErr(key string) uint64
	GetBoolNoErr(key string) bool

	Set(key, value string)
	SetInt(key string, value int)
	SetInt8(key string, value int8)
	SetInt16(key string, value int16)
	SetInt32(key string, value int32)
	SetInt64(key string, value int64)
	SetUInt(key string, value uint)
	SetUInt8(key string, value uint8)
	SetUInt16(key string, value uint16)
	SetUInt32(key string, value uint32)
	SetUInt64(key string, value uint64)
	SetBool(key string, value bool)
}
