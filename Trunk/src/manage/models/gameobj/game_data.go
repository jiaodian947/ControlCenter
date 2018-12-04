package gameobj

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"io"
	"log"
	"manage/protocol"
)

func NewGameData() *GameData {
	g := &GameData{}
	g.Attrs = NewAttrList()
	g.Records = NewRecordList()
	return g
}

func NewGameDataFromBinary(data []byte) *GameData {
	g := NewGameData()
	if len(data) == 0 {
		g.DataVersion = GLOBAL_DATA_VERSION
		return g
	}
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
		}

		io.Copy(&out, r)
		ar := protocol.NewLoadArchiver(out.Bytes())
		g.Load(ar)
	} else {
		ar := protocol.NewLoadArchiver(data)
		g.Load(ar)
	}

	return g
}

type GameData struct {
	Compress    bool
	DataVersion int32
	Attrs       *AttrList
	Records     *RecordList
}

func (g *GameData) NeedCompress() bool {
	return g.Compress
}

func (g *GameData) Load(ar *protocol.LoadArchive) {
	MustRead(ar, &g.DataVersion)
	var attrcount int32
	var reccount int32
	MustRead(ar, &attrcount)
	MustRead(ar, &reccount)
	g.LoadAttr(ar, attrcount)
	g.LoadRecord(ar, reccount)
}

func (g *GameData) LoadAttr(ar *protocol.LoadArchive, attrcount int32) {
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
			val, err := ar.ReadCString()
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

func (g *GameData) LoadRecord(ar *protocol.LoadArchive, reccount int32) {
	for i := 0; i < int(reccount); i++ {
		var recname string
		var maxrow, rows, cols, keys int32
		MustRead(ar, &recname)
		MustRead(ar, &maxrow)
		MustRead(ar, &rows)
		MustRead(ar, &cols)
		MustRead(ar, &keys)
		rec := NewRecord(int(cols), int(maxrow))
		rec.SetName(recname)
		column := rec.Column()
		for c := 0; c < int(cols); c++ {
			var coltype uint8
			MustRead(ar, &coltype)
			column.SetType(c, int(coltype))
		}
		rec.ClearKeys()
		for c := 0; c < int(keys); c++ {
			var index int32
			var ci int8
			MustRead(ar, &index)
			MustRead(ar, &ci)
			rec.AddKey(index, ci)
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
					val, err := ar.ReadCString()
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

func (g *GameData) Store(ar *protocol.StoreArchive) error {
	if err := ar.Write(g.DataVersion); err != nil {
		return err
	}
	if err := ar.Write(int32(g.Attrs.Count())); err != nil {
		return err
	}
	if err := ar.Write(int32(g.Records.Count())); err != nil {
		return err
	}
	if err := g.StoreAttr(ar); err != nil {
		return err
	}
	if err := g.StoreRecord(ar); err != nil {
		return err
	}

	return nil
}

func (g *GameData) StoreAttr(ar *protocol.StoreArchive) error {
	attrcount := g.Attrs.Count()
	for i := 0; i < attrcount; i++ {
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

func (g *GameData) StoreRecord(ar *protocol.StoreArchive) error {
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
		if err := ar.Write(int32(rec.Keys())); err != nil {
			return err
		}
		col := rec.Column()
		for c := 0; c < rec.Cols(); c++ {
			if err := ar.Write(int8(col.Type(c))); err != nil {
				return err
			}
		}

		for c := 0; c < rec.Keys(); c++ {
			ki := rec.Key(c)
			if ki == nil {
				panic("index error")
			}
			if err := ar.Write(ki.Index); err != nil {
				return err
			}
			if err := ar.Write(ki.CaseInsensitive); err != nil {
				return err
			}
		}

		for r := 0; r < rec.RowCount(); r++ {
			row := rec.Row(r)
			for c := 0; c < rec.Cols(); c++ {
				if err := ar.Write(row.Value(c)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
