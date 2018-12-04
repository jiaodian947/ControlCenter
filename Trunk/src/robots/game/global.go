package game

import (
	"sync"
)

type PropInfo struct {
	Name string
	Type int8
}

type PropTable struct {
	Props []PropInfo
	KI    map[string]int
}

type RecInfo struct {
	Name    string
	Cols    uint16
	ColType []uint8
}

type RecTable struct {
	Recs []RecInfo
	KI   map[string]int
}

var (
	PropTables *PropTable
	RecTables  *RecTable
	Mtx        sync.RWMutex
)

func CreatePropTables() bool {
	Mtx.Lock()
	if PropTables != nil {
		return false
	}
	PropTables = &PropTable{}
	Mtx.Unlock()
	return true
}
