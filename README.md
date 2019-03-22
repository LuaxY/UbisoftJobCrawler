Ubisoft Job Crawler
===================

Crawl SmartRecruiters website for new Ubisoft jobs on target city.  
Send new offers to IFTTT WebHook

```
go build -o ubisoft main.go
./ubisoft -city {CITY} -ifttt-event {EVENT} -ifttt-key {KEY}
```

You can create cron task to run crawler every day at 10 AM

```
0 10 * * * /path/to/ubisoft -city Montreal -ifttt-event UbisoftJob -ifttt-key {KEY} >/dev/null 2>&1 
```