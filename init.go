package cron

// Cron instance
var CronInst Cron

func init() {
	CronInst = NewCron()
	CronInst.Start()
}
