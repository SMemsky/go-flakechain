package portable

import (
	"errors"
	"reflect"
)

const (
	storageSignature = 0x0102010101011101

	serializeTypeInt64		= 1
	serializeTypeInt32		= 2
	serializeTypeInt16		= 3
	serializeTypeInt8		= 4
	serializeTypeUint64		= 5
	serializeTypeUint32		= 6
	serializeTypeUint16		= 7
	serializeTypeUint8		= 8
	serializeTypeFloat64	= 9
	serializeTypeString		= 10
	serializeTypeBool		= 11
	serializeTypeObject		= 12
	serializeTypeArray		= 13	// TODO

	serializeArrayMask		= 0x80

	portableVarint8		= 0
	portableVarint16	= 1
	portableVarint32	= 2
	portableVarint64	= 4
)

var (
	ErrBadVarint	= errors.New("storages/portable: varint is too big")
	ErrBadRoot		= errors.New("storages/portable: root is not a structure")
	ErrSecName		= errors.New("storages/portable: section name is too long")

	ErrSigMismatch	= errors.New("storages/portable: signature mismatch")
	ErrVerMismatch	= errors.New("storages/portable: version mismatch")

	ErrEntryCount	= errors.New("storages/portable: some entries are missing")
	ErrEntryMissing	= errors.New("storages/portable: entry is missing or duplicated")
	ErrTypeMismatch	= errors.New("storages/portable: type mismatch")

	ErrUnknownType	= errors.New("storages/portable: unknown serialize type")

	ErrBadArray		= errors.New("storages/portable: unknown slice type")
	ErrBadKind		= errors.New("storages/portable: array kind mismatch")
)

var (
	// SerializeType to reflect.Kind for array encoding (hacky-hacky)
	serializeType2Kind = map[uint8]reflect.Kind{
		serializeTypeInt64: reflect.Int64,
		serializeTypeInt32: reflect.Int32,
		serializeTypeInt16: reflect.Int16,
		serializeTypeInt8: reflect.Int8,
		serializeTypeUint64: reflect.Uint64,
		serializeTypeUint32: reflect.Uint32,
		serializeTypeUint16: reflect.Uint16,
		serializeTypeUint8: reflect.Uint8,
		serializeTypeFloat64: reflect.Float64,
		serializeTypeString: reflect.String,
		serializeTypeBool: reflect.Bool,
		serializeTypeObject: reflect.Struct,
	}
)

type storageHeader struct {
	Signature	uint64	// Always the same
	Version		uint8	// Always 1
}
