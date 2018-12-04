package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type StoreArchive struct {
	buffer *bytes.Buffer
}

func NewStoreArchiver(data []byte) *StoreArchive {
	ar := &StoreArchive{}
	ar.buffer = bytes.NewBuffer(data)
	return ar
}

func (ar *StoreArchive) Data() []byte {
	return ar.buffer.Bytes()
}

func (ar *StoreArchive) Len() int {
	return ar.buffer.Len()
}

func (ar *StoreArchive) WriteAt(offset int, val interface{}) error {
	if offset >= ar.buffer.Len() {
		return fmt.Errorf("offset out of range")
	}

	data := ar.buffer.Bytes()
	tmp := bytes.NewBuffer(data[offset:offset])
	switch val.(type) {
	case int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return binary.Write(tmp, binary.LittleEndian, val)
	case int:
		return binary.Write(tmp, binary.LittleEndian, int32(val.(int)))
	default:
		return fmt.Errorf("unsupport type")
	}
}

func (ar *StoreArchive) Write(val interface{}) error {
	switch val.(type) {
	case bool, int8, int16, int32, int64, uint8, uint16, uint32, uint64, float32, float64:
		return binary.Write(ar.buffer, binary.LittleEndian, val)
	case int:
		return binary.Write(ar.buffer, binary.LittleEndian, int64(val.(int)))
	case string:
		return ar.WriteCStringWithLen(val.(string))
	case []byte:
		return ar.WriteData(val.([]byte))
	case *VarMessage:
		return ar.WriteVarmsg(val.(*VarMessage))
	default:
		return fmt.Errorf("unsupport type")
	}
}

func (ar *StoreArchive) WriteVarmsg(v *VarMessage) error {
	count := v.Count()
	var err error
	err = ar.Write(uint16(count))
	if err != nil {
		return err
	}
	for i := 0; i < count; i++ {
		ar.Write(uint8(v.Type(i)))
		switch v.Type(i) {
		case VTYPE_BOOL:
			err = ar.Write(v.BoolVal(i))
		case VTYPE_BYTE:
			err = ar.Write(v.ByteVal(i))
		case VTYPE_WORD:
			err = ar.Write(v.WordVal(i))
		case VTYPE_INT:
			err = ar.Write(v.Int32Val(i))
		case VTYPE_INT64:
			err = ar.Write(v.Int64Val(i))
		case VTYPE_FLOAT:
			err = ar.Write(v.FloatVal(i))
		case VTYPE_DOUBLE:
			err = ar.Write(v.DoubleVal(i))
		case VTYPE_STRING:
			err = ar.Write(v.StringVal(i))
		default:
			return fmt.Errorf("unsupport type")
		}
		if err != nil {
			return err
		}
	}
	return err
}

// 写入包含前置长度的字符串,结尾加\0
func (ar *StoreArchive) WriteCStringWithLen(val string) error {
	data := []byte(val)
	size := len(data) + 1 //包含结尾的0
	err := binary.Write(ar.buffer, binary.LittleEndian, int32(size))
	if err != nil {
		return err
	}
	if _, err = ar.buffer.Write(data); err != nil {
		return err
	}

	if err = ar.buffer.WriteByte(0); err != nil { //结尾写0
		return err
	}
	return nil
}

// 写入字符中，结尾加\0
func (ar *StoreArchive) WriteCString(val string) error {
	data := []byte(val)
	if _, err := ar.buffer.Write(data); err != nil {
		return err
	}

	if err := ar.buffer.WriteByte(0); err != nil { //结尾写0
		return err
	}
	return nil
}

func (ar *StoreArchive) WriteData(data []byte) error {
	err := ar.Write(uint16(len(data)))
	if err != nil {
		return err
	}
	_, err = ar.buffer.Write(data)
	return err
}

type LoadArchive struct {
	reader *bytes.Reader
	s      []byte
}

func NewLoadArchiver(data []byte) *LoadArchive {
	ar := &LoadArchive{}
	ar.reader = bytes.NewReader(data)
	ar.s = data
	return ar
}

func (ar *LoadArchive) Source() []byte {
	return ar.s
}

func (ar *LoadArchive) Position() int {
	return int(ar.reader.Size()) - ar.reader.Len()
}

