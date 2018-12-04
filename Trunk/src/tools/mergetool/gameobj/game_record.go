package gameobj

import (
	"sort"
)

type Row struct {
	Values []*Any
}

func NewRow(cols int) *Row {
	r := &Row{}
	r.Values = make([]*Any, cols)
	return r
}

func (r *Row) SetValue(col int, value *Any) {
	if col >= len(r.Values) {
		panic("col exceed")
	}

	r.Values[col] = value
}

func (r *Row) Value(col int) *Any {
	if col >= len(r.Values) {
		panic("col exceed")
	}

	return r.Values[col]
}

type Column struct {
	ColType []int
}

func NewColumn(cols int) *Column {
	c := &Column{}
	c.ColType = make([]int, cols)
	return c
}

func (c *Column) SetType(col int, typ int) {
	if col >= len(c.ColType) {
		panic("col exceed")
	}
	c.ColType[col] = typ
}

func (c *Column) Type(col int) int {
	if col >= len(c.ColType) {
		panic("col exceed")
	}
	return c.ColType[col]
}

type KeyInfo struct {
	Index           int32
	CaseInsensitive int8
}

type Record struct {
	RecName  string
	Columns  *Column
	Rows     []*Row
	cols     int
	MaxRow   int
	sort_key int
	keys     []KeyInfo
}

func NewRecord(cols int, maxrows int) *Record {
	r := &Record{}
	r.Columns = NewColumn(cols)
	r.cols = cols
	if maxrows == 0 {
		maxrows = 256
	}
	r.MaxRow = maxrows
	r.Rows = make([]*Row, 0, maxrows)
	r.keys = make([]KeyInfo, 0, 256)
	return r
}

func (r *Record) Keys() int {
	return len(r.keys)
}

func (r *Record) ClearKeys() {
	r.keys = r.keys[:0]
}

func (r *Record) AddKey(index int32, case_insensitive int8) {
	r.keys = append(r.keys, KeyInfo{index, case_insensitive})
}

func (r *Record) Key(index int) *KeyInfo {
	if index >= 0 && index < len(r.keys) {
		return &r.keys[index]
	}

	return nil
}

func (r *Record) Clear() {
	r.Rows = r.Rows[:0]
}

func (r *Record) Column() *Column {
	return r.Columns
}

func (r *Record) SetName(name string) {
	r.RecName = name
}

func (r *Record) Name() string {
	return r.RecName
}

func (r *Record) MaxRows() int {
	return r.MaxRow
}

func (r *Record) Cols() int {
	return r.cols
}

func (r *Record) RowCount() int {
	return len(r.Rows)
}

func (r *Record) AddRow(row int) int {
	newrow := NewRow(r.cols)
	if row <= -1 || row >= len(r.Rows) {
		r.Rows = append(r.Rows, newrow)
		return len(r.Rows) - 1
	}

	r.Rows = append(r.Rows, newrow)
	copy(r.Rows[row+1:], r.Rows[row:])
	r.Rows[row] = newrow
	return row
}

func (r *Record) AddRowValue(row int, newrow *Row) int {
	if row <= -1 || row >= len(r.Rows) {
		r.Rows = append(r.Rows, newrow)
		return len(r.Rows) - 1
	}

	r.Rows = append(r.Rows, newrow)
	copy(r.Rows[row+1:], r.Rows[row:])
	r.Rows[row] = newrow
	return row
}

func (r *Record) Row(row int) *Row {
	if row < 0 || row >= len(r.Rows) {
		panic("row exceed")
	}

	return r.Rows[row]
}

func (r *Record) Limit() {
	if len(r.Rows) > r.MaxRow {
		r.Rows = r.Rows[:r.MaxRow]
	}
}

func (r *Record) Len() int {
	return len(r.Rows)
}

func (r *Record) Swap(i, j int) {
	r.Rows[i], r.Rows[j] = r.Rows[j], r.Rows[i]
}

func (r *Record) Less(i, j int) bool {
	return r.Rows[i].Value(r.sort_key).Less(r.Rows[j].Value(r.sort_key))
}

func (r *Record) Sort(key int, mode string) bool {
	r.sort_key = key
	if r.sort_key < 0 || r.sort_key >= r.cols {
		return false
	}

	if mode == "des" {
		sort.Sort(sort.Reverse(r))
	} else {
		sort.Sort(r)
	}

	return true
}

type RecordList struct {
	records []*Record
	ki      map[string]int
}

func NewRecordList() *RecordList {
	r := &RecordList{}
	r.records = make([]*Record, 0, 16)
	r.ki = make(map[string]int)
	return r
}

func (r *RecordList) Size() int {
	return len(r.records)
}

func (r *RecordList) Count() int {
	return len(r.ki)
}

func (r *RecordList) Clear() {
	r.records = r.records[:0]
	r.ki = make(map[string]int)
}

func (r *RecordList) NameList() []string {
	ret := make([]string, 0, len(r.ki))
	size := len(r.records)
	for i := 0; i < size; i++ {
		if r.records[i] != nil {
			ret = append(ret, r.records[i].RecName)
		}
	}
	return ret
}

func (r *RecordList) AddRecord(name string, rec *Record) {
	index := len(r.records)
	r.records = append(r.records, rec)
	r.ki[name] = index
}

func (r *RecordList) RemoveRecord(name string) {
	if index, has := r.ki[name]; has {
		delete(r.ki, name)
		r.records[index] = nil
	}
}

func (r *RecordList) Record(name string) *Record {
	if index, has := r.ki[name]; has {
		return r.records[index]
	}
	return nil
}

func (r *RecordList) RecordByIndex(index int) *Record {
	if index < 0 || index >= len(r.records) {
		return nil
	}

	return r.records[index]
}
