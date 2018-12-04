package game

import (
	"fmt"
	"robots/protocol"
	"robots/utils"
	"time"
)

type PosInfo struct {
	Vec    Vector3D
	Orient float32
}

type DestInfo struct {
	Vec         Vector3D
	Orient      float32
	MoveSpeed   float32
	RotateSpeed float32
	JumpSpeed   float32
	Mode        int32
}

type GameObject struct {
	ObjId      uint64
	ConfigId   int32
	Attr       *AttrList
	Childs     []*GameObject
	KI         map[uint64]int
	Pos        *PosInfo
	Dest       *DestInfo
	Mode       int32
	MotionTime time.Time
}

func NewGameObject() *GameObject {
	g := &GameObject{}
	g.MotionTime = time.Now()
	g.Init()
	return g
}

func (g *GameObject) Init() {
	g.Attr = NewAttrList()
	g.Childs = make([]*GameObject, 0, 1024)
	g.KI = make(map[uint64]int)
}

func (g *GameObject) Motion(dest *DestInfo) {
	if dest.Mode == E_MODE_STOP {
		g.Pos.Vec = dest.Vec
		g.Pos.Orient = dest.Orient
		g.Location(g.Pos)
		return
	}
	g.Dest = dest
	g.Mode = dest.Mode
	g.MotionTime = time.Now()
}

func (g *GameObject) Location(pos *PosInfo) {
	g.Pos = pos
	g.Dest = nil
	g.MotionTime = time.Now()
}

func (g *GameObject) Position() Vector3D {
	if g.Dest == nil || g.MotionTime.Unix() == 0 {
		return g.Pos.Vec
	}

	dis := g.Dest.Vec.Distance(g.Pos.Vec)
	if dis <= 0.000001 {
		g.Location(&PosInfo{Vec: g.Dest.Vec, Orient: g.Dest.Orient})
		return g.Pos.Vec
	}
	now := time.Now()
	diff := now.Sub(g.MotionTime).Seconds()
	spd := float32(diff) * g.Dest.MoveSpeed

	if spd >= dis {
		g.Location(&PosInfo{Vec: g.Dest.Vec, Orient: g.Dest.Orient})
		return g.Pos.Vec
	}

	f := spd / dis
	g.Pos.Vec = Lerp(g.Pos.Vec, g.Dest.Vec, f)
	g.MotionTime = now
	return g.Pos.Vec
}

func (g *GameObject) AddObject(obj *GameObject) {
	index := -1
	for k, v := range g.Childs {
		if v == nil {
			index = k
			break
		}
	}
	if index == -1 {
		index = len(g.Childs)
		g.Childs = append(g.Childs, obj)

	} else {
		g.Childs[index] = obj
	}

	g.KI[obj.ObjId] = index
}

func (g *GameObject) FindObject(objid uint64) *GameObject {
	for _, v := range g.Childs {
		if v != nil && v.ObjId == objid {
			return v
		}
	}

	return nil
}

func (g *GameObject) RemoveObject(obj uint64) {
	if index, has := g.KI[obj]; has {
		delete(g.KI, obj)
		g.Childs[index] = nil
	}
}

func (g *GameObject) LoadAttr(ar *utils.LoadArchive) {
	var props uint16
	CheckErr(ar.Read(&props))
	var index uint16
	for i := 0; i < int(props); i++ {
		CheckErr(ar.Read(&index))
		if int(index) >= len(PropTables.Props) {
			panic(fmt.Sprintf("index error, %d", index))
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

		old := g.Attr.GetAttr(attr.Name())
		if old == nil {
			g.Attr.AddAttr(attr)
			continue
		}
		if old.Type() != attr.Type() {
			panic("type not match")
		}
		old.Val = attr.Val
	}
}
