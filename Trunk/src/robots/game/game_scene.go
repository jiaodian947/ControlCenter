package game

type GameScene struct {
	GameObject
}

func NewGameScene() *GameScene {
	s := &GameScene{}
	s.Init()
	return s
}
