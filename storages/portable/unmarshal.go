package portable

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) error {
	b := bytes.NewBuffer(data)
	return decode(b, v)
}

func decode(r io.Reader, v interface{}) error {
	header := &storageHeader{}
	if err := binary.Read(r, binary.LittleEndian, header); err != nil {
		return err
	}
	if header.Signature != storageSignature {
		return ErrSigMismatch
	}
	if header.Version != 1 {
		return ErrVerMismatch
	}
	rv := reflect.Indirect(reflect.ValueOf(v))
	if err := decodeStruct(r, rv); err != nil {
		return err
	}
	return nil
}

func decodeStruct(r io.Reader, v reflect.Value) error {
	if v.Kind() != reflect.Struct {
		return ErrBadRoot
	}

	t := v.Type()
	l := v.NumField()
	usedFields := make(map[string]struct{}) // imitate a set with zero byte structs
	for i := 0; i < l; i++ {
		if tag, ok := t.Field(i).Tag.Lookup("store"); ok {
			usedFields[tag] = struct{}{}
		}
	}

	c, err := decodeVarint(r)
	if err != nil {
		return err
	} else if c < uint64(len(usedFields)) {
		return ErrEntryCount
	}

	for i := uint64(0); i < c; i++ {
		var strSize uint8
		binary.Read(r, binary.LittleEndian, &strSize)
		strBuf := make([]byte, strSize)
		binary.Read(r, binary.LittleEndian, strBuf)
		str := string(strBuf)
		if _, prs := usedFields[str]; !prs {
			return fmt.Errorf("%s: %s", ErrEntryMissing, str)
		}
		delete(usedFields, str)
		for j := 0; j < l; j++ {
			if tag, ok := t.Field(j).Tag.Lookup("store"); ok && tag == str {
				if err := decodeEntry(r, v.Field(j)); err != nil {
					return err
				}
				break
			}
		}
	}

	if len(usedFields) != 0 {
		return ErrEntryMissing
	}

	return nil
}

func decodeEntry(r io.Reader, v reflect.Value) error {
	var valueType uint8
	if err := binary.Read(r, binary.LittleEndian, &valueType); err != nil {
		return err
	}
	if valueType&serializeArrayMask != 0 {
		return decodeArray(r, v, valueType&^serializeArrayMask)
	}
	return decodeValue(r, v, valueType)
}

// TODO: indexing varint on 32-bit targets
func decodeArray(r io.Reader, v reflect.Value, valueType uint8) error {
	if v.Kind() != reflect.Slice {
		return ErrBadArray
	}
	if serializeType2Kind[valueType] != v.Type().Elem().Kind() {
		return fmt.Errorf("%s: %s", ErrBadKind, v.Type().Elem().Kind())
	}

	count, err := decodeVarint(r)
	if err != nil {
		return err
	}
	v.Set(reflect.MakeSlice(v.Type(), int(count), int(count)))
	for i := uint64(0); i < count; i++ {
		decodeValue(r, v.Index(int(i)), valueType)
	}
	return nil
}

func decodeValue(r io.Reader, v reflect.Value, valueType uint8) error {
	switch valueType {
	case serializeTypeObject:
		return decodeStruct(r, v)
	case serializeTypeString:
		length, err := decodeVarint(r)
		if err != nil {
			return err
		}
		strBuffer := make([]byte, length)
		if err := binary.Read(r, binary.LittleEndian, strBuffer); err != nil {
			return err
		}
		v.SetString(string(strBuffer))
	case serializeTypeInt64:
		var value int64
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetInt(value)
	case serializeTypeInt32:
		var value int32
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetInt(int64(value))
	case serializeTypeInt16:
		var value int16
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetInt(int64(value))
	case serializeTypeInt8:
		var value int8
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetInt(int64(value))
	case serializeTypeUint64:
		var value uint64
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetUint(value)
	case serializeTypeUint32:
		var value uint32
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetUint(uint64(value))
	case serializeTypeUint16:
		var value uint16
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetUint(uint64(value))
	case serializeTypeUint8:
		var value uint8
		if err := binary.Read(r, binary.LittleEndian, &value); err != nil {
			return err
		}
		v.SetUint(uint64(value))
	default:
		return ErrUnknownType
	}
	return nil
}

func decodeVarint(r io.Reader) (uint64, error) {
	buf := make([]byte, 8) // 8 bytes at most
	if _, err := io.ReadFull(r, buf[:1]); err != nil {
		return 0, err
	}
	varintType := uint8(buf[0]) & 0x3
	if varintType == portableVarint8 {
		return uint64(buf[0]) >> 2, nil
	} else if varintType == portableVarint16 {
		if _, err := io.ReadFull(r, buf[1:2]); err != nil {
			return 0, err
		}
		b := bytes.NewBuffer(buf[:2])
		var value uint16
		binary.Read(b, binary.LittleEndian, &value)
		return uint64(value) >> 2, nil
	} else if varintType == portableVarint32 {
		if _, err := io.ReadFull(r, buf[1:4]); err != nil {
			return 0, err
		}
		b := bytes.NewBuffer(buf[:4])
		var value uint32
		binary.Read(b, binary.LittleEndian, &value)
		return uint64(value) >> 2, nil
	} else {
		if _, err := io.ReadFull(r, buf[1:8]); err != nil {
			return 0, err
		}
		b := bytes.NewBuffer(buf[:8])
		var value uint64
		binary.Read(b, binary.LittleEndian, &value)
		return value >> 2, nil
	}
}
