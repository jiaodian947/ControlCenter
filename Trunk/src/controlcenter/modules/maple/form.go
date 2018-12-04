package maple

type ServerForm struct {
	Id           int
	GameId       int    `valid:"Required"`
	DistrictId   int    `valid:"Required"`
	ServerName   string `valid:"Required"`
	ServerType   int
	ServerStatus int
	MaxPlayer    int    `valid:"Required"`
	ServerIp     string `valid:"Required"`
	ServerPort   int    `valid:"Required"`
	Comment      string
}

type ServerEditForm struct {
	ServerId     int    `valid:"Required"`
	Id           int    `valid:"Required"`
	GameId       int    `valid:"Required"`
	DistrictId   int    `valid:"Required"`
	ServerName   string `valid:"Required"`
	ServerType   int
	ServerStatus int
	MaxPlayer    int    `valid:"Required"`
	ServerIp     string `valid:"Required"`
	ServerPort   int    `valid:"Required"`
	Comment      string
}

type GameForm struct {
	Id       int
	GameName string `valid:"Required"`
	Comment  string
}

type DistrictForm struct {
	Id           int
	DistrictName string `valid:"Required"`
	GameId       int    `valid:"Required"`
	Group        int
	Comment      string
}

type DistrictEditForm struct {
	DistrictId   int    `valid:"Required"`
	Id           int    `valid:"Required"`
	DistrictName string `valid:"Required"`
	GameId       int    `valid:"Required"`
	Group        int
	Comment      string
}
