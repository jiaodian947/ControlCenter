package game

import (
	"robots/protocol"
	"robots/utils"
)

type ViewItem struct {
	ViewId   uint16
	ObjId    int
	ConfigId int32
	Attr     *AttrList
}

func NewViewItem() *ViewItem {
	item := &ViewItem{}
	item.Attr = NewAttrList()
	return item
}

func (v *ViewItem) LoadAttr(ar *utils.LoadArchive) {
	var props uint16
	CheckErr(ar.Read(&props))
	var index uint16
	for i := 0; i < int(props); i++ {
		CheckErr(ar.Read(&index))
		if int(index) >= len(PropTables.Props) {
			panic("index error")
		}
		attr := NewAny()
		attr.SetName(PropTables.Props[index].Name)
		switch int(PropTables.Props[index].Type) {
		case protocol.SC_TYPE_BYTE:
			val, err := ar.ReadUInt8()
			if err != nil {
				panic(err)
			}
			attr.SetInt(int32(val))
		case protocol.SC_TYPE_WORD:
			val, err := ar.ReadUInt16()
			if err != nil {
				panic(err)
			}
			attr.SetInt(int32(val))
		case protocol.SC_TYPE_DWORD:
			val, err := ar.ReadInt32()
			if err != nil {
				panic(err)
			}
			attr.SetInt(val)
		case protocol.SC_TYPE_QWORD:
			val, err := ar.ReadInt64()
			if err != nil {
				panic(err)
			}
			attr.SetInt64(val)
		case protocol.SC_TYPE_FLOAT:
			val, err := ar.ReadFloat32()
			if err != nil {
				panic(err)
			}
			attr.SetFloat(val)
		case protocol.SC_TYPE_DOUBLE:
			val, err := ar.ReadFloat64()
			if err != nil {
				panic(err)
			}
			attr.SetDouble(val)
		case protocol.SC_TYPE_STRING:
			fallthrough
		case protocol.SC_TYPE_WIDESTR:
			val, err := ar.ReadCStringWithLen()
			if err != nil {
				panic(err)
			}
			attr.SetString(val)
		case protocol.SC_TYPE_OBJECT:
			val, err := ar.ReadUInt64()
			if err != nil {
				panic(err)
			}
			attr.SetObjectId(val)
		}

		old := v.Attr.GetAttr(attr.Name())
		if old == nil {
			v.Attr.AddAttr(attr)
			continue
		}
		if old.Type() != attr.Type() {
			panic("type not match")
		}
		old.Val = attr.Val
	}
}

type View struct {
	ViewId   uint16
	Cap      uint16
	ConfigId int32
	Attr     *AttrList
	Childs   []*ViewItem
}

func (v *View) LoadAttr(ar *utils.LoadArchive) {
	var props uint16
	CheckErr(ar.Read(&props))
	var index uint16
	for i := 0; i < int(props); i++ {
		CheckErr(ar.Read(&index))
		if int(index) >= len(PropTables.Props) {
			panic("index error")
		}
		attr := NewAny()
		attr.SetName(PropTables.Props[index].Name)
		switch int(PropTables.Props[index].Type) {
		case protocol.SC_TYPE_BYTE:
			val, err := ar.ReadUInt8()
			if err != nil {
				panic(err)
			}
			attr.SetInt(int32(val))
		case protocol.SC_TYPE_WORD:
			val, err := ar.ReadUInt16()
			if err != nil {
				panic(err)
			}
			attr.SetInt(int32(val))
		case protocol.SC_TYPE_DWORD:
			val, err := ar.ReadInt32()
			if err != nil {
				panic(err)
			}
			attr.SetInt(val)
		case protocol.SC_TYPE_QWORD:
			val, err := ar.ReadInt64()
			if err != nil {
				panic(err)
			}
			attr.SetInt64(val)
		case protocol.SC_TYPE_FLOAT:
			val, err := ar.ReadFloat32()
			if err != nil {
				panic(err)
			}
			attr.SetFloat(val)
		case protocol.SC_TYPE_DOUBLE:
			val, err := ar.ReadFloat64()
			if err != nil {
				panic(err)
			}
			attr.SetDouble(val)
		case protocol.SC_TYPE_STRING:
			fallthrough
		case protocol.SC_TYPE_WIDESTR:
			val, err := ar.ReadCStringWithLen()
			if err != nil {
				panic(err)
			}
			attr.SetString(val)
		case protocol.SC_TYPE_OBJECT:
			val, err := ar.ReadUInt64()
			if err != nil {
				panic(err)
			}
			attr.SetObjectId(val)
		}

		old := v.Attr.GetAttr(attr.Name())
		if old == nil {
			v.Attr.AddAttr(attr)
			continue
		}
		if old.Type() != attr.Type() {
			panic("type not match")
		}
		old.Val = attr.Val
	}
}

func NewView(cap uint16) *View {
	v := &View{}
	v.Attr = NewAttrList()
	v.Cap = cap
	v.Childs = make([]*ViewItem, cap)
	return v
}

func (v *View) GetItemByIndex(index int) *ViewItem {
	if index >= len(v.Childs) {
		panic("index error")
	}

	return v.Childs[index]
}

func (v *View) AddItem(index int, item *ViewItem) {
	if index >= len(v.Childs) {
		panic("index error")
	}

	v.Childs[index] = item
}

func (v *View) RemoveItem(index int) {
	if index >= len(v.Childs) {
		panic("index error")
	}
	v.Childs[index] = nil
}

func (v *View) Exchange(src, dest int) {
	if src >= len(v.Childs) || dest >= len(v.Childs) {
		panic("index error")
	}
	v.Childs[src], v.Childs[dest] = v.Childs[dest], v.Childs[src]
}
