package service

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"
)

type workers struct {
	*sync.RWMutex
	w      []*worker
	status map[int]wStatus
}

type worker struct {
	id    int
	jobs  chan job
	stats chan wStats
	p     *params
}
type wStatus struct {
	src, dst string
}
type wStats struct {
	workerId int
	src      string
	dst      string
	size     int64
	time     time.Duration
	status   string
	err      error
	msgErr   string
}

func (w *worker) run(ctx context.Context, wg *sync.WaitGroup, pool *workers) {
	w.stats <- wStats{workerId: w.id, status: "launched"}
	rand.Seed(time.Now().UnixNano())
	dst := w.p.Dst
	cnt := len(dst) - 1
	for {
		select {
		case job := <-w.jobs:
			info, err := os.Stat(job.src)
			if err != nil {
				w.stats <- wStats{workerId: w.id, status: "error", err: err, msgErr: "", src: job.src}
				continue
			}
			var target string
			if cnt == 0 {
				target = dst[0]
			} else {
				target = dst[rand.Intn(cnt)]
			}

			// TODO анализ свободного пространства и учёт текущих событий записи
			/*
				if b, err := exec.Command("df", target).Output(); err != nil {
					fmt.Println("Failed to initiate command:", err)
				} else {
					fmt.Println(info.ModTime(), info.Size(), "mv", target, "\n", string(b))
					s := strings.Split(string(b), " ")

					if len(s) > 5 {
						if size, err := strconv.Atoi(s[len(s)-4]); err == nil {
							size64 := int64(size) * 1048576

						} else {
							w.stats <- wStats{workerId: w.id, status: "error", src: job.src, err: err, msgErr: string(b)}
						}
					}
					fmt.Println()
				}
			*/
			s := wStats{}
			s.workerId = w.id
			s.src = job.src
			s.dst = target
			s.size = info.Size()
			s.status = "startMV"
			s.time = 0
			fmt.Println(s.String())
			w.stats <- s

			startJob := time.Now()
			if b, err := exec.Command("mv", s.src, s.dst).Output(); err != nil {
				fmt.Println("Failed to initiate command:", err)
				w.stats <- wStats{workerId: w.id, status: "error", err: err, msgErr: string(b), src: s.src, dst: s.dst}
				continue
			}
			s.time = time.Since(startJob)
			s.status = "doneMV"
			fmt.Println(s.String())
			w.stats <- s

		case <-ctx.Done():
			w.stats <- wStats{workerId: w.id, status: "stopped"}
			wg.Done()
			return
		}
	}
}

func (s *wStats) String() string {
	switch s.status {
	case "launched", "stopped":
		return fmt.Sprintf("w%d %v", s.workerId, s.status)
	case "error":
		return fmt.Sprintf("w%d %v %v %v %v %v", s.workerId, s.status, s.err.Error(), s.msgErr)
	}
	return fmt.Sprintf("w%d %v %v %v from: %s to: %s", s.workerId, s.time, s.size, s.status, s.src, s.dst)
}
