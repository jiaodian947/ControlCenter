package protocol

const (
	CTOS_MSG_BEGIN    = 0
	CTOS_LOGIN        = 1
	CTOS_CREATE_ROLE  = 2
	CTOS_DELETE_ROLE  = 3
	CTOS_CHOOSE_ROLE  = 4
	CTOS_WORLD_INFO   = 5
	CTOS_READY        = 6
	CTOS_CUSTOM       = 7
	CTOS_REQUEST_MOVE = 8
	CTOS_SELECT       = 9
	CTOS_SPEECH       = 10
	CTOS_GET_VERIFY   = 11
	CTOS_RET_ENCODE   = 12
	CTOS_MSG_END      = 99

	STOC_MSG_BEGIN          = 100
	STOC_LOGIN_SUCCEED      = 101
	STOC_PROPERTY_TABLE     = 102
	STOC_RECORD_TABLE       = 103
	STOC_ENTRY_SCENE        = 104
	STOC_EXIT_SCENE         = 105
	STOC_ADD_OBJECT         = 106
	STOC_REMOVE_OBJECT      = 107
	STOC_SCENE_PROPERTY     = 108
	STOC_OBJECT_PROPERTY    = 109
	STOC_CREATE_VIEW        = 110
	STOC_DELETE_VIEW        = 111
	STOC_VIEW_PROPERTY      = 112
	STOC_VIEW_ADD           = 113
	STOC_VIEW_REMOVE        = 114
	STOC_VIEW_CHANGE        = 115
	STOC_RECORD_ADD_ROW     = 116
	STOC_RRECORD_DEL_ROW    = 117
	STOC_RECORD_GRID        = 118
	STOC_RECORD_CLEAR       = 119
	STOC_SPEECH             = 120
	STOC_SYSTEM_INFO        = 121
	STOC_MENU               = 122
	STOC_CLEAR_MENU         = 123
	STOC_CUSTOM             = 124
	STOC_LOCATION           = 125
	STOC_MOVING             = 126
	STOC_ALL_DEST           = 127
	STOC_WARNING            = 128
	STOC_FROM_GMCC          = 129
	STOC_ALL_PROP           = 130
	STOC_ADD_MORE_OBJECT    = 131
	STOC_REMOVE_MORE_OBJECT = 132
	STOC_SERVER_INFO        = 133
	STOC_SET_VERIFY         = 134
	STOC_SET_ENCODE         = 135
	STOC_ERROR_CODE         = 136
	STOC_WORLD_INFO         = 137
	STOC_IDLE               = 138
	STOC_QUEUE              = 139
	STOC_TERMINATE          = 140
	STOC_LINK_TO            = 141
	STOC_UNLINK             = 142
	STOC_LINK_MOVE          = 143
	STOC_CREATE_BLKING      = 144
	STOC_ROLE_TRANSFERED    = 145
	STOC_LOGIN_STRING       = 146
	STOC_MSG_END            = 199

	STOC_COMPRESSED_MSG_BEGIN = 200
	STOC_CP_ADD_OBJECT        = 201 // 压缩的添加可见对象消息
	STOC_CP_RECORD_ADD_ROW    = 202 // 压缩的表格添加行消息
	STOC_CP_VIEW_ADD          = 203 // 压缩的容器添加对象消息
	STOC_CP_CUSTOM            = 204 // 压缩的自定义消息
	STOC_CP_ALL_DEST          = 205 // 压缩的对象移动消息
	STOC_CP_ALL_PROP          = 206 // 压缩的多个对象的属性改变信息
	STOC_CP_ADD_MORE_OBJECT   = 207 // 压缩的增加多个对象
	STOC_COMPRESSED_MSG_END   = 299
)

const (
	SC_TYPE_UNKNOWN = iota // 未知
	SC_TYPE_BYTE           // 一字节
	SC_TYPE_WORD           // 二字节
	SC_TYPE_DWORD          // 四字节
	SC_TYPE_QWORD          // 八字节
	SC_TYPE_FLOAT          // 浮点四字节
	SC_TYPE_DOUBLE         // 浮点八字节
	SC_TYPE_STRING         // 字符串，前四个字节为长度
	SC_TYPE_WIDESTR        // UNICODE宽字符串，前四个字节为长度
	SC_TYPE_OBJECT         // 对象号
	SC_TYPE_MAX
)
