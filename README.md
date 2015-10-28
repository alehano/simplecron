
Simple Cron. Every job runs in goroutine.


 Patterns:
 12 - at 12
 1,2,3 - at 1 or 2 or 3
 * - every hour/min
 */15 - every 15 hours/min



WEEKDAYS:
Sunday = 0
Monday = 1
Tuesday
Wednesday
Thursday
Friday
Saturday


MONTHS
January = 0
February = 1
March
April
May
June
July
August
September
October
November
December


Example:

    cron.CronInst.AddJob("job1", "job one descr", cron.NewJob(&CronRunner{}, "17", "3,6", "*", "0,3"))

Run CronRun() func on CronRunner every Sunday and  Wednesday at 3:17 and 6:17 am.


WARNING:
To run on multiple machines, use cron as a microservice.
