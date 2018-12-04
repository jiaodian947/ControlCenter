package gameobj

import (
	"manage/protocol"
	"testing"
)

func TestRecord_Sort(t *testing.T) {
	r1 := NewRecord(1, 5)
	r1.Column().SetType(0, protocol.VTYPE_INT)

	for i := 0; i < 5; i++ {
		r := r1.AddRow(-1)
		row := r1.Row(r)
		value := NewAny()
		value.SetInt(int32(i*2 + 1))
		row.SetValue(0, value)
	}

	for i := 0; i < 5; i++ {
		r := r1.AddRow(-1)
		row := r1.Row(r)
		value := NewAny()
		value.SetInt(int32(i * 2))
		row.SetValue(0, value)
	}

	for i := 0; i < r1.RowCount(); i++ {
		t.Log(i, r1.Row(i).Value(0).Value())
	}

	r1.Sort(0, "desc")

	for i := 0; i < r1.RowCount(); i++ {
		t.Log(i, r1.Row(i).Value(0).Value())
	}
}
