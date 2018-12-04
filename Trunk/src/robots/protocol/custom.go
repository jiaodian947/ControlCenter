package protocol

const (
	CHATTYPE_SILENCE = 0
	// 可视范围聊天
	CHATTYPE_VISUALRANGE = 1
	// 普通场景聊天
	CHATTYPE_SCENE = 2
	// 私聊
	CHATTYPE_WHISPER = 3
	// 组队聊天
	CHATTYPE_TEAM = 4
	// 好友聊天
	CHATTYPE_FRIEND = 5
	// 公会聊天
	CHATTYPE_GUILD = 6
	// 联盟聊天
	CHATTYPE_UNION = 7
	// npc组织聊天
	CHATTYPE_ORGANIZATION = 8
	// 竞技频道
	CHATTYPE_SPORTS = 9
	// 战场
	CHATTYPE_BATTLE = 10
	// 世界聊天
	CHATTYPE_WORLD       = 11
	CHATTYPE_WORLD_FAILD = 12
	// 小喇叭
	CHATTYPE_SMALL_SPEAKER = 13
	// 离线时头顶信息
	CHATTYPE_OFFLINE = 14
	// 自动播放剧情台词
	CHATTYPE_MOVIE = 15
	// 团队聊天
	CHATTYPE_BIG_TEAM = 16
	// 喊话聊天
	CHATTYPE_SHOUT = 17
	// 队伍招募
	CHATTYPE_RECRUIT = 18

	// 私聊回显
	CHATTYPE_WHISPER_ECHO = 19
	// 对所有好友聊天回显
	CHATTYPE_FRIEND_ECHO = 20
	// 超链接发送
	CHATTYPE_HYPERLINK = 21
	// 系统聊天
	CHATTYPE_SYSTEM = 22
	// 电台聊天
	CHATTYPE_VIDEO = 23

	CHATTYPE_FRIEND_GROUP = 24

	CHATTYPE_TEAM_MSG = 25
)

