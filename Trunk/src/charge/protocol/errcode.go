package protocol

const (
	CAS_ERR_ALREADY_DONE          = 53001 //交易已经完成了，重复验证
	CAS_ERR_VERIFY_FAILED         = 53002 //交易凭证无效
	CAS_ERR_VERIFY_TOO_OFTEN      = 53003 //同一个凭证验证太频繁
	CAS_ERR_BILL_NOT_MATCH        = 53004 //同一个记录基本信息不匹配
	CAS_ERR_BILL_ERROR            = 53005 //系统错误，请稍候重试
	CAS_ERR_TRANSACTION_NOT_FOUND = 53006 //交易没有找到
)
