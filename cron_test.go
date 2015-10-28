package cron

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	//	"time"
	"strconv"
)

type MyType struct {
	val string
}

var isErr = true

func (mt *MyType) CronRun() {
	isErr = false
	fmt.Printf("COMMAND RAN: %s", mt.val)
}

func TestCron(t *testing.T) {
	cron := NewCron()
	j := NewJob(&MyType{"DONE"}, "test job", "",
		strconv.Itoa(int(time.Now().Minute())),
		strconv.Itoa(int(time.Now().Hour())),
		strconv.Itoa(int(time.Now().Month())),
		strconv.Itoa(int(time.Now().Weekday())))
	cron.AddJob(j)
	t.Log("Start cron")
	t.Log("M:" + j.minute)
	t.Log("H:" + j.hour)
	t.Log("MON:" + j.month)
	t.Log("WD:" + j.weekday)
	cron.Start()
	time.Sleep(time.Second * 1)
	if isErr {
		t.Error("COMMAND NOT RUN")
	}
	t.Log("Stop")
}

func TestCheckValue(t *testing.T) {

	var cases = []struct {
		pattern string
		value   int
		res     bool
	}{
		{"1", 1, true},
		{"5", 5, true},
		{"60", 61, false},
		{"*/1", 0, true},
		{"*/1", 1, true},
		{"*/1", 2, true},
		{"*/1", 3, true},
		{"*/1", 4, true},
		{"*/15", 0, true},
		{"*/15", 15, true},
		{"*/15", 30, true},
		{"*/15", 45, true},
		{"*/15", 60, true},
		{"*/15", 10, false},
		{"*/15", 14, false},
		{"*/15", 149, false},
		{"*/2", 0, true},
		{"*/2", 2, true},
		{"*/2", 4, true},
		{"*/2", 128, true},
		{"*/2", 1, false},
		{"*/2", 3, false},
		{"*/2", 5, false},
		{"*/2", 7, false},
		{"*/2", 127, false},
		{"*/2", 129, false},
		{"*/11", 0, true},
		{"*/11", 11, true},
		{"*/11", 22, true},
		{"*/11", 33, true},
		{"*/11", 12, false},
		{"*/11", 14, false},
		{"*/11", 133, false},
		{"1,2,3", 0, false},
		{"1,2,3", 123, false},
		{"1,2,3", 4, false},
		{"1,2,3", 1, true},
		{"1,2,3", 2, true},
		{"1,2,3", 3, true},
		{"12,60", 12, true},
		{"12,60", 60, true},
		{"12,60", 120, false},
		{"12,60", 600, false},
		{"12,60", 0, false},
		{"*", 0, true},
		{"*", 123, true},
		{"*", 1, true},
		{"*", 12213123, true},
	}

	for _, c := range cases {
		assert.Equal(t, c.res, isNow(c.pattern, c.value), "wrong: "+c.pattern)
	}

}

func TestPattern(t *testing.T) {

	var cases = []struct {
		pattern string
		res     bool
	}{
		{"*", true},
		{"0", true},
		{"1", true},
		{"5", true},
		{"50", true},
		{"60", true},
		{"a5", false},
		{"1,2,3,4,5", true},
		{"1,2,3", true},
		{"1,2,3", true},
		{"a,1,2,3", false},
		{"1,2,3a", false},
		{"*/12", true},
		{"*/123", false},
		{"*/", false},
		{"*/abc", false},
	}

	for _, c := range cases {
		assert.Equal(t, c.res, isPatternValid(c.pattern), "wrong: "+c.pattern)
	}

}
