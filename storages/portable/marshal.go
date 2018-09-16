package portable

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	if err := encode(&b, v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func encode(w io.Writer, v interface{}) error {
	if err := binary.Write(w, binary.LittleEndian, storageHeader{storageSignature, 1}); err != nil {
		return err
	}

	rv := reflect.Indirect(reflect.ValueOf(v))
	if err := encodeStruct(w, rv); err != nil {
		return err
	}

	return nil
}

func encodeStruct(w io.Writer, v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return ErrBadRoot
	}

	t := v.Type()
	l := v.NumField()
	entryCount := uint64(0)
	for i := 0; i < l; i++ {
		if _, ok := t.Field(i).Tag.Lookup("store"); ok {
			entryCount++
		}
	}

	if err := encodeVarint(w, entryCount); err != nil {
		return err
	}

	for i := 0; i < l; i++ {
		if tag, ok := t.Field(i).Tag.Lookup("store"); ok {
			if len(tag) > 0xff {
				return ErrSecName
			}
			if err := binary.Write(w, binary.LittleEndian, uint8(len(tag))); err != nil {
				return err
			}
			if err := binary.Write(w, binary.LittleEndian, []byte(tag)); err != nil {
				return err
			}
			if err := encodeValue(w, v.Field(i)); err != nil {
				return err
			}
		}
	}

	return nil
}

func encodeValue(w io.Writer, v reflect.Value) error {
	switch v.Kind() {
	case reflect.Struct:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeObject)); err != nil {
			return err
		}
		return encodeStruct(w, v)
	case reflect.String:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeString)); err != nil {
			return err
		}
		if err := encodeVarint(w, uint64(len(v.String()))); err != nil {
			return err
		}
		if err := binary.Write(w, binary.LittleEndian, []byte(v.String())); err != nil {
			return err
		}
	case reflect.Int64:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeInt64)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, int64(v.Int()))
	case reflect.Int32:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeInt32)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, int32(v.Int()))
	case reflect.Int16:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeInt16)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, int16(v.Int()))
	case reflect.Int8:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeInt8)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, int8(v.Int()))
	case reflect.Uint64:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeUint64)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint64(v.Uint()))
	case reflect.Uint32:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeUint32)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint32(v.Uint()))
	case reflect.Uint16:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeUint16)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint16(v.Uint()))
	case reflect.Uint8:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeUint8)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, uint8(v.Uint()))
	case reflect.Float64:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeFloat64)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, float64(v.Float()))
	case reflect.Bool:
		if err := binary.Write(w, binary.LittleEndian, uint8(serializeTypeBool)); err != nil {
			return err
		}
		return binary.Write(w, binary.LittleEndian, bool(v.Bool()))
	}

	return nil
}

func encodeVarint(w io.Writer, value uint64) error {
	if value <= 0x3f {
		return binary.Write(w, binary.LittleEndian, uint8(value << 2) | uint8(portableVarint8))
	} else if value <= 0x3fff {
		return binary.Write(w, binary.LittleEndian, uint16(value << 2) | uint16(portableVarint16))
	} else if value <= 0x3fffffff {
		return binary.Write(w, binary.LittleEndian, uint32(value << 2) | uint32(portableVarint32))
	} else if value <= 0x3fffffffffffffff {
		return binary.Write(w, binary.LittleEndian, uint64(value << 2) | uint64(portableVarint64))
	}

	return ErrBadVarint
}
