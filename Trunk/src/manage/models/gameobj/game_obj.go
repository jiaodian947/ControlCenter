package gameobj

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"log"
	"manage/protocol"
)

func NewGameObject() *GameObject {
	g := &GameObject{}
	g.Attrs = NewAttrList()
	g.Records = NewRecordList()
	return g
}

func NewGameObjectFromBinary(data []byte) *GameObject {
	g := NewGameObject()
	version := binary.LittleEndian.Uint32(data)
	compress := false
	if version == COMPRESSED_DATA_VERSION {
		compress = true
	}
	g.Compress = compress
	if compress {
		srcLen := binary.LittleEndian.Uint32(data[4:])
		b := bytes.NewBuffer(data[8:])
		var out bytes.Buffer
		r, err := zlib.NewReader(b)
		if err != nil {
			log.Println(err, data[:8], version, srcLen, len(data))
			return nil
			//panic(err)
		}

		io.Copy(&out, r)
		if srcLen != uint32(out.Len()) {
			panic("len not match")
		}
		ar := protocol.NewLoadArchiver(out.Bytes())
		g.Load(ar)
	} else {
		ar := protocol.NewLoadArchiver(data)
		g.Load(ar)
	}

	return g
}

type GameObject struct {
	Compress    bool
	DataVersion int32
	ClassType   int32
	Script      string
	Config      int64
	ChildCount  int32
	Cap         int32
	Pos         int32
	Attrs       *AttrList
	Records     *RecordList
	Childs      []*GameObject
}

func (g *GameObject) NeedCompress() bool {
	return g.Compress
}

func MustRead(ar *protocol.LoadArchive, val interface{}) {
	err := ar.Read(val)
	if err != nil {
		panic(err)
	}
}

func (g *GameObject) Load(ar *protocol.LoadArchive) {
	MustRead(ar, &g.DataVersion)
	MustRead(ar, &g.ClassType)
	MustRead(ar, &g.Script)
	MustRead(ar, &g.Config)
	var attrcount int32
	var reccount int32
	MustRead(ar, &attrcount)
	MustRead(ar, &reccount)
	MustRead(ar, &g.ChildCount)
	MustRead(ar, &g.Cap)
	MustRead(ar, &g.Pos)
	g.LoadAttr(ar, attrcount)
	g.LoadRecord(ar, reccount)
	if g.ChildCount > 0 {
		g.Childs = make([]*GameObject, 0, g.ChildCount)
		for i := 0; i < int(g.ChildCount); i++ {
			obj := NewGameObject()
			obj.Load(ar)
			g.Childs = append(g.Childs, obj)
		}
	}
}

func (g *GameObject) LoadAttr(ar *protocol.LoadArchive, attrcount int32) {
	for i := 0; i < int(attrcount); i++ {
		var key string
		var typ uint8
		MustRead(ar, &key)
		attr := NewAny()
		attr.SetName(key)
		MustRead(ar, &typ)
		switch typ {
		case protocol.VTYPE_INT:
			val, err := ar.ReadInt32()
			if err != nil {
				panic(err)
			}
			attr.SetInt(val)
		case protocol.VTYPE_INT64:
			val, err := ar.ReadInt64()
			if err != nil {
				panic(err)
			}
			attr.SetInt64(val)
		case protocol.VTYPE_FLOAT:
			val, err := ar.ReadFloat32()
			if err != nil {
				panic(err)
			}
			attr.SetFloat(val)
		case protocol.VTYPE_DOUBLE:
			val, err := ar.ReadFloat64()
			if err != nil {
				panic(err)
			}
			attr.SetDouble(val)
		case protocol.VTYPE_STRING:
			fallthrough
		case protocol.VTYPE_WIDESTR:
			val, err := ar.ReadCStringWithLen()
			if err != nil {
				panic(err)
			}
			attr.SetString(val)
		case protocol.VTYPE_OBJECT:
			val, err := ar.ReadUInt64()
			if err != nil {
				panic(err)
			}
			attr.SetObjectId(val)
		default:
			panic("unsupport type")

		}
		g.Attrs.AddAttr(attr)
	}
}

