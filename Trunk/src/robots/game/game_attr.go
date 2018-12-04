package game

import (
	"fmt"
	"robots/protocol"
	"robots/utils"
)

type Any struct {
	typ int
	Key string
	Val interface{}
}

func NewAny() *Any {
	a := &Any{}
	return a
}

func (a *Any) SetName(name string) {
	a.Key = name
}

func (a *Any) Name() string {
	return a.Key
}

func (a *Any) Type() int {
	return a.typ
}

func (a *Any) Value() interface{} {
	return a.Val
}

func (a *Any) SetInt(value int32) {
	a.typ = protocol.E_VTYPE_INT
	a.Val = value
}

func (a *Any) SetInt64(value int64) {
	a.typ = protocol.E_VTYPE_INT64
	a.Val = value
}

func (a *Any) SetFloat(value float32) {
	a.typ = protocol.E_VTYPE_FLOAT
	a.Val = value
}

func (a *Any) SetDouble(value float64) {
	a.typ = protocol.E_VTYPE_DOUBLE
	a.Val = value
}

func (a *Any) SetString(value string) {
	a.typ = protocol.E_VTYPE_STRING
	a.Val = value
}

func (a *Any) SetObjectId(value uint64) {
	a.typ = protocol.E_VTYPE_OBJECT
	a.Val = value
}

func (a *Any) Int() int32 {
	if a.typ == protocol.E_VTYPE_INT {
		return a.Val.(int32)
	}
	return 0
}

func (a *Any) Int64() int64 {
	if a.typ == protocol.E_VTYPE_INT64 {
		return a.Val.(int64)
	}
	return 0
}

func (a *Any) Float() float32 {
	if a.typ == protocol.E_VTYPE_FLOAT {
		return a.Val.(float32)
	}
	return 0
}

func (a *Any) Double() float64 {
	if a.typ == protocol.E_VTYPE_DOUBLE {
		return a.Val.(float64)
	}
	return 0
}

func (a *Any) String() string {
	if a.typ == protocol.E_VTYPE_STRING {
		return a.Val.(string)
	}
	return ""
}

func (a *Any) ObjectId() uint64 {
	if a.typ == protocol.E_VTYPE_OBJECT {
		return a.Val.(uint64)
	}
	return 0
}

func (a *Any) Less(other *Any) bool {
	if a.typ != other.typ {
		return false
	}

	switch a.typ {
	case protocol.E_VTYPE_INT:
		return a.Int() < other.Int()
	case protocol.E_VTYPE_INT64:
		return a.Int64() < other.Int64()
	case protocol.E_VTYPE_FLOAT:
		return a.Float() < other.Float()
	case protocol.E_VTYPE_DOUBLE:
		return a.Double() < other.Double()
	case protocol.E_VTYPE_STRING:
		fallthrough
	case protocol.E_VTYPE_WIDESTR:
		return a.String() < other.String()
	}
	return false
}

func (a *Any) ToString() string {
	return utils.AsString(a.Val)
}

func (a *Any) FromString(s string) error {
	switch a.typ {
	case protocol.E_VTYPE_INT:
		var val int32
		if err := utils.ParseStrNumber(s, &val); err != nil {
			return err
		}
		a.Val = val
	case protocol.E_VTYPE_INT64:
		var val int64
		if err := utils.ParseStrNumber(s, &val); err != nil {
			return err
		}
		a.Val = val
	case protocol.E_VTYPE_FLOAT:
		var val float32
		if err := utils.ParseStrNumber(s, &val); err != nil {
			return err
		}
		a.Val = val
	case protocol.E_VTYPE_DOUBLE:
		var val float64
		if err := utils.ParseStrNumber(s, &val); err != nil {
			return err
		}
		a.Val = val
	case protocol.E_VTYPE_STRING:
		fallthrough
	case protocol.E_VTYPE_WIDESTR:
		a.Val = s
	case protocol.E_VTYPE_OBJECT:
		var val uint64
		if err := utils.ParseStrNumber(s, &val); err != nil {
			return err
		}
		a.Val = val
	}
	return nil
}

func (a *Any) Add(other *Any) bool {
	if a.typ != other.typ {
		return false
	}

	switch a.typ {
	case protocol.E_VTYPE_INT:
		a.Val = a.Val.(int) + other.Val.(int)
	case protocol.E_VTYPE_INT64:
		a.Val = a.Val.(int64) + other.Val.(int64)
	case protocol.E_VTYPE_FLOAT:
		a.Val = a.Val.(float32) + other.Val.(float32)
	case protocol.E_VTYPE_DOUBLE:
		a.Val = a.Val.(float64) + other.Val.(float64)
	case protocol.E_VTYPE_STRING:
		a.Val = a.Val.(string) + other.Val.(string)
	default:
		return false
	}
	return true
}

func (a *Any) Clear() {

	switch a.typ {
	case protocol.E_VTYPE_INT:
		a.Val = 0
	case protocol.E_VTYPE_INT64:
		a.Val = int64(0)
	case protocol.E_VTYPE_FLOAT:
		a.Val = float32(0)
	case protocol.E_VTYPE_DOUBLE:
		a.Val = float64(0)
	case protocol.E_VTYPE_STRING:
		a.Val = ""
	case protocol.E_VTYPE_OBJECT:
		a.Val = uint64(0)
	}

}

type AttrList struct {
	Attr []*Any
	ki   map[string]int
}

func (a *AttrList) Print() {
	fmt.Println("attrs:")
	for k := range a.Attr {
		fmt.Println(a.Attr[k].Key, ":", a.Attr[k].Val)
	}
}
func NewAttrList() *AttrList {
	a := &AttrList{}
	a.Attr = make([]*Any, 0, 128)
	a.ki = make(map[string]int)
	return a
}

func (g *AttrList) Count() int {
	return len(g.Attr)
}

func (g *AttrList) Clear() {
	g.Attr = g.Attr[:0]
	g.ki = make(map[string]int)
}

func (g *AttrList) AddAttr(attr *Any) {
	idx := len(g.Attr)
	g.Attr = append(g.Attr, attr)
	g.ki[attr.Name()] = idx
}

func (g *AttrList) GetAttr(name string) *Any {
	if index, has := g.ki[name]; has {
		return g.Attr[index]
	}
	return nil
}

func (g *AttrList) GetAttrByIndex(index int) *Any {
	if index >= 0 && index < len(g.Attr) {
		return g.Attr[index]
	}
	return nil
}