const (
	//记录玩家消息到观察者客户端
	CLIENT_CUSTOMMSG_JOB_CHANGE = 1
	// 顶号情况下客户端通知服务器准备完毕
	CLIENT_CUSTOMMSG_ON_READY = 2
	// 防沉迷信息
	CLIENT_CUSTOMMSG_UNENTHRALL = 3
	//返回MessageBox的应答结果，消息格式：
	//unsigned msgid  int msgboxid  int answer(客户端点击的按纽索引)
	CLIENT_CUSTOMMSG_ANSWERMESSAGEBOX = 5
	//返回InputBox的应答结果，消息格式：
	//unsigned msgid  int inputboxid  wstring input(客户端输入的内容)
	CLIENT_CUSTOMMSG_ANSWERINPUTBOX = 6

	// 客户端剧情结束
	CLIENT_CUSTOMMSG_END_MOVIE = 7
	// 客户端请求剧情
	CLIENT_CUSTOMMSG_START_MOVIE = 8

	// 公会相关消息
	CLIENT_CUSTOMMSG_GUILD = 70

	// 客户端请求击杀公会Boss
	CLIENT_CUSTOMMSG_FIGHT_GUILD_BOSS = 71

	// 客户端请求开启公会Boss
	CLIENT_CUSTOMMSG_OPEN_GUILD_BOSS = 72

	// 客户端请求公会Boss信息
	CLIENT_CUSTOMMSG_GUILD_BOSS_INFO = 73

	// 客户端请求上交公会兽粮
	CLIENT_CUSTOMMSG_GUILD_BOSS_FOOD = 74

	/********************************************************
	* @brief   : 聊天消息
	* @details : ID  100~
	************************************************************************************************/
	//发送GM命令，格式：int msgid  wstring info
	CLIENT_CUSTOMMSG_GM = 100

	//聊天命令，格式：int msgid  int chat_type  wstring content(用空格分隔的一个宽字符串)  [wstring targetname](私聊专用：目标名称)
	CLIENT_CUSTOMMSG_CHAT = 110

	// 聊天室消息
	// 客户端请求创建聊天室： unsigned msgid  wstring wsTitleName  wstring wsPlayerName(邀请的玩家姓名)
	CLIENT_CREATE_CHATROOM = 120
	// 吧主请求解散聊天室 ：unsigned msgid
	CLIENT_DESTORY_CHATROOM = 121
	// 吧主踢人请求：unsigned msgid  wstring target
	CLIENT_KICK_CHATROOM = 122
	// 吧主提升或降低茶友权限: unsigned msgid  wstring target  unsigned postid
	CLIENT_UPLEVEL_CHATROOM = 123
	// 申请加入聊天室: unsigned msgid  int nChatRoomID
	CLIENT_JOIN_CHATROOM = 124
	// 邀请加入聊天室:unsigned msgid  wstring target
	CLIENT_INVITE_CHATROOM = 125
	// 申请退出聊天室: unsigned msgid  wstring target
	CLIENT_LEVEL_CHATROOM = 126
	// 发送聊天内容: unsigned msgid  wstring content
	CLIENT_CONTENT_CHATROOM = 127

	/********************************************************
	* @brief   : 场景消息
	* @details : ID  150~
	************************************************************************************************/
	//大地图使用道具传送
	CLIENT_CUSTOMMSG_BIG_MAP_TRANS = 150
	//通知服务器向客户端发送已激活传送点列表
	CLIENT_CUSTOMMSG_GET_ACTIVE_TRANS_POINTS = 151
	//更新场景种子数据
	CLIENT_CUSTOMMSG_EDITOR_UPDATAE_SCENE_OBJ = 155

	/********************************************************
	* @brief   : 战斗技能消息
	* @details : ID  200~
	************************************************************************************************/
	//使用技能，消息格式：unsigned msgid  AutoStr skillid  ...
	CLIENT_CUSTOMMSG_USE_SKILL = 200
	//打断吟唱和引导技能
	CLIENT_CUSTOMMSG_INTERRUPT_CURRENT_SKILL = 201
	// 特殊技能客户端判断的目标 消息格式：unsigned msgid  AutoStr skillid   ObjId target...
	CLIENT_CUSTOMMSG_SKILL_HIT_TARGET = 202
	// 特殊技能状态同步 消息格式：unsigned msgid  AutoStr skillid
	CLIENT_CUSTOMMSG_SKILL_SYCSTATE = 203

	/********************************************************
	* @brief   : 玩家消息
	* @details : ID  300~
	************************************************************************************************/
	//切换移动状态，格式：int msgid  int move_type(0:跑步  1:走路)
	CLIENT_CUSTOMMSG_MOVE_TYPE = 300
	// 移动同步消息int msgid move_mode param_num param1 param2 ...各种同步参数
	CLIENT_CUSTOMMSG_REQUEST_MOVE = 301
	// NPC移动同步消息int msgid  assist object  move_mode x z (就为了加伙伴协助技能而加的，奇奇怪怪的！！)
	CLIENT_CUSTOMMSG_ASSIST_NPC_MOVE = 302

	//复活命令，消息格式：unsigned msgid  int relive_type
	CLIENT_CUSTOMMSG_RELIVE = 310

	//选中返回人物选择界面
	CLIENT_CUSTOMMSG_CHECKED_SELECTROLE = 360

	CLIENT_CUSTOMMSG_OPTION = 361

	CLIENT_CUSTOMMSG_TO_BORN = 390
	// 组队系统用 151 ~ 180
	// 创建队伍
	CLIENT_CUSTOMMSG_TEAM = 400 // 组队消息

	// 组队副本
	CLIENT_CUSTOMMSG_MULTISCENE = 401

	// 爬塔
	CLIENT_CUSTOMMSG_TOWERSCENE = 402

	//请求接取转职的下一阶段任务
	CLIENT_CUSTOMMSG_TRANSFER_PROFESSION_NEXT_TASK = 420

	CLIENT_CUSTOMMSG_TRANSFER_PROFESSION_COMPLETE = 421

	// 接收邀请
	CLIENT_CUSTOMMSG_TEAM_INVITE_ACCEPT = 422

	// 清空邀请列表
	CLIENT_CUSTOMMSG_TEAM_INVITE_CLEAR = 423

	// 催促
	CLIENT_CUSTOMMSG_TEAM_HURRY_UP = 424

	// 计费商城
	CLIENT_CUSTOMMSG_CHARG_SHOP = 461

	// 充值相关
	CLIENT_CUSTOMMSG_PAY = 462

	/********************************************************
	* @brief   : 通用消息
	* @details : ID  500~
	************************************************************************************************/
	//用户交互请求，消息格式：unsigned msgid  int request_type  wstring name(被请求者角色名称)
	CLIENT_CUSTOMMSG_REQUEST = 500
	//用户交互请求回应，消息格式：unsigned msgid  int request_type  wstring name(被请求者角色名称)  int result(0:拒绝 1:同意 2:超时)
	CLIENT_CUSTOMMSG_ANSWER     = 501
	CLIENT_CUSTOMMSG_ANSWER_ALL = 502
	//用户设置对其他用户的请求开关，消息格式：unsigned msgid  int newvalue(每位对应request_type的一种类型);
	CLIENT_CUSTOMMSG_ANSWER_RESULT = 503

	CLIENT_CUSTOMMSG_MALL = 508 // 商城消息

	CLIENT_CUSTOMMSG_INFO = 509 // 信息消息

	// 伙伴消息
	CLIENT_CUSTOMMSG_PARTNER = 510

	// 好友消息
	CLIENT_CUSTOMMSG_FRIEND = 511

	// 章回体副本消息
	CLINET_CUSTOMMSG_CLONE_SCENE = 512

	// 好友推荐
	CLIENT_CUSTOMMSG_RECOMMEND = 513
	// 技能界面相关消息
	CLIENT_CUSTOMMSG_SKILL = 514

	/*!
	* @brief	使用技能时 对象产生了位移
	* @param	int64		技能的uid
	* @param	string		技能的id
	* @param	ObjId	位移的目标
	* @param	float x z orient
	 */
	CLIENT_CUSTOMMSG_SKILL_LOCATE_OBJECT = 515

	// 道具相关消息
	CLIENT_CUSTOMMSG_TOOL_ITEM = 516

	// 装备消息
	CLIENT_CUSTOMMSG_EQUIPMENT = 517

	// 委托消息
	CLIENT_CUSTOMMSG_DELEGATION = 518

	// 客户端主动切换场景
	CLIENT_CUSTOMMSG_SWITCH_SCENE = 519

	// 兑换
	CLIENT_CUSTOMMSG_EXCHANGE = 520

	// 查询服务器时间
	CLIENT_CUSTOMMSG_GET_SERVERTIME = 521

	// 客户端拾取掉落物品
	CLIENT_CUSTOMMSG_PICK_UP_DROP_ITEM = 522

	// 客户端触发器消息
	CLIENT_CUSTOMMSG_TRIGGER_EVENT = 523

	//容器操作相关命令
	//客户端上传使用物品的消息，消息格式：unsigned msgid int srcviewid int srcpos int dstviewid int dstpos
	CLIENT_CUSTOMMSG_USEITEM_ON_ITEM = 529
	//客户端上传移动物品的消息，消息格式：unsigned msgid  int srcviewid  int srcpos  int destviewid  int destpos
	CLIENT_CUSTOMMSG_MOVEITEM = 530
	//客户端上传使用物品的消息，消息格式：unsigned msgid  int srcviewid  int srcpos
	CLIENT_CUSTOMMSG_USEITEM = 531
	//客户端上传丢弃物品的消息，消息格式：unsigned msgid  int srcviewid  int srcpos  int amount(数量，为0时丢弃全部物品)
	CLIENT_CUSTOMMSG_CHUCKITEM = 532
	//客户端上传拆分物品的消息，消息格式：unsigned msgid  int srcviewid  int srcpos  int destviewid  int destpos  int amount
	CLIENT_CUSTOMMSG_SPLITITEM = 533
	//客户端上传删除物品的消息，消息格式：unsigned msgid  int srcviewid  int srcpos  int amount(数量，为0时丢弃全部物品)
	CLIENT_CUSTOMMSG_DELETEITEM = 534
	//客户端上传锁定物品的消息，消息格式：unsigned msgid  int srcviewid  int srcpos
	CLIENT_CUSTOMMSG_LOCKITEM = 535
	//客户端上传整理容器中的物品，消息格式：unsigned msgid  int srcviewid  int beginpos  int endpos
	CLIENT_CUSTOMMSG_ARANGEITEM = 536
	//一键出售杂物
	CLIENT_CUSTOMMSG_SELL_SUNDRIES = 550

	/********************************************************
	* @brief   : 玩法消息
	* @details : ID  650~
	************************************************************************************************/
	//邮件模块
	CLIENT_CUSTOMMSG_POST = 450

	//通报系统
	CLIENT_CUSTOMMSG_NOTIFY = 460

	// 任务链相关消息
	CLIENT_CUSTOMMSG_QUEST_LOOP = 879

	//任务相关消息 消息格式：unsigned msgid  int submsgid  int questid
	CLIENT_CUSTOMMSG_QUEST_MSG = 880

	//采集任务相关消息 消息格式
	CLIENT_CUSTOMMSG_GATHER_MSG = 881

	// 勾玉模块消息 消息格式 一级消息，二级消息（后面具体的查看具体的二级消息枚举注释）
	CLIENT_CUSTOMMSG_MAGATAMA_MSG = 882

	// 忍者竞技模块消息
	CLIENT_CUSTOMMSG_NINJA_ARENA_MSG = 883

	// 伙伴抽取
	CLIENT_CUSTOMMSG_PARTNER_PRIZE_MSG = 884

	// 查看别人的信息
	CLIENT_CUSTOMMSG_QUERY_PLAYER_INFO_MSG = 885
	// 切磋
	CLIENT_CUSTOMMSG_PK_MSG = 886

	CLIENT_CUSTOMMSG_GROUP_PK_MSG = 890
	// 组队PK消息
	CLIENT_CUSTOMMSG_TEAM_PK_MSG = 894
	/********************************************************
	* @brief   : 综合消息
	* @details : ID  900~
	************************************************************************************************/
	//世界排名消息
	CLIENT_CUSTOMMSG_WORLD_RANK = 963

	//////////////////////////////////////////////////////////////////////////
	CLIENT_CUSTOMMSG_DEL_ALL_LETTER = 970 // 删除所有邮件

	CLIENT_CUSTOMMSG_READ_LETTER = 971 // 读取邮件

	CLIENT_CUSTOMMSG_DEL_LETTER = 972 // 删除邮件

	CLIENT_CUSTOMMSG_GET_LETTER_APPENDIX = 973 // 获取邮件附件

	CLIENT_CUSTOMMSG_OPEN_MAIL_UI_PANEL = 974

	// gmcc 消息
	CLIENT_CUSTOMMSG_GMCC_MSG = 1020
	//后台模块消息
	CLIENT_CUSTOMMSG_BACKSTAGE = 1021

	/********************************************************
	* @brief   : 测试专用消息
	* @details : ID  1023~
	************************************************************************************************/
	//压力测试专用消息int msgid AutoStr test_type..各种测试对应参数
	CLIENT_CUSTOMMSG_STRESS_TEST = 1023

	//楼下的注意 这个编号在 1-1023 之内 超过这个编号注册回调会失败
	CLIENT_CUSTOMMSG_MAX = 1024
)

