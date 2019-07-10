package gocrontab

import (
	"reflect"
	"sort"
	"time"

	"github.com/pkg/errors"
)

var loc = time.Local

type Scheduler struct {
	jobs []*Job

	running bool

	stop chan struct{}
}

type Job struct {

	// 任务名
	jobName string

	// 存储函数Map
	jobFunc interface{}

	jobParams []interface{}

	interval uint64 // 设定为秒

	unit string // 单位

	period time.Duration // 周期

	lastTime time.Time // 最后运行时间

	nextTime time.Time // 下次运行时间
}

// Run 运行 Job
// 使用 reflect 进行调用
func (j Job) Run() {
	fn := reflect.ValueOf(j.jobFunc)

	fparam := make([]reflect.Value, len(j.jobParams))
	for key, value := range j.jobParams {
		fparam[key] = reflect.ValueOf(value)
	}

	fn.Call(fparam)
}

// NewSchedule 新建
func NewSchedule() *Scheduler {
	return &Scheduler{
		jobs:    nil,
		running: false,
		stop:    make(chan struct{}, 1),
	}
}

// sort

// Len ...
func (s *Scheduler) Len() int {
	return len(s.jobs)
}

// Swap ...
func (s *Scheduler) Swap(i, j int) {
	s.jobs[i], s.jobs[j] = s.jobs[j], s.jobs[i]
}

// Less 判断 i 是否在 j 之前
func (s *Scheduler) Less(i, j int) bool {
	return s.jobs[i].nextTime.Before(s.jobs[j].nextTime)
}

// NewJob
func (s *Scheduler) NewJob(Name string) *Job {
	job := &Job{
		interval: 0,
		period:   0,
		jobName:  Name,
	}
	s.jobs = append(s.jobs, job)

	return job
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
	s.stop <- struct{}{}
}

// RemoveJob ...
func (s *Scheduler) RemoveJob(name string) bool {
	pos := s.posJob(name)
	if pos < 0 {
		return false
	}

	s.jobs = s.jobs[:pos+copy(s.jobs[pos:], s.jobs[pos+1:])]

	return true
}

func (s *Scheduler) run() {
	for {
		sort.Sort(s)

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
					if now.After(s.jobs[idx].nextTime) {
						go s.jobs[idx].Run()

						s.jobs[idx].lastTime = time.Now()
						s.jobs[idx].shouldNextTime()
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

// job

// Do ...
func (j *Job) Do(jobFunc interface{}, params ...interface{}) {
	j.jobFunc = jobFunc
	j.jobParams = params

	j.shouldNextTime()
}

func (j *Job) shouldNextTime() {
	if j.lastTime.IsZero() {
		j.lastTime = time.Now()
	}

	if j.period == 0 {
		switch j.unit {
		case "seconds":
			j.period = time.Duration(j.interval)
		case "minutes":
			j.period = time.Duration(j.interval * 60)
		case "hours":
			j.period = time.Duration(j.interval * 60 * 60)
		case "days":
			j.period = time.Duration(j.interval * 60 * 60 * 24)
		}
	}

	j.nextTime = j.lastTime.Add(j.period * time.Second)
}

// Every
func (j *Job) Every(interval uint64) *Job {
	j.interval = interval
	return j
}

// Seconds unit
func (j *Job) Seconds() *Job {
	j.unit = "seconds"
	return j
}

// Minutes unit
func (j *Job) Minutes() *Job {
	j.unit = "minutes"
	return j
}

// Hours unit
func (j *Job) Hours() *Job {
	j.unit = "minutes"
	return j
}

// Days unit
func (j *Job) Days() *Job {
	j.unit = "days"
	return j
}

// At unit 时间：小时分钟，18, 20
func (j *Job) At(hour, min uint) *Job {
	if (hour < 0 || hour > 23) || (min < 0 || min > 59) {
		panic(errors.New("wrong time range"))
	}

	at := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), int(hour), int(min), 0, 0, loc)
	if j.unit == "days" && j.interval == 1 {
		if time.Now().After(at) {
			j.lastTime = at
		} else {
			dayDuration, _ := time.ParseDuration("-24h")
			j.lastTime = at.Add(dayDuration)
		}
	}

	return j
}

// GetName ...
func (j *Job) GetName() string {
	return j.jobName
}

// Rename ...
func (j *Job) Rename(name string) {
	j.jobName = name
}