func (g *GameObject) LoadRecord(ar *protocol.LoadArchive, reccount int32) {
	for i := 0; i < int(reccount); i++ {
		var recname string
		var maxrow, rows, cols int32
		MustRead(ar, &recname)
		MustRead(ar, &maxrow)
		MustRead(ar, &rows)
		MustRead(ar, &cols)
		rec := NewRecord(int(cols), int(maxrow))
		rec.SetName(recname)
		column := rec.Column()
		for c := 0; c < int(cols); c++ {
			var coltype uint8
			MustRead(ar, &coltype)
			column.SetType(c, int(coltype))
		}
		for r := 0; r < int(rows); r++ {
			row := NewRow(int(cols))
			for c := 0; c < int(cols); c++ {
				switch column.Type(c) {
				case protocol.VTYPE_INT:
					attr := NewAny()
					val, err := ar.ReadInt32()
					if err != nil {
						panic(err)
					}
					attr.SetInt(val)
					row.SetValue(c, attr)
				case protocol.VTYPE_INT64:
					attr := NewAny()
					val, err := ar.ReadInt64()
					if err != nil {
						panic(err)
					}
					attr.SetInt64(val)
					row.SetValue(c, attr)
				case protocol.VTYPE_FLOAT:
					attr := NewAny()
					val, err := ar.ReadFloat32()
					if err != nil {
						panic(err)
					}
					attr.SetFloat(val)
					row.SetValue(c, attr)
				case protocol.VTYPE_DOUBLE:
					attr := NewAny()
					val, err := ar.ReadFloat64()
					if err != nil {
						panic(err)
					}
					attr.SetDouble(val)
					row.SetValue(c, attr)
				case protocol.VTYPE_STRING:
					fallthrough
				case protocol.VTYPE_WIDESTR:
					attr := NewAny()
					val, err := ar.ReadCStringWithLen()
					if err != nil {
						panic(err)
					}
					attr.SetString(val)
					row.SetValue(c, attr)
				case protocol.VTYPE_OBJECT:
					attr := NewAny()
					val, err := ar.ReadUInt64()
					if err != nil {
						panic(err)
					}
					attr.SetObjectId(val)
					row.SetValue(c, attr)
				default:
					panic("unsupport type")
				}

			}
			rec.AddRowValue(r, row)
		}
		g.Records.AddRecord(recname, rec)
	}
}

func (g *GameObject) Store(ar *protocol.StoreArchive) error {
	if err := ar.Write(g.DataVersion); err != nil {
		return err
	}
	if err := ar.Write(g.ClassType); err != nil {
		return err
	}
	if err := ar.Write(g.Script); err != nil {
		return err
	}
	if err := ar.Write(g.Config); err != nil {
		return err
	}
	if err := ar.Write(int32(g.Attrs.Count())); err != nil {
		return err
	}
	if err := ar.Write(int32(g.Records.Count())); err != nil {
		return err
	}
	if err := ar.Write(int32(len(g.Childs))); err != nil {
		return err
	}
	if err := ar.Write(g.Cap); err != nil {
		return err
	}
	if err := ar.Write(g.Pos); err != nil {
		return err
	}
	if err := g.StoreAttr(ar); err != nil {
		return err
	}
	if err := g.StoreRecord(ar); err != nil {
		return err
	}
	if len(g.Childs) > 0 {
		for _, v := range g.Childs {
			if err := v.Store(ar); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *GameObject) StoreAttr(ar *protocol.StoreArchive) error {
	size := g.Attrs.Count()
	for i := 0; i < size; i++ {
		attr := g.Attrs.GetAttrByIndex(i)
		if err := ar.Write(attr.Name()); err != nil {
			return err
		}
		if err := ar.Write(int8(attr.Type())); err != nil {
			return err
		}
		if err := ar.Write(attr.Value()); err != nil {
			return err
		}
	}
	return nil
}

func (g *GameObject) StoreRecord(ar *protocol.StoreArchive) error {
	size := g.Records.Size()
	for i := 0; i < size; i++ {
		rec := g.Records.RecordByIndex(i)
		if rec == nil {
			continue
		}
		if err := ar.Write(rec.Name()); err != nil {
			return err
		}
		if err := ar.Write(int32(rec.MaxRows())); err != nil {
			return err
		}
		if err := ar.Write(int32(rec.RowCount())); err != nil {
			return err
		}
		if err := ar.Write(int32(rec.Cols())); err != nil {
			return err
		}
		col := rec.Column()
		for c := 0; c < rec.Cols(); c++ {
			if err := ar.Write(int8(col.Type(c))); err != nil {
				return err
			}
		}

		for r := 0; r < rec.RowCount(); r++ {
			row := rec.Row(r)
			for c := 0; c < rec.Cols(); c++ {
				if err := ar.Write(row.Value(c).Value()); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
