package cron

import (
	"regexp"
	"sync"
	"time"

	"fmt"
	"strconv"
	"strings"
)

func NewCron() Cron {
	return Cron{jobs: make(map[string]*Job)}
}

// Patterns:
// 12 - at 12
// 1,2,3 - at 1 or 2 or 3
// * - every hour/min
// */15 - every 15 hours/min
func NewJob(r CronRunner, name, descr, minute, hour, month, weekday string) *Job {
	return &Job{runner: r, name: name, descr: descr, minute: minute, hour: hour, month: month, weekday: weekday}
}

type CronRunner interface {
	CronRun()
}

type Job struct {
	runner  CronRunner
	name    string
	descr   string
	minute  string
	hour    string
	month   string
	weekday string
}

type Cron struct {
	jobs  map[string]*Job
	mutex sync.Mutex
}

func (c Cron) String() string {
	res := fmt.Sprintf("\tCron jobs:\n")
	res += fmt.Sprintf("\t---------:\n")
	for name, job := range c.jobs {
		res += fmt.Sprintf("\t\t%s: (m:%s, h:%s, month:%s, weekday:%s)  %s\n", name, job.minute, job.hour, job.month, job.weekday, job.descr)
	}
	res += fmt.Sprintf("\tNow time: %s", time.Now())
	return res
}

func (c *Cron) AddJob(j *Job) {
	c.mutex.Lock()
	if _, ok := c.jobs[j.name]; ok {
		panic("Cron job already exists: " + j.name)
	}
	if !isPatternValid(j.minute) ||
		!isPatternValid(j.hour) ||
		!isPatternValid(j.month) ||
		!isPatternValid(j.weekday) {
		panic(fmt.Sprintf("Cron: %s pattern not valid", j.name))
	}
	c.jobs[j.name] = j
	c.mutex.Unlock()
}

func (c *Cron) RemoveJob(name string, j Job) {
	c.mutex.Lock()
	delete(c.jobs, name)
	c.mutex.Unlock()

}

func (c *Cron) RunByNameAsync(name string) {
	go func() {
		c.jobs[name].runner.CronRun()
	}()
}

func (c *Cron) RunByName(name string) {
	c.jobs[name].runner.CronRun()
}

func (c *Cron) Start() {
	go func() {
		for {
			for jobName, job := range c.jobs {
				if isNow(job.weekday, int(time.Now().Weekday())) &&
					isNow(job.month, int(time.Now().Month())) &&
					isNow(job.hour, int(time.Now().Hour())) &&
					isNow(job.minute, int(time.Now().Minute())) {
					c.RunByNameAsync(jobName)
				}
			}
			time.Sleep(time.Minute)
		}
	}()
}

// Compares nowVal with pattern
func isNow(pattern string, nowVal int) bool {
	if pattern == "*" {
		return true
	}
	if strings.HasPrefix(pattern, "*/") {
		val, err := strconv.Atoi(strings.TrimPrefix(pattern, "*/"))
		if err != nil {
			return false
		}
		if nowVal%val == 0 {
			return true
		} else {
			return false
		}
	}
	vals := strings.Split(pattern, ",")
	for _, val := range vals {
		valInt, err := strconv.Atoi(val)
		if err != nil {
			return false
		}
		if valInt == nowVal {
			return true
		}
	}
	return false
}

// Note: don't check length of value
func isPatternValid(p string) bool {
	matched, err := regexp.MatchString(`^\*$|^(\*\/?[0-9]{1,2})$|^([0-9],?)+$`, p)
	if err != nil {
		return false
	}
	return matched
}
