package protocol

import (
	"fmt"
	"robots/utils"
)

func PubArgs(ar *utils.StoreArchive, args *VarMessage) {
	ar.Write(uint16(args.Size))
	for i := 0; i < args.Size; i++ {
		typ := args.Type(i)
		ar.Write(int8(typ))
		switch typ {
		case E_VTYPE_BYTE:
			ar.Write(args.ByteVal(i))
		case E_VTYPE_INT:
			ar.Write(args.Int32Val(i))
		case E_VTYPE_INT64:
			ar.Write(args.Int64Val(i))
		case E_VTYPE_FLOAT:
			ar.Write(args.FloatVal(i))
		case E_VTYPE_DOUBLE:
			ar.Write(args.DoubleVal(i))
		case E_VTYPE_STRING:
			ar.Write(args.StringVal(i))
		case E_VTYPE_OBJECT:
			ar.Write(args.ObjectVal(i))
		default:
			panic(fmt.Sprint("not support type", typ))
		}
	}
}

func ParseArgs(ar *utils.LoadArchive) *VarMessage {
	nums, err := ar.ReadUInt16()
	if err != nil {
		panic(err)
	}
	args := NewVarMsg(int(nums))
	for i := 0; i < int(nums); i++ {
		typ, err := ar.ReadInt8()
		if err != nil {
			panic(err)
		}
		switch typ {
		case E_VTYPE_BYTE:
			val, err := ar.ReadInt8()
			if err != nil {
				panic(err)
			}
			args.AddByte(val)
		case E_VTYPE_INT:
			val, err := ar.ReadInt32()
			if err != nil {
				panic(err)
			}
			args.AddInt32(val)
		case E_VTYPE_INT64:
			val, err := ar.ReadInt64()
			if err != nil {
				panic(err)
			}
			args.AddInt64(val)
		case E_VTYPE_FLOAT:
			val, err := ar.ReadFloat32()
			if err != nil {
				panic(err)
			}
			args.AddFloat(val)
		case E_VTYPE_DOUBLE:
			val, err := ar.ReadFloat64()
			if err != nil {
				panic(err)
			}
			args.AddDouble(val)
		case E_VTYPE_STRING:
			val, err := ar.ReadCStringWithLen()
			if err != nil {
				panic(err)
			}
			args.AddString(val)
		case E_VTYPE_OBJECT:
			val, err := ar.ReadUInt64()
			if err != nil {
				panic(err)
			}
			args.AddObject(val)
		default:
			panic(fmt.Sprint("not support type", typ))
		}
	}
	return args
}
