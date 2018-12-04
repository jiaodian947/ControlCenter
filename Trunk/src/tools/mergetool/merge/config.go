package merge

import (
	"encoding/xml"
	"io/ioutil"
)

const (
	M_INSERT = "insert"
	M_MASTER = "master"
	M_EMPTY  = "empty"
	M_CLEAR  = "clear"
	M_MERGE  = "merge"
	M_ADD    = "add"
)

type ThreadInfos struct {
	Threads    int `xml:"init_threads,attr"`
	MaxThreads int `xml:"max_threads,attr"`
	IncThreads int `xml:"inc_threads,attr"`
}

type DBInfo struct {
	Name       string `xml:"name,attr"`
	DataSource string `xml:"connect_string,attr"`
}

type Table struct {
	Name       string `xml:"name,attr"`
	Mode       string `xml:"mode,attr"`
	InsertCols string `xml:"insert_cols,attr"`
	Condition  string `xml:"condition,attr"`
}

type ConflictTableCol struct {
	Name      string `xml:"name,attr"`
	Value     string `xml:"value,attr"`
	Id        string `xml:"id,attr"`
	Func      string `xml:"func,attr"`
	Condition string `xml:"condition,attr"`
}

type ConflictTable struct {
	Name    string             `xml:"name,attr"`
	Refer   string             `xml:"refer,attr"`
	Columns []ConflictTableCol `xml:"Column"`
}

type GameData struct {
	Mode      string     `xml:"def_rec_mode"`
	GameAttrs []GameAttr `xml:"GameAttr"`
	GameRecs  []GameRec  `xml:"GameRec"`
	nameIdx   map[string]int
}

func (g *GameData) Prepare() {
	g.nameIdx = make(map[string]int)
	for k, v := range g.GameRecs {
		g.nameIdx[v.Name] = k
	}
}

func (g *GameData) GetRecByName(name string) *GameRec {
	if index, has := g.nameIdx[name]; has {
		return &g.GameRecs[index]
	}
	return nil
}

type GameAttr struct {
	Name        string `xml:"name,attr"`
	Mode        string `xml:"mode,attr"`
	Value       string `xml:"value,attr"`
	Discription string `xml:"discription,attr"`
}

type GameCol struct {
	Index       int    `xml:"index,attr"`
	Value       string `xml:"value,attr"`
	Discription string `xml:"discription,attr"`
}

type GameRec struct {
	Name   string    `xml:"name,attr"`
	Mode   string    `xml:"mode,attr"`
	Key    string    `xml:"key,attr"`
	Sort   string    `xml:"sort,attr"`
	MaxRow string    `xml:"MaxRow,attr"`
	Cols   []GameCol `xml:"GameCol"`
}

type GameObj struct {
	GameAttrs []GameAttr `xml:"GameAttr"`
	GameRecs  []GameRec  `xml:"GameRec"`
}

type ResolveTableCol struct {
	Name     string    `xml:"name,attr"`
	Value    string    `xml:"value,attr"`
	GameObj  *GameObj  `xml:"GameObj"`
	GameData *GameData `xml:"GameData"`
}

type ResolveTable struct {
	Name    string            `xml:"name,attr"`
	Columns []ResolveTableCol `xml:"Column"`
}

type MergeTableCol struct {
	Name     string    `xml:"name,attr"`
	Key      string    `xml:"key,attr"`
	GameData *GameData `xml:"GameData"`
}

type MergeTable struct {
	KeyName string          `xml:"keyname,attr"`
	Name    string          `xml:"name,attr"`
	Mode    string          `xml:"mode,attr"`
	Columns []MergeTableCol `xml:"Column"`

	keyIdx map[string]int
}

func (m *MergeTable) Prepare() {
	if m.KeyName == "" {
		m.KeyName = "s_name"
	}

	m.keyIdx = make(map[string]int)
	for k, v := range m.Columns {
		m.keyIdx[v.Key] = k
	}
}

func (m *MergeTable) FindKeyIndex(key string) int {
	if index, has := m.keyIdx[key]; has {
		return index
	}
	return -1
}

func (m *MergeTable) GetColumn(index int) *MergeTableCol {
	if index != -1 && index < len(m.Columns) {
		return &m.Columns[index]
	}
	return nil
}

type Sql struct {
	Sql string `xml:"sql,attr"`
}

type Config struct {
	Type        string          `xml:"type,attr"`
	Version     string          `xml:"version,attr"`
	Discription string          `xml:"discription,attr"`
	ThreadInfos ThreadInfos     `xml:"ThreadInfos"`
	DBInfos     []DBInfo        `xml:"DBInfos>DBInfo"`
	Tables      []Table         `xml:"Tables>Table"`
	Conflict    []ConflictTable `xml:"Conflict>Table"`
	Resolve     []ResolveTable  `xml:"Resolve>Table"`
	Merge       []MergeTable    `xml:"Merge>Table"`
	Sqls        []Sql           `xml:"SQLS>SQL"`
}

func NewConfig() *Config {
	c := &Config{}
	return c
}

func (c *Config) HasMergeTable(tbl string) bool {
	if tbl == "" {
		return false
	}
	for _, v := range c.Merge {
		if v.Name == tbl {
			return true
		}
	}
	return false
}

func (c *Config) Load(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	return xml.Unmarshal(data, c)
}
