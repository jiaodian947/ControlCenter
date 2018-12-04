package gameobj

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

const (
	COMPRESSED_DATA_VERSION = 0x3050495A
	GLOBAL_DATA_VERSION = 0x30303031
	MAX_DATA_LEN            = 0x1000000
)
