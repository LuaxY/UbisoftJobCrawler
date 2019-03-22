Ubisoft Job Crawler
===================

Crawl SmartRecruiters website for new Ubisoft jobs on target city.  
Send new offers to IFTTT WebHook

```
go build -o ubisoft main.go
./ubisoft -city {CITY} -ifttt-event {EVENT} -ifttt-key {KEY}
```