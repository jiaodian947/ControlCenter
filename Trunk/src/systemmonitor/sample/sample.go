package sample

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"
)

type ProcSample struct {
	Time time.Time
	Pid  int32
	Usr  float64
	Rss  uint64
}

type GameSample struct {
	Time      time.Time
	TotalUser int32
	Online    int32
	MaxOnline int32
	MaxTime   time.Time
	Average   int32
	Scenes    []*SceneInfo
}

type SceneInfo struct {
	SceneId   int32
	SceneName string
	Players   int32
}

type Sample struct {
	cache      chan *ProcSample
	saveToFile chan *ProcSample
	gameSample *GameSample
	gameInfo   *GameInfo
	ctx        context.Context
	cancel     context.CancelFunc
	addr       string
	port       int
	files      map[int]io.WriteCloser
}

func NewSample(maxcache int, addr string, port int) *Sample {
	s := &Sample{}
	s.cache = make(chan *ProcSample, maxcache)
	s.saveToFile = make(chan *ProcSample, maxcache)
	s.gameSample = &GameSample{Scenes: make([]*SceneInfo, 0, 1024)}
	s.ctx = context.Background()
	s.gameInfo = NewGameInfo(s)
	s.addr = addr
	s.port = port
	return s
}

var splitter = regexp.MustCompile("[, ] *")

// ParsePidList take a string of process ids and converts it into a list of int pids
func ParsePidList(s string) ([]int, error) {
	parts := splitter.Split(s, -1)
	ret := make([]int, len(parts))
	for pos, part := range parts {
		if len(part) == 0 {
			continue
		}
		num, err := strconv.Atoi(part)
		if err != nil {
			panic(err)
		}
		ret[pos] = num
	}
	return ret, nil
}

func (s *Sample) StoreSample(ps *ProcSample) {
	s.saveToFile <- ps
	for {
		select {
		case s.cache <- ps:
			return
		default:
			<-s.cache
		}
	}
}

func (s *Sample) GameSample() *GameSample {
	return s.gameSample
}

func (s *Sample) BeginGameInfo() {
	s.gameSample.Time = time.Now()
}

func (s *Sample) StoreOnlinePlayer(players int32) {
	s.gameSample.Online = players
}

func (s *Sample) StoreMaxPlayer(players int32) {
	s.gameSample.TotalUser = players
}

func (s *Sample) StoreSceneInfo(scene []*SceneInfo) {
	s.gameSample.Scenes = s.gameSample.Scenes[:0]
	s.gameSample.Scenes = append(s.gameSample.Scenes, scene...)
}

func (s *Sample) Query() <-chan *ProcSample {
	return s.cache
}

func (s *Sample) Sample(pids string) {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}
	s.gameInfo.SetAddrPort(s.addr, s.port)
	ctx, cancel := context.WithCancel(s.ctx)
	s.cancel = cancel
	go s.BeginSample(ctx, 200, 5, 100, 2048, "", pids, 100)
	go s.gameInfo.Run(ctx)
	go s.WriteSample(ctx, pids)
}

func (s *Sample) WriteSample(ctx context.Context, pid string) {
	pids, err := ParsePidList(pid)
	if err != nil {
		return
	}

	s.files = make(map[int]io.WriteCloser)
	t := time.Now()
	for _, v := range pids {
		f, err := os.Create(fmt.Sprintf("%d-%d%d%d-%d%d%d.sample", v, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
		if err != nil {
			panic(err)
		}

		s.files[v] = f
	}
L:
	for {
		select {
		case <-ctx.Done():
			break L
		case ps := <-s.saveToFile:
			s.WriteToFile(ps)
		}
	}

	for _, f := range s.files {
		f.Close()
	}

	s.files = make(map[int]io.WriteCloser)
}

func (s *Sample) WriteToFile(ps *ProcSample) {
	if f, ok := s.files[int(ps.Pid)]; ok {
		line := fmt.Sprintf("%d,%.2f,%d\n", ps.Time.Unix(), ps.Usr, ps.Rss)
		f.Write([]byte(line))
	}
}

func (s *Sample) StopSample() {
	if s.cancel != nil {
		s.cancel()
		s.cancel = nil
	}

	for {
		select {
		case <-s.cache:
		default:
			return
		}
	}
}
