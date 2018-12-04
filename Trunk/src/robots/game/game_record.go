package game

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

type Record struct {
	RecName string
	Columns *Column
	Rows    []*Row
	cols    int
	MaxRow  int
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
	return r
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

func (r *Record) DelRow(row int) {
	if row < 0 || row >= len(r.Rows) {
		return
	}

	copy(r.Rows[row:], r.Rows[row+1:])
	r.Rows = r.Rows[:len(r.Rows)-1]
}

func (r *Record) Row(row int) *Row {
	if row < 0 || row >= len(r.Rows) {
		panic("row exceed")
	}

	return r.Rows[row]
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
