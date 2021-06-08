package service

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	log "github.com/n404an/gomv/logs"
)

type Overlord struct {
	wg    *sync.WaitGroup
	wPool *workers
	jobs  chan job
	stats chan wStats
	l     *log.Logger
	ctx   context.Context
	p     *params
}

type params struct {
	Src      []string
	Dst      []string
	FileMask string
	Parallel int
	NoRand   bool
}
type job struct {
	src  string
	size int64
}

func Start(ctx context.Context) error {
	p := &params{
		Src: make([]string, 0),
		Dst: make([]string, 0),
	}

	if err := p.getArgs(); err != nil {
		return err
	}
	if err := p.checkArgs(); err != nil {
		return err
	}

	overlord := newOverlord(p, ctx)
	defer overlord.l.Close()

	overlord.wg.Add(3)
	go overlord.collectStats()
	go overlord.runWorkers()
	go overlord.collectJobs()

	for {
		select {
		case <-ctx.Done():
			return nil
		}
	}
	overlord.wg.Wait()
	return nil
}

func (o *Overlord) collectStats() {
	defer o.wg.Done()
	for {
		select {
		case stat := <-o.stats:
			o.l.Info(stat.String())
		case <-o.ctx.Done():
			return
		}
	}
}

func (o *Overlord) runWorkers() {
	defer o.wg.Done()

	wg := &sync.WaitGroup{}

	for _, w := range o.wPool.w {
		wg.Add(1)
		go w.run(o.ctx, wg, o.wPool)
	}
	wg.Wait()
	<-o.ctx.Done()
}

func (o *Overlord) collectJobs() {
	defer o.wg.Done()

	allFiles := make(map[string]struct{}, 0)

	searchNewJobs := func() {
		fmt.Println("search new jobs...")
		files := make(map[string]struct{}, 0)
		for _, s := range o.p.Src {
			match, _ := filepath.Glob(filepath.Join(s, o.p.FileMask))
			for _, v := range match {
				if _, ok := allFiles[v]; !ok {
					fmt.Println("new job", v)
					o.l.Info("new job", v)
					allFiles[v] = struct{}{}
					files[v] = struct{}{}
				}
			}
		}
		for k, _ := range files {
			o.jobs <- job{src: k}
		}
	}

	searchNewJobs()

	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:
			searchNewJobs()
		case <-o.ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func newOverlord(p *params, ctx context.Context) *Overlord {
	jobs := make(chan job, 10000000)
	stats := make(chan wStats, 100)

	w := make([]*worker, 0)
	for i := 0; i < p.Parallel; i++ {
		w = append(w, &worker{
			id:    i,
			jobs:  jobs,
			stats: stats,
			p:     p.copy()})
	}
	return &Overlord{
		wg: &sync.WaitGroup{},
		wPool: &workers{
			RWMutex: &sync.RWMutex{},
			w:       w,
			status:  make(map[int]wStatus, 0),
		},
		jobs:  jobs,
		stats: stats,
		l:     log.NewLogger(),
		ctx:   ctx,
		p:     p,
	}
}

func (p *params) copy() *params {
	n := &params{}
	n.Src = make([]string, 0)
	n.Dst = make([]string, 0)

	n.Src = append(n.Src, p.Src...)
	n.Dst = append(n.Dst, p.Dst...)

	n.FileMask = p.FileMask
	n.NoRand = p.NoRand
	n.Parallel = p.Parallel
	return n
}
