package protocol

type PERSISTID struct {
	Ident  uint32
	Serial uint32
}

const (
	VTYPE_UNKNOWN  = iota // 未知
	VTYPE_BOOL            // 布尔
	VTYPE_BYTE            // 1字节
	VTYPE_WORD            // 2字节
	VTYPE_INT             // 32位整数
	VTYPE_INT64           // 64位整数
	VTYPE_FLOAT           // 单精度浮点数
	VTYPE_DOUBLE          // 双精度浮点数
	VTYPE_STRING          // 字符串
	VTYPE_WIDESTR         // 宽字符串
	VTYPE_OBJECT          // 对象号
	VTYPE_POINTER         // 指针
	VTYPE_USERDATA        // 用户数据
	VTYPE_MAX
)

type VarMessage struct {
	ConnId int64
	Serial int64
	Size   int
	typ    []int
	data   []interface{}
}

func NewVarMsg(cap int) *VarMessage {
	v := &VarMessage{}
	if cap == 0 {
		cap = 8
	}
	v.typ = make([]int, 0, cap)
	v.data = make([]interface{}, 0, cap)
	return v
}

func (m *VarMessage) Clear() {
	m.Size = 0
	m.typ = m.typ[:0]
	m.data = m.data[:0]
}

func (m *VarMessage) Count() int {
	return len(m.data)
}

func (m *VarMessage) Type(index int) int {
	if index >= m.Size {
		return VTYPE_UNKNOWN
	}

	return m.typ[index]
}

func (m *VarMessage) RawValue(index int) interface{} {
	if index >= m.Size {
		return nil
	}

	return m.data[index]
}

func (m *VarMessage) AddBool(val bool) {
	m.Size++
	m.typ = append(m.typ, VTYPE_BOOL)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddByte(val int8) {
	m.Size++
	m.typ = append(m.typ, VTYPE_BYTE)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddWord(val int16) {
	m.Size++
	m.typ = append(m.typ, VTYPE_WORD)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddInt32(val int32) {
	m.Size++
	m.typ = append(m.typ, VTYPE_INT)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddInt64(val int64) {
	m.Size++
	m.typ = append(m.typ, VTYPE_INT64)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddFloat(val float32) {
	m.Size++
	m.typ = append(m.typ, VTYPE_FLOAT)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddDouble(val float64) {
	m.Size++
	m.typ = append(m.typ, VTYPE_DOUBLE)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddString(val string) {
	m.Size++
	m.typ = append(m.typ, VTYPE_STRING)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddObject(val PERSISTID) {
	m.Size++
	m.typ = append(m.typ, VTYPE_OBJECT)
	m.data = append(m.data, val)
}

func (m *VarMessage) BoolVal(index int) bool {
	if index >= m.Size {
		return false
	}

	switch val := m.data[index].(type) {
	case bool:
		return val
	case int:
		return val != 0
	case int64:
		return val != 0
	default:
		return false
	}
}

func (m *VarMessage) ByteVal(index int) int8 {
	if index >= m.Size {
		return 0
	}

	switch val := m.data[index].(type) {
	case int8:
		return val
	default:
		return 0
	}
}

func (m *VarMessage) WordVal(index int) int16 {
	if index >= m.Size {
		return 0
	}

	switch val := m.data[index].(type) {
	case int16:
		return val
	default:
		return 0
	}
}

func (m *VarMessage) Int32Val(index int) int32 {
	if index >= m.Size {
		return 0
	}

	switch val := m.data[index].(type) {
	case bool:
		if val {
			return 1
		}
		return 0
	case int32:
		return int32(val)
	case int64:
		return int32(val)
	case float32:
		return int32(val)
	case float64:
		return int32(val)
	default:
		return 0
	}
}

func (m *VarMessage) Int64Val(index int) int64 {
	if index >= m.Size {
		return 0
	}

	switch val := m.data[index].(type) {
	case bool:
		if val {
			return 1
		}
		return 0
	case int:
		return int64(val)
	case int64:
		return int64(val)
	case float32:
		return int64(val)
	case float64:
		return int64(val)
	default:
		return 0
	}
}

func (m *VarMessage) DoubleVal(index int) float64 {
	if index >= m.Size {
		return 0
	}

	switch val := m.data[index].(type) {
	case bool:
		if val {
			return 1
		}
		return 0
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return float64(val)
	default:
		return 0
	}
}

func (m *VarMessage) FloatVal(index int) float32 {
	if index >= m.Size {
		return 0
	}

	switch val := m.data[index].(type) {
	case bool:
		if val {
			return 1
		}
		return 0
	case int:
		return float32(val)
	case int64:
		return float32(val)
	case float32:
		return float32(val)
	case float64:
		return float32(val)
	default:
		return 0
	}
}

func (m *VarMessage) StringVal(index int) string {
	if index >= m.Size {
		return ""
	}

	switch val := m.data[index].(type) {
	case string:
		return val
	default:
		return ""
	}
}

func (m *VarMessage) ObjectVal(index int) PERSISTID {
	if index >= m.Size {
		return PERSISTID{}
	}

	switch val := m.data[index].(type) {
	case PERSISTID:
		return val
	default:
		return PERSISTID{}
	}
}
