package protocol

const (
	E_VTYPE_UNKNOWN  = iota // 未知
	E_VTYPE_BOOL            // 布尔
	E_VTYPE_BYTE            // 1字节
	E_VTYPE_WORD            // 2字节
	E_VTYPE_INT             // 32位整数
	E_VTYPE_INT64           // 64位整数
	E_VTYPE_FLOAT           // 单精度浮点数
	E_VTYPE_DOUBLE          // 双精度浮点数
	E_VTYPE_STRING          // 字符串
	E_VTYPE_WIDESTR         // 宽字符串
	E_VTYPE_OBJECT          // 对象号
	E_VTYPE_POINTER         // 指针
	E_VTYPE_USERDATA        // 用户数据
	E_VTYPE_MAX
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

func (m *VarMessage) Type(index int) int {
	if index >= m.Size {
		return E_VTYPE_UNKNOWN
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
	m.typ = append(m.typ, E_VTYPE_BOOL)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddByte(val int8) {
	m.Size++
	m.typ = append(m.typ, E_VTYPE_BYTE)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddWord(val int16) {
	m.Size++
	m.typ = append(m.typ, E_VTYPE_WORD)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddInt32(val int32) {
	m.Size++
	m.typ = append(m.typ, E_VTYPE_INT)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddInt64(val int64) {
	m.Size++
	m.typ = append(m.typ, E_VTYPE_INT64)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddFloat(val float32) {
	m.Size++
	m.typ = append(m.typ, E_VTYPE_FLOAT)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddDouble(val float64) {
	m.Size++
	m.typ = append(m.typ, E_VTYPE_DOUBLE)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddString(val string) {
	m.Size++
	m.typ = append(m.typ, E_VTYPE_STRING)
	m.data = append(m.data, val)
}

func (m *VarMessage) AddObject(val uint64) {
	m.Size++
	m.typ = append(m.typ, E_VTYPE_OBJECT)
	m.data = append(m.data, val)
}

func (m *VarMessage) BoolVal(index int) bool {
	if index >= m.Size {
		return false
	}

	switch val := m.data[index].(type) {
	case bool:
		return val
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
	case int32:
		return val
	default:
		return 0
	}
}

func (m *VarMessage) Int64Val(index int) int64 {
	if index >= m.Size {
		return 0
	}

	switch val := m.data[index].(type) {
	case int64:
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
	case float32:
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

func (m *VarMessage) ObjectVal(index int) uint64 {
	if index >= m.Size {
		return 0
	}

	switch val := m.data[index].(type) {
	case uint64:
		return val
	default:
		return 0
	}
}