func (ar *LoadArchive) AvailableBytes() int {
	return ar.reader.Len()
}

func (ar *LoadArchive) Size() int {
	return int(ar.reader.Size())
}

func (ar *LoadArchive) Seek(offset int, whence int) (int, error) {
	ret, err := ar.reader.Seek(int64(offset), whence)
	return int(ret), err
}

func (ar *LoadArchive) Read(val interface{}) (err error) {
	switch val.(type) {
	case *bool, *int8, *int16, *int32, *int64, *uint8, *uint16, *uint32, *uint64, *float32, *float64:
		return binary.Read(ar.reader, binary.LittleEndian, val)
	case *int:
		var out int64
		err = binary.Read(ar.reader, binary.LittleEndian, &out)
		if err != nil {
			return err
		}
		*(val.(*int)) = int(out)
		return nil
	case *string:
		inst := val.(*string)
		*inst, err = ar.ReadCStringWithLen()
		return err
	case *[]byte:
		inst := val.(*[]byte)
		*inst, err = ar.ReadData()
		return err

	default:
		return fmt.Errorf("unsupport type")
	}
}

func (ar *LoadArchive) ReadVarMsg() (*VarMessage, error) {
	count, err := ar.ReadUInt16()
	if err != nil {
		return nil, err
	}

	var typ uint8

	varmsg := NewVarMsg(int(count))
	for i := 0; i < int(count); i++ {
		ar.Read(&typ)
		switch typ {
		case VTYPE_BOOL:
			var value bool
			err = ar.Read(&value)
			if err != nil {
				return nil, err
			}
			varmsg.AddBool(value)
		case VTYPE_BYTE:
			var value int8
			err = ar.Read(&value)
			if err != nil {
				return nil, err
			}
			varmsg.AddByte(value)
		case VTYPE_WORD:
			var value int16
			err = ar.Read(&value)
			if err != nil {
				return nil, err
			}
			varmsg.AddWord(value)
		case VTYPE_INT:
			var value int32
			err = ar.Read(&value)
			if err != nil {
				return nil, err
			}
			varmsg.AddInt32(value)
		case VTYPE_INT64:
			var value int64
			err = ar.Read(&value)
			if err != nil {
				return nil, err
			}
			varmsg.AddInt64(value)
		case VTYPE_FLOAT:
			var value float32
			err = ar.Read(&value)
			if err != nil {
				return nil, err
			}
			varmsg.AddFloat(value)
		case VTYPE_DOUBLE:
			var value float64
			err = ar.Read(&value)
			if err != nil {
				return nil, err
			}
			varmsg.AddDouble(value)
		case VTYPE_STRING:
			var value string
			err = ar.Read(&value)
			if err != nil {
				return nil, err
			}
			varmsg.AddString(value)
		default:
			return nil, fmt.Errorf("unsupport type")
		}
	}
	return varmsg, err
}

func (ar *LoadArchive) ReadInt8() (val int8, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadUInt8() (val uint8, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadInt16() (val int16, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadUInt16() (val uint16, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadInt32() (val int32, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadUInt32() (val uint32, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadInt64() (val int64, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadUInt64() (val uint64, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadFloat32() (val float32, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadFloat64() (val float64, err error) {
	err = ar.Read(&val)
	return
}

func (ar *LoadArchive) ReadCStringWithLen() (val string, err error) {
	var size int32
	binary.Read(ar.reader, binary.LittleEndian, &size)
	if size == 0 {
		val = ""
		return
	}
	data := make([]byte, size)
	_, err = ar.reader.Read(data)
	if err != nil {
		return
	}
	if data[size-1] == 0 {
		data = data[:size-1]
	}
	val = string(data)
	return
}

func (ar *LoadArchive) ReadCString() (val string, err error) {
	buf := make([]byte, 0, 128)
	ch, err := ar.reader.ReadByte()
	for ch != 0 && err == nil {
		buf = append(buf, ch)
		ch, err = ar.reader.ReadByte()
	}
	val = string(buf)
	return
}

func (ar *LoadArchive) ReadData() (data []byte, err error) {
	var l uint16
	l, err = ar.ReadUInt16()
	data = make([]byte, int(l))
	_, err = ar.reader.Read(data)
	return data, err
}
