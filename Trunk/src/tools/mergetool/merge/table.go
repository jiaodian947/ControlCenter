package merge

import (
	"database/sql"
	"fmt"
	"sort"
)

// 数据库表字段信息
type FieldInfo struct {
	Field   string         //字段名
	Type    string         //类型
	Null    string         //是否可以为空
	Key     string         //主键类型
	Default sql.NullString //默认值
	Extra   string
}

// 判断两个字段是否相同
func (f *FieldInfo) Equal(other FieldInfo) bool {
	return f.Field == other.Field &&
		f.Type == other.Type &&
		f.Key == other.Key
}

type Fields []FieldInfo

func (f Fields) Len() int {
	return len(f)
}

func (f Fields) Less(i, j int) bool {
	return f[i].Field < f[j].Field
}

func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// 数据库表信息
type TableInfo struct {
	Table  string //表名
	Create string //建表sql
	Fields Fields //字段信息
	Keys   []string
}

// 获取建表sql
func (t *TableInfo) LoadCreateSql(db *sql.DB) error {
	r, err := db.Query(fmt.Sprintf("SHOW CREATE TABLE `%s`", t.Table))
	if err != nil {
		return err
	}
	defer r.Close()
	if !r.Next() {
		return fmt.Errorf("load create sql error")
	}

	var table, create string
	if err := r.Scan(&table, &create); err != nil {
		return err
	}

	t.Create = create
	return nil
}

// 加载表格信息
func (t *TableInfo) LoadInfo(db *sql.DB) error {
	r, err := db.Query(fmt.Sprintf("SHOW COLUMNS FROM `%s`", t.Table))
	if err != nil {
		return fmt.Errorf("table not found")
	}

	defer r.Close()
	t.Fields = make([]FieldInfo, 0, 32)
	t.Keys = make([]string, 0, 2)
	for r.Next() {
		ci := FieldInfo{}
		if err := r.Scan(&ci.Field, &ci.Type, &ci.Null, &ci.Key, &ci.Default, &ci.Extra); err != nil {
			panic(err)
		}
		t.Fields = append(t.Fields, ci)
		if ci.Key == "PRI" || ci.Key == "UNI" {
			t.Keys = append(t.Keys, ci.Field)
		}
	}
	sort.Sort(t.Fields)
	return nil
}

// 获取一个主键
func (t *TableInfo) OneKey() string {
	if len(t.Keys) > 0 {
		return t.Keys[0]
	}
	return ""
}

func NewTableInfo(db *sql.DB, tbl string, loadcreate bool) (*TableInfo, error) {
	ti := &TableInfo{}
	ti.Table = tbl
	if err := ti.LoadInfo(db); err != nil {
		return nil, err
	}
	if loadcreate {
		if err := ti.LoadCreateSql(db); err != nil {
			return nil, err
		}
	}
	return ti, nil
}

// 获取表格行数量
func TableRows(db *sql.DB, table string) int {
	r, err := db.Query(fmt.Sprintf("SELECT COUNT(1) FROM `%s`", table))
	if err != nil {
		panic(err)
	}
	defer r.Close()
	if r.Next() {
		var count int
		r.Scan(&count)
		return count
	}
	return 0
}
