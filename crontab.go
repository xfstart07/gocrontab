package gocrontab

import (
	"fmt"
	"sort"
	"time"
)

type Scheduler struct {
	jobs []*Job

	size int

	running bool

	stop chan bool
}

type Job struct {

	// 任务名
	jobName string

	// 存储函数Map
	jobFunc FuncJob

	interval uint64 // 设定为秒

	period time.Duration // 周期

	lastTime time.Time // 最后运行时间

	nextTime time.Time // 下次运行时间
}

type FuncJob func()

func (f FuncJob) Run() { f() }

func NewSchedule() *Scheduler {
	return &Scheduler{
		jobs:    nil,
		running: false,
		stop:    make(chan bool, 1),
	}
}

// sort

func (s *Scheduler) Len() int {
	return len(s.jobs)
}

func (s *Scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

// 判断 i 是否在 j 之前
func (s *Scheduler) Less(i, j int) bool {
	return s.jobs[i].nextTime.Before(s.jobs[j].nextTime)
}

// add job
func (s *Scheduler) AddJob(interval uint64, Name string, jobFunc func()) {
	job := &Job{
		interval: interval,
		period:   time.Duration(interval),
		lastTime: time.Now(),
		jobName:  Name,
		jobFunc:  FuncJob(jobFunc),
	}
	job.nextTime = job.lastTime.Add(job.period * time.Second)

	s.jobs = append(s.jobs, job)
}

// Entries ...
func (s *Scheduler) Entries() []*Job {
	return s.jobs
}

// Start ...
func (s *Scheduler) Start() {
	s.running = true
	go s.run()
}

// Stop ...
func (s *Scheduler) Stop() {
	if !s.running {
		return
	}
	s.running = false
	s.stop <- true
}

func (s *Scheduler) RemoveJob(name string) bool {
	pos := s.posJob(name)
	if pos == -1 {
		return false
	}

	size := len(s.jobs)
	for i := (pos + 1); i < size; i++ {
		s.jobs[i-1] = s.jobs[i]
	}
	s.jobs[size-1] = nil

	return true
}

func (s *Scheduler) run() {

	for {
		sort.Sort(s)

		for idx := range s.jobs {
			fmt.Println(s.jobs[idx].jobName)
		}

		var timer *time.Timer
		now := time.Now()
		if len(s.jobs) == 0 || s.jobs[0].nextTime.IsZero() {
			timer = time.NewTimer(10000 * time.Hour)
		} else {
			timer = time.NewTimer(s.jobs[0].nextTime.Sub(now))
		}

		for {
			// 使用 select 接收 channel 信号
			select {
			case now := <-timer.C:

				for idx := range s.jobs {
					fmt.Println(s.jobs[idx].jobName)
					if now.After(s.jobs[idx].nextTime) {
						s.jobs[idx].jobFunc.Run()
						s.jobs[idx].lastTime = now
						s.jobs[idx].nextTime = s.jobs[idx].lastTime.Add(s.jobs[idx].period * time.Second)
					} else {
						break
					}
				}

			case <-s.stop:
				timer.Stop()
				return
			}

			break
		}

	}
}

func (s *Scheduler) posJob(name string) int {
	for idx := range s.jobs {
		if s.jobs[idx].jobName == name {
			return idx
		}
	}

	return -1
}