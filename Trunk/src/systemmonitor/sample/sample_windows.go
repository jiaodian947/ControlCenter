// +build windows

package sample

import (
	"context"
	"log"
	"time"

	"github.com/shirou/gopsutil/process"
)

func (s *Sample) BeginSample(ctx context.Context, interval, samples, topN, maxProcsToScan int, usrOnly, pidOnly string, jiffy int) {
	pids, err := ParsePidList(pidOnly)
	if err != nil {
		log.Fatal(err)
	}

	ps := make([]*process.Process, len(pids))
	for k := range ps {
		ps[k], err = process.NewProcess(int32(pids[k]))
		if err != nil {
			log.Fatal(err)
		}
	}

	tick := time.Duration(interval*samples) * time.Millisecond
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for _, p := range ps {
			process := p
			go func() {
				c, err := process.CPUPercent()
				if err != nil {
					return
				}
				m, err := process.MemoryInfo()
				if err != nil {
					return
				}

				ps := &ProcSample{}
				ps.Time = time.Now()
				ps.Usr = c
				ps.Rss = m.RSS / 4096
				ps.Pid = process.Pid
				s.StoreSample(ps)
			}()
		}
		time.Sleep(tick)
	}
}
