package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/mileusna/crontab"
	"github.com/mmcdole/gofeed/rss"
	"github.com/namsral/flag"
)

const (
	envCronInterval = "CRON_INTERVAL_MIN"
	envHost         = "DB_HOST"
	envUser         = "DB_USER"
	envPass         = "DB_PASS"

	defaultIntervalMin = 1
	cronPattern        = "*/%v * * * *"
)

type service struct {
	interval int
	db       IngestorDB
	cron     *crontab.Crontab
	client   *http.Client
	parser   rss.Parser
}

func main() {
	log.Println("service started")

	var interval int
	var host, user, pass string

	flag.IntVar(&interval, envCronInterval, defaultIntervalMin, "cron interval (min) - using default otherwise")
	flag.StringVar(&host, envHost, "", "host:port for postgres")
	flag.StringVar(&user, envUser, "", "username for postgres")
	flag.StringVar(&pass, envPass, "", "password for postgres")

	flag.Parse()

	if len(user) == 0 {
		log.Fatal("required auth credentials not set. quitting.")
	}

	dbInstance, err := InitPostgresDB(host, user, pass)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to database.")

	srv := service{
		interval: interval,
		db:       dbInstance,
		cron:     crontab.New(),
		client:   &http.Client{},
		parser:   rss.Parser{},
	}

	srv.start()
}

func (srv service) start() {
	srv.cron.MustAddJob(fmt.Sprintf(cronPattern, srv.interval), srv.startProcessing)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	srv.db.close()
	srv.cron.Clear()

	log.Println("service stopped")
}

func (s service) startProcessing() {
	log.Println("process all feeds:")
	feeds, err := s.db.getFeeds()
	if err != nil {
		log.Println("error getting feeds: ", err)
	}
	for _, f := range feeds {
		go s.fetchAndProcessFeed(f)
	}
}

func (s service) fetchAndProcessFeed(f Feed) {
	log.Println("fetching ", f.URL)
	resp, err := s.client.Get(f.URL)
	if err != nil {
		log.Println("error fetching ", f.URL, err)
		return
	}
	defer resp.Body.Close()

	jsonFeed, err := s.parser.Parse(resp.Body)
	if err != nil {
		log.Println("error parsing feed: ", err)
		return
	}
	log.Println("found items #", len(jsonFeed.Items))
	for _, item := range jsonFeed.Items {
		news := News{
			Title:       item.Title,
			URL:         item.Link,
			GUID:        item.GUID.Value,
			Description: item.Description,
			Category:    f.Id,
			PubDate:     item.PubDateParsed,
		}
		err = s.db.storeNews(news)
		if err != nil {
			log.Println("error storing news: ", err)
		}
	}
}
