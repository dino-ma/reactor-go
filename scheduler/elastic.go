package scheduler

import (
	"math"

	"github.com/panjf2000/ants"
)

var elastic Scheduler

func init() {
	elastic = NewElastic(math.MaxInt32)
}

type elasticScheduler struct {
	pool *ants.Pool
}

func (p *elasticScheduler) Close() error {
	return p.pool.Release()
}

func (p *elasticScheduler) Do(job Job) {
	if err := p.pool.Submit(job); err != nil {
		panic(err)
	}
}

func (p *elasticScheduler) Worker() Worker {
	return p
}

func NewElastic(size int) Scheduler {
	pool, err := ants.NewPool(size)
	if err != nil {
		panic(err)
	}
	return &elasticScheduler{
		pool: pool,
	}
}

func Elastic() Scheduler {
	return elastic
}
