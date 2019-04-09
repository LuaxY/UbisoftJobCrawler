package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type Job struct {
	ID       int64
	Link     string
	Title    string
	City     string
	Contract string
}

const (
	LastIdFile = "last_job_id.txt"
	MaxPages   = 10
)

var (
	city       string
	IFTTTEvent string
	IFTTTKey   string
)

func main() {
	flag.StringVar(&city, "city", "Montreal", "Target city")
	flag.StringVar(&IFTTTEvent, "ifttt-event", "", "IFTTT WebHook Event")
	flag.StringVar(&IFTTTKey, "ifttt-key", "", "IFTTT WebHook Key")
	flag.Parse()

	os.Chdir(os.Getenv("HOME"))

	var allJobs []*Job

	for i := 0; i < MaxPages; i++ {
		jobs, err := GetJobs(city, i)

		if _, ok := err.(*NoJob); ok {
			break
		} else if err != nil {
			panic(err)
		}

		allJobs = append(allJobs, jobs...)
	}

	var previousLastId int64

	content, err := ioutil.ReadFile(LastIdFile)

	if err == nil {
		previousLastId, _ = strconv.ParseInt(string(content), 10, 64)
	}

	for _, job := range allJobs {
		if job.ID <= previousLastId {
			break
		}

		msg := fmt.Sprintf("%s - %s - %s - %s", job.Title, job.City, job.Contract, job.Link)
		http.Post("https://maker.ifttt.com/trigger/"+IFTTTEvent+"/with/key/"+IFTTTKey, "application/json", bytes.NewBuffer([]byte(`{"value1":"`+msg+`"}`)))
	}

	if len(allJobs) > 0 {
		ioutil.WriteFile(LastIdFile, []byte(strconv.FormatInt(allJobs[0].ID, 10)), os.ModePerm)
	}
}

func GetJobs(city string, page int) ([]*Job, error) {
	res, err := http.Get(fmt.Sprintf("https://careers.smartrecruiters.com/Ubisoft2/?search=&page=%d&location=%s", page, url.QueryEscape(city)))

	if err != nil {
		return nil, fmt.Errorf("unable to get url: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)

	if err != nil {
		return nil, fmt.Errorf("unable to parse html: %v", err)
	}

	var jobs []*Job

	doc.Find(".opening-job.job").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a").AttrOr("href", "")
		title := s.Find(".details-title").Text()
		city := s.Find(".details-desc .desc-item").Eq(0).Text()
		contract := s.Find(".details-desc .desc-item").Eq(1).Text()

		if title == "" {
			return
		}

		var id int64

		_, err := fmt.Sscanf(link, "https://jobs.smartrecruiters.com/Ubisoft2/%d-", &id)

		if err != nil {
			return
		}

		jobs = append(jobs, &Job{
			ID:       id,
			Link:     link,
			Title:    title,
			City:     city,
			Contract: contract,
		})
	})

	if len(jobs) <= 0 {
		return nil, &NoJob{}
	}

	return jobs, nil
}

type NoJob struct {
}

func (e *NoJob) Error() string {
	return ""
}
