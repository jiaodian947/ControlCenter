package game

const (
	E_MODE_STOP   = iota // 停止
	E_MODE_MOTION        // 地表移动
	E_MODE_JUMP          // 跳跃
	E_MODE_JUMPTO        // 改变跳跃的目标方向
	E_MODE_FLY           // 空中移动
	E_MODE_SWIM          // 水中移动
	E_MODE_DRIFT         // 水面移动
	E_MODE_CLIMB         // 爬墙
	E_MODE_SLIDE         // 在不可行走范围内滑行
	E_MODE_SINK          // 下沉
	E_MODE_LOCATE        // 强制定位
)

const (
	E_DATAGRID_ADD_ROW     = iota // 添加行
	E_DATAGRID_REMOVE_ROW         // 删除行
	E_DATAGRID_CLEAR_ROW          // 清空
	E_DATAGRID_GRID_CHANGE        // 表元数据改变
	E_DATAGRID_SET_ROW            // 某一行的数据改变
)
