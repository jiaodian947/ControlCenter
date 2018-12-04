package gameobj

type ObjID struct {
	Ident  uint32
	Serial uint32
}

func NewObjId() ObjID {
	id := ObjID{}
	return id
}
