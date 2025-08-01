package worker

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/tobiashort/ansi-go"
	"github.com/tobiashort/isatty-go"
)

type Pool interface {
	GetWorker() Worker
	Wait()
}

type pool struct {
	cap     int
	workers []*worker
	mutex   sync.Mutex
	wg      sync.WaitGroup
}

func NewPool(cap int) Pool {
	pool := &pool{}
	pool.cap = cap
	pool.workers = make([]*worker, cap)
	pool.mutex = sync.Mutex{}
	pool.wg = sync.WaitGroup{}
	if isatty.IsTerminal() {
		for range cap {
			fmt.Print(ansi.EraseEntireLine)
			fmt.Println()
		}
		fmt.Print(ansi.MoveCursorUp(pool.cap))
	}
	for i := range cap {
		pool.workers[i] = &worker{
			num:  i + 1,
			pool: pool,
			done: true,
			msg:  "",
		}
		pool.print(pool.workers[i])
	}
	return pool
}

func (p *pool) GetWorker() Worker {
	for {
		for _, worker := range p.workers {
			if worker.done {
				worker.done = false
				p.wg.Add(1)
				return worker
			}
		}
	}
}

func (p *pool) Wait() {
	p.wg.Wait()
	if isatty.IsTerminal() {
		fmt.Print(ansi.EraseFromCursorToEndOfScreen)
	}
}

func (p *pool) done() {
	p.wg.Done()
}

func (p *pool) print(w *worker) {
	if isatty.IsTerminal() {
		p.mutex.Lock()
		defer p.mutex.Unlock()
		capLen := len(strconv.Itoa(p.cap))
		prefixFmt := fmt.Sprintf("[Worker %%%dd] ", capLen)
		fmt.Print(ansi.MoveCursorDown(w.num))
		fmt.Print(ansi.EraseEntireLine)
		fmt.Printf(prefixFmt, w.num)
		fmt.Print(w.msg)
		fmt.Print(ansi.MoveCursorToColumn(0))
		fmt.Print(ansi.MoveCursorUp(w.num))
	}
}

func (p *pool) log(w *worker) {
	p.mutex.Lock()
	if isatty.IsTerminal() {
		fmt.Print(ansi.EraseEntireLine)
	}
	fmt.Printf("%s\n", w.msg)
	w.msg = ""
	if isatty.IsTerminal() {
		for range p.cap {
			fmt.Print(ansi.EraseEntireLine)
			fmt.Println()
		}
		fmt.Print(ansi.MoveCursorUp(p.cap))
	}
	p.mutex.Unlock()
	for _, w := range p.workers {
		p.print(w)
	}
}
