package utils

import "sync"

type limitGroup struct {
	mutex *sync.Mutex

	wg    sync.WaitGroup
	c     *sync.Cond
	limit int
}

func NewLimitGroup(limit int) *limitGroup {
	// To avoid lock state, if limit is to small set it to big number.
	if limit <= 0 {
		limit = 100
	}

	mutex := new(sync.Mutex)

	return &limitGroup{
		mutex: mutex,

		wg:    sync.WaitGroup{},
		c:     sync.NewCond(mutex),
		limit: limit,
	}
}

func (lg *limitGroup) Add() {
	lg.mutex.Lock()
	defer lg.mutex.Unlock()

	for lg.limit < 1 {
		lg.c.Wait()
	}

	lg.limit -= 1
	lg.wg.Add(1)
}

func (lg *limitGroup) Done() {
	lg.mutex.Lock()
	defer lg.mutex.Unlock()

	lg.limit++
	lg.c.Signal()
	lg.wg.Done()
}

func (lg *limitGroup) Wait() { lg.wg.Wait() }
