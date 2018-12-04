// +build linux

package sample

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	lib "github.com/uber-common/cpustat/lib"
)

func checkPrivs() {
	if os.Geteuid() != 0 {
		fmt.Println("This program uses the netlink taskstats inteface, so it must be run as root.")
		os.Exit(1)
	}
}

func formatMem(num uint64) string {
	letter := string("K")

	num = num * 4
	if num >= 1000 {
		num = (num + 512) / 1024
		letter = "M"
		if num >= 10000 {
			num = (num + 512) / 1024
			letter = "G"
		}
	}
	return fmt.Sprintf("%d%s", num, letter)
}

func formatNum(num uint64) string {
	if num > 1000000 {
		return fmt.Sprintf("%dM", num/1000000)
	}
	if num > 1000 {
		return fmt.Sprintf("%dK", num/1000)
	}
	return fmt.Sprintf("%d", num)
}

func trim(num float64, max int) string {
	var str string
	if num >= 1000.0 {
		str = fmt.Sprintf("%d", int(num+0.5))
	} else {
		str = fmt.Sprintf("%.1f", num)
	}
	if len(str) > max {
		if str[max-1] == 46 { // ASCII .
			return str[:max-1]
		}
		return str[:max]
	}
	if str == "0.0" {
		return "0"
	}
	return str
}

func trunc(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length]
}

func textInit(interval, samples, topN int, filters lib.Filters) {
	fmt.Printf("sampling interval:%s, summary interval:%s (%d samples), showing top %d procs,",
		time.Duration(interval)*time.Millisecond,
		time.Duration(interval*samples)*time.Millisecond,
		samples, topN)
	fmt.Print(" user filter:")
	if len(filters.User) == 0 {
		fmt.Print("all")
	} else {
		fmt.Print(strings.Join(filters.UserStr, ","))
	}
	fmt.Print(", pid filter:")
	if len(filters.Pid) == 0 {
		fmt.Print("all")
	} else {
		fmt.Print(strings.Join(filters.PidStr, ","))
	}
	fmt.Println()
}

func (s *Sample) dumpStats(infoMap lib.ProcInfoMap, list lib.Pidlist, procSum lib.ProcSampleMap,
	procHist lib.ProcStatsHistMap, taskHist lib.TaskStatsHistMap,
	sysSum *lib.SystemStats, sysHist *lib.SystemStatsHist, jiffy, interval, samples int) {

	scaleSum := func(val float64, count int64) float64 {
		valSec := val / float64(jiffy)
		sampleSec := float64(interval) * float64(count) / 1000.0
		ret := (valSec / sampleSec) * 100
		return ret
	}

	for _, pid := range list {
		sampleCount := procHist[pid].Ustime.TotalCount()
		if proc, ok := procSum[pid]; ok {
			ps := &ProcSample{}
			ps.Time = time.Now()
			ps.Usr = scaleSum(float64(proc.Proc.Utime), sampleCount)
			ps.Rss = proc.Proc.Rss
			ps.Pid = int32(pid)
			s.StoreSample(ps)
		}

	}
}