const (
	//系统消息，格式：int msgid  int tipstype  string stringid  ...(参数表)
	SERVER_CUSTOMMSG_SYSINFO = 1
	//服务器同步时间，格式：int msgid  int64 server_time
	SERVER_CUSTOMMSG_ASYN_TIME = 2

	//通知客户端弹出一个MessageBox，消息格式：
	//unsigned msgid  int msgboxid  int count(按纽个数)  ...(按纽的提示信息文本ID)  string stringid(提示信息文本ID)  ...(提示信息参数)
	SERVER_CUSTOMMSG_SHOWMESSAGEBOX = 3
	//通知客户端弹出一个InputBox，消息格式：
	//unsigned msgid  int inputboxid  int type(0:只能数值 1:文本)  string stringid(提示信息文本ID)  ...(提示信息参数)
	SERVER_CUSTOMMSG_SHOWINPUTBOX = 4

	// 播放剧情
	SERVER_CUSTOMMSG_PLAY_SCENARIO = 5
	// 停止剧情
	SERVER_CUSTOMMSG_STOP_SCENARIO = 6
	// 暂停剧情
	SERVER_CUSTOMMSG_PAUSE_SCENARIO = 7
	// 继续剧情
	SERVER_CUSTOMMSG_CONTINUE_SCENARIO = 8
	// 系统消息 扩展
	SERVER_CUSTOMMSG_SYSINFO_EX = 19

	// 防沉迷信息
	SERVER_CUSTOMMSG_UNENTHRALL = 30
	// 发布新闻公告
	SERVERMSG_PUSLISH_NEWS = 31
	//notify
	SERVER_CUSTOMMSG_NOTIFY = 32
	// 公会相关下发消息主消息 int sub_id  ...
	SERVER_CUSTOMMSG_GUILD = 34

	// 伙伴消息
	SERVER_CUSTOMMSG_PARTNER = 40

	// 公会BOSS信息
	SERVER_CUSTOMMSG_GUILD_BOSS_INFO = 45

	// 好友消息
	SERVER_CUSTOMMSG_FRIEND = 50

	// 虚拟表用 61 ~ 70
	SERVER_CUSTOMMSG_VIRTUAL_RECORD_ADD = 60

	SERVER_CUSTOMMSG_VIRTUAL_RECORD_CLEAR = 61

	SERVER_CUSTOMMSG_VIRTUAL_RECORD_FINISH = 62

	// 组队系统用 110 ~ 130
	// 创建队伍
	SERVER_CUSTOMMSG_TEAM_CREATE = 111

	// 加入队伍
	SERVER_CUSTOMMSG_TEAM_JOIN = 112

	// 离开队伍
	SERVER_CUSTOMMSG_TEAM_LEAVE = 113

	// 解散队伍
	SERVER_CUSTOMMSG_TEAM_DISBAND = 114

	// 请求队伍
	SERVER_CUSTOMMSG_TEAM_REQUEST = 115

	// 踢出队伍
	SERVER_CUSTOMMSG_TEAM_KICKOFF = 116

	// 开始招募
	SERVER_CUSTOMMSG_TEAM_RECRUIT_START = 117

	// 结束招募
	SERVER_CUSTOMMSG_TEAM_RECRUIT_FINISH = 118

	// 开始匹配
	SERVER_CUSTOMMSG_TEAM_MATCH_START = 119

	// 结束匹配
	SERVER_CUSTOMMSG_TEAM_MATCH_FINISH = 120

	// 设置是否申请是否立即入队
	SERVER_CUSTOMMSG_TEAM_SET_IMMEDIATELY = 121

	// 转让队长
	SERVER_CUSTOMMSG_TEAM_TRANSFER = 122

	// 匹配结束
	SERVER_CUSTOMMSG_TEAM_MATCHTIMEOUT = 123

	// 多人副本
	SERVER_CUSTOMMSG_MULTISCENE = 124

	//
	SERVER_CUSTOMMSG_TOWERSCENE = 125

	SERVER_CUSTOMMSG_TEAM_CONVEY        = 126
	SERVER_CUSTOMMSG_TEAM_CONVEY_FINISH = 127
	SERVER_CUSTOMMSG_TEAM_FALLIN_WALK   = 128

	/********************************************************
	* @brief   : 聊天消息
	* @details : ID  150~
	************************************************************************************************/
	//聊天消息
	SERVER_CUSTOMMSG_CHAT = 150
	//聊天界面接收物品超链接信息
	SERVER_CUSTOMMSG_CHATHYPER = 151

	// 聊天室信息
	// 进入聊天室
	SERVER_CUSTOMMSG_ENTER_CHATROOM = 185
	// 离开聊天室
	SERVER_CUSTOMMSG_LEAVE_CHATROOM = 186
	// 更新聊天室聊天内容 int msgid  wstring wsContent
	SERVER_CUSTOMMSG_UPDATE_CONTENT = 187
	// 发送聊天室系统提示信息  int msgid  int tipstype(0，表示0个参数，1表示1个参数，2表示2个参数)  string stringid  ...(参数表)
	SERVER_CUSTOMMSG_CHATROOM_SYSINFO = 188

	/********************************************************
	* @brief   : 消息
	* @details : ID  200~
	************************************************************************************************/
	//用户交互请求，消息格式：unsigned msgid  int request_type  wstring name(请求者角色名称)  ...(参数，根据request_type的不同而不同)
	SERVER_CUSTOMMSG_REQUEST = 208

	// 播放一次性动作
	SERVER_CUSTOMMSG_PLAYER_ACTION = 210
	//通知客户端马上要进行切场景，可以准备进行预处理unsigned msgid  float x  float y  float z
	SERVER_CUSTOMMSG_SWITCH_SCENE_BEGIN = 211
	//通知客户端玩家进入场景准备就绪
	SERVER_CUSTOMMSG_PLAYER_ONREADY = 212

	//通知客户端玩家继续连接
	SERVER_CUSTOMMSG_PLAYER_ONCONTINUE = 213

	// 通知客户端开始一段剧情  参数格式: string cmd_id  Object NPC  int movie_id
	SERVER_CUSTOMMSG_START_MOVIE = 214

	//通知客户端有一段剧情等待播放，参数格式：string cmd_id  Object NPC  int movie_id
	SERVER_CUSTOMMSG_MOVIE_REQUEST = 215

	// 通知客户端显示副本倒计时 string msg， int time  int type(CloneSceneModule.h中有定义枚举)
	SERVER_CUSTOMMSG_PREPARE_EXIT_CLONE = 216

	//客户端显示副本奖励确认对话框
	SERVER_CUSTOMMSG_CLONE_COUNT = 217

	// 通知客户端显示或关闭副本相关界面    int flag(0标示关闭  1标示显示)
	SERVER_CUSTOMMSG_SET_CLONE_FORM = 219

	//通知客户端获得奖励 消息格式; unsigned msgid，int award_id
	SERVER_CUSTOMMSG_REAP_AWARD = 220

	// 允许（禁止）移动
	SERVER_CUSTOMMSG_ALLOW_MOVE = 221

	// 通知客户端玩家复活
	SERVER_CUSTOMMSG_RELIVE_END = 222

	// 通知客户端弹复活界面
	SERVER_CUSTOMMSG_NOTIFY_DIE = 223

	// 通知客户端玩家死亡
	SERVER_CUSTOMMSG_PLAYER_DIE = 224

	//自定义特效，消息格式：unsigned msgid  object target  ...(由客户端程序确定)
	SERVER_CUSTOMMSG_EFFECT = 301

	//开始念咒，消息格式：unsigned msgid  int tick(时间，单位ms)
	SERVER_CUSTOMMSG_BEGIN_CURSE = 310

	//结束念咒，消息格式：unsigned msgid
	SERVER_CUSTOMMSG_END_CURSE = 311

	// 任务系统主消息 string msg， int subType.....
	SERVER_CUSTOMMSG_QUEST_MSG = 312

	SERVER_CUSTOMMSG_SYSTEM_MAIL = 313

	// 任务链主消息
	SERVER_CUSTOMMSG_LOOPQUEST_MSG = 314

	// 设置NPC碰撞
	SERVER_CUSTOMMSG_SET_NEED_COLLIDE = 350
	// 预加载npc相关消息
	SERVER_CUSTOMMSG_PRELOAD_NPC = 351
	// 通知客户端NPC说话: string msg_id  object npc  string talk_id
	SERVER_CUSTOMMSG_NPC_TALK = 352

	// 章回体副本消息
	SERVER_CUSTOMMSG_CLONE_SCENE = 353

	// 通知客户端掉落
	SERVER_CUSTOMMSG_ITEM_DROP = 354

	// 同步客户端时间
	SERVER_CUSTOMMSG_SERVER_TIME = 355

	// 触发事件
	SERVER_CUSTOMMSG_TRIGGER_EVENT = 356

	// 忍者的委托
	SERVER_CUSTOMMSG_DAILY_ACTIVITY = 357

	// 忍者竞技
	SERVER_CUSTOMMSG_NINJA_ARENA = 358

	// 勾玉
	SERVER_CUSTOMMSG_MAGATAMA = 359

	// 伙伴抽取（十连抽）
	SERVER_CUSTOMMSG_PARTNER_DRAW_PRIZE = 360

	// 切磋
	SERVER_CUSTOMMSG_PK = 361

	// 世界排行榜
	SERVER_CUSTOMMSG_WORLD_RANK = 400

	// 对象原地死亡: string msg_id  object npc
	SERVER_CUSTOMMSG_SITU_DEAD = 410

	SERVER_CUSTOMMSG_CHARG_SHOP = 450

	// 查看人物信息
	SERVER_CUSTOMMSG_QUERY_PLAYER_INFO = 451

	// 分组PK相关消息
	SERVER_CUSTOMMSG_GROUP_PK = 454

	// 跨服帮助场景相关消息
	SERVER_CUSTOMMSG_CROSS_HELP = 457

	/********************************************************
	* @brief   : 技能消息
	* @details : ID  500~
	************************************************************************************************/
	// 技能冷却时间
	SERVER_CUSTOMMSG_SKILL_CD_TIME = 700
	// 技能相关 UUID  skill_id  skill_stage  params
	SERVER_CUSTOMMSG_SKILL = 701

	//技能开始释放string skillid  object target
	SERVER_CUSTOMMSG_SKILL_SUCCESS = 702

	// 技能释放失败消息 string skillid
	SERVER_CUSTOMMSG_SKILL_FAIL = 703

	/*!
	* @brief	技能效果产生了Motion
	* @param	ObjId 目标
	* @param	float x
	* @param   float y
	* @param	float z
	* @param	float speed
	 */
	SERVER_CUSTOMMSG_SKILL_MOTION = 704

	// 自定义消息通知客户端切换坐标
	SERVER_CUSTOMMSG_SWITCH_LOCATION = 705

	// 服务器返回领取附件成功get_apdix
	SERVER_CUSTOM_GET_APDIX_SUCCESS = 710

	SERVER_CUSTOMMSG_CHECK_CHAT_GUILD = 713

	// 移动相关
	SERVER_CUSTOMMSG_MOVE     = 720
	SERVER_CUSTOMMSG_RELOCATE = 721
	SERVER_CUSTOMMSG_ROTATE   = 722

	/********************************************************
	* @brief   : GMCC消息
	* @details : ID  900~
	************************************************************************************************/
	//Gm记录某个玩家的消息
	SERVER_CUSTOMMSG_MSG_LOG = 900

	//楼下的注意 这个编号在 1-1023 之内 超过这个编号注册回调会失败
	SERVER_CUSTOMMSG_MAX = 1024
)