func (s *Sample) BeginSample(ctx context.Context, interval, samples, topN, maxProcsToScan int, usrOnly, pidOnly string, jiffy int) {

	checkPrivs()

	if interval < 10 {
		fmt.Println("The minimum sampling interval is 10ms")
		os.Exit(1)
	}
	intervalms := uint32(interval)
	filters := lib.FiltersInit(usrOnly, pidOnly)
	nlConn := lib.NLInit()

	textInit(interval, samples, topN, filters)

	infoMap := make(lib.ProcInfoMap)

	procCur := lib.NewProcSampleList(maxProcsToScan)
	procPrev := lib.NewProcSampleList(maxProcsToScan)
	procSum := make(lib.ProcSampleMap)
	procHist := make(lib.ProcStatsHistMap)
	taskHist := make(lib.TaskStatsHistMap)

	var sysCur lib.SystemStats
	var sysPrev lib.SystemStats
	var sysSum *lib.SystemStats
	var sysHist *lib.SystemStatsHist

	var t1, t2 time.Time
	var err error

	// run all scans one time to establish a baseline
	pids := make(lib.Pidlist, 0, maxProcsToScan)

	t1 = time.Now()
	lib.GetPidList(&pids, maxProcsToScan)
	lib.ProcStatsReader(pids, filters, &procPrev, infoMap)
	lib.TaskStatsReader(nlConn, pids, &procPrev)
	err = lib.SystemStatsReader(&sysPrev)
	if err != nil {
		panic(err)
	}

	sysSum = &lib.SystemStats{}
	sysHist = lib.NewSysStatsHist()
	t2 = time.Now()

	targetSleep := time.Duration(interval) * time.Millisecond
	adjustedSleep := targetSleep - t2.Sub(t1)

	topPids := make(lib.Pidlist, topN)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		for count := 0; count < samples; count++ {
			time.Sleep(adjustedSleep)

			t1 = time.Now()
			lib.GetPidList(&pids, maxProcsToScan)

			lib.ProcStatsReader(pids, filters, &procCur, infoMap)
			lib.TaskStatsReader(nlConn, pids, &procCur)

			procDelta := make(lib.ProcSampleMap, len(pids))
			lib.ProcStatsRecord(intervalms, procCur, procPrev, procSum, procDelta)
			lib.UpdateProcStatsHist(procHist, procDelta)
			lib.TaskStatsRecord(intervalms, procCur, procPrev, procSum, procDelta)
			lib.UpdateTaskStatsHist(taskHist, procDelta)

			procPrev, procCur = procCur, procPrev

			if err = lib.SystemStatsReader(&sysCur); err != nil {
				log.Fatal(err)
			}
			sysDelta := lib.SystemStatsRecord(intervalms, &sysCur, &sysPrev, sysSum)
			lib.UpdateSysStatsHist(sysHist, sysDelta)
			sysPrev = sysCur

			t2 = time.Now()
			adjustedSleep = targetSleep - t2.Sub(t1)
		}

		topHist := sortList(procHist, taskHist, topN)
		topPids = topPids[:len(topHist)]
		for i := 0; i < len(topHist) && i < topN; i++ {
			topPids[i] = topHist[i].pid
		}

		s.dumpStats(infoMap, topPids, procSum, procHist, taskHist, sysSum, sysHist, jiffy, interval, samples)
		procHist = make(lib.ProcStatsHistMap)
		taskHist = make(lib.TaskStatsHistMap)
		procSum = make(lib.ProcSampleMap)
		sysHist = lib.NewSysStatsHist()
		sysSum = &lib.SystemStats{}
		t2 = time.Now()
		adjustedSleep = targetSleep - t2.Sub(t1)
		// If we can't keep up, try to buy ourselves a little headroom by sleeping for a magic number of ms
		if adjustedSleep <= 0 {
			adjustedSleep = time.Duration(100) * time.Millisecond
		}
	}
}

// Wrapper to sort histograms by max but remember which pid they are
type sortHist struct {
	pid  int
	proc *lib.ProcStatsHist
	task *lib.TaskStatsHist
}

// ByMax sorts stats by max usage
type ByMax []*sortHist

func (m ByMax) Len() int {
	return len(m)
}
func (m ByMax) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
func (m ByMax) Less(i, j int) bool {
	var maxI, maxJ float64

	// We might have proc stats but no taskstats because of unfortuante timing
	if m[i].task == nil || m[j].task == nil {
		maxI = maxList([]float64{
			float64(m[i].proc.Ustime.Max()),
			float64(m[i].proc.Cutime.Max()+m[i].proc.Cstime.Max()) / 1000,
		})
		maxJ = maxList([]float64{
			float64(m[j].proc.Ustime.Max()),
			float64(m[j].proc.Cutime.Max()+m[j].proc.Cstime.Max()) / 1000,
		})
	} else {
		maxI = maxList([]float64{
			float64(m[i].proc.Ustime.Max()),
			float64(m[i].proc.Cutime.Max()+m[i].proc.Cstime.Max()) / 100,
			float64(m[i].task.Cpudelay.Max()) / 1000 / 1000,
			float64(m[i].task.Iowait.Max()) / 1000 / 1000,
			float64(m[i].task.Swap.Max()) / 1000 / 1000,
		})
		maxJ = maxList([]float64{
			float64(m[j].proc.Ustime.Max()),
			float64(m[j].proc.Cutime.Max()+m[j].proc.Cstime.Max()) / 100,
			float64(m[j].task.Cpudelay.Max()) / 1000 / 1000,
			float64(m[j].task.Iowait.Max()) / 1000 / 1000,
			float64(m[j].task.Swap.Max()) / 1000 / 1000,
		})
	}
	return maxI > maxJ
}
func maxList(list []float64) float64 {
	ret := list[0]
	for i := 1; i < len(list); i++ {
		if list[i] > ret {
			ret = list[i]
		}
	}
	return ret
}

func sortList(procHist lib.ProcStatsHistMap, taskHist lib.TaskStatsHistMap, limit int) []*sortHist {
	var list []*sortHist

	// let's hope that pid is in both sets, otherwise this will blow up
	for pid, hist := range procHist {
		list = append(list, &sortHist{pid, hist, taskHist[pid]})
	}
	sort.Sort(ByMax(list))

	if len(list) > limit {
		list = list[:limit]
	}

	return list
}
