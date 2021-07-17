package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mileusna/crontab"
	"github.com/mmcdole/gofeed/rss"
)

const (
	feedValue = `<?xml version="1.0" encoding="UTF-8"?>
	<rss xmlns:atom="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" version="2.0">
	  <channel>
		<atom:link href="https://feeds.skynews.com/feeds/rss/uk.xml" rel="self" type="application/rss+xml"/>
		<title>UK News - The latest headlines from the UK | Sky News</title>
		<link>https://news.sky.com/uk</link>
		<image>
		  <title>UK News - The latest headlines from the UK | Sky News</title>
		  <url>https://e3.365dm.com/skynews/logo.png</url>
		  <link>https://news.sky.com/uk</link>
		</image>
		<description>Expert comment and analysis on the latest UK news, with headlines from England, Scotland, Northern Ireland and Wales.</description>
		<language>en-GB</language>
		<copyright>Copyright 2021, Sky UK. All Rights Reserved.</copyright>
		<lastBuildDate>Sat, 17 Jul 2021 13:35:00 +0100</lastBuildDate>
		<category>Sky News</category>
		<ttl>1</ttl>
		<item>
		  <title>Due your COVID-19 jab? Pop-up vaccine clinics open at shops, parks, and stadiums this weekend</title>
		  <link>http://news.sky.com/story/covid-19-pop-up-vaccine-clinics-open-at-shops-parks-and-stadiums-this-weekend-12357592</link>
		  <description>Pop-up clinics will be offering coronavirus jabs this weekend as part of a major drive to get more people vaccinated, so those enjoying the sunshine by going shopping or visiting the park will be able to get one on the go.</description>
		  <pubDate>Sat, 17 Jul 2021 10:04:00 +0100</pubDate>
		  <guid>http://news.sky.com/story/covid-19-pop-up-vaccine-clinics-open-at-shops-parks-and-stadiums-this-weekend-12357592</guid>
		  <enclosure url="http://e3.365dm.com/21/07/70x70/skynews-covid-coronavirus-vaccine_5449859.jpg?20210717115239" length="0" type="image/jpeg"/>
		  <media:description type="html">Behind enemy lines: An Arsenal fan gets their COVID vaccine at the Tottenham Hotspur Stadium</media:description>
		  <media:thumbnail url="http://e3.365dm.com/21/07/70x70/skynews-covid-coronavirus-vaccine_5449859.jpg?20210717115239" width="70" height="70"/>
		  <media:content type="image/jpeg" url="http://e3.365dm.com/21/07/70x70/skynews-covid-coronavirus-vaccine_5449859.jpg?20210717115239"/>
		</item>
		<item>
		  <title>Health Secretary Sajid Javid tests positive for coronavirus and has 'mild' symptoms </title>
		  <link>http://news.sky.com/story/health-secretary-sajid-javid-tests-positive-for-coronavirus-and-has-mild-symptoms-12357736</link>
		  <description>Health Secretary Sajid Javid has tested positive for coronavirus and is experiencing &quot;mild&quot; symptoms.</description>
		  <pubDate>Sat, 17 Jul 2021 13:35:00 +0100</pubDate>
		  <guid>http://news.sky.com/story/health-secretary-sajid-javid-tests-positive-for-coronavirus-and-has-mild-symptoms-12357736</guid>
		  <enclosure url="http://e3.365dm.com/21/07/70x70/skynews-health-secretary-sajid-javid_5438856.jpg?20210706130113" length="0" type="image/jpeg"/>
		  <media:description type="html">The Health Secretary Sajid Javid speaks to the House of Commons</media:description>
		  <media:thumbnail url="http://e3.365dm.com/21/07/70x70/skynews-health-secretary-sajid-javid_5438856.jpg?20210706130113" width="70" height="70"/>
		  <media:content type="image/jpeg" url="http://e3.365dm.com/21/07/70x70/skynews-health-secretary-sajid-javid_5438856.jpg?20210706130113"/>
		</item>
		<item>
		  <title>Wales eases COVID rules - but full unlocking still weeks away despite lowest case rate in UK</title>
		  <link>http://news.sky.com/story/covid-19-wales-eases-restrictions-again-but-full-unlocking-still-weeks-away-despite-lowest-case-rate-in-uk-12357520</link>
		  <description>Wales has further eased coronavirus restrictions, with up to six people now allowed to meet inside homes or holiday accommodation.</description>
		  <pubDate>Sat, 17 Jul 2021 03:32:00 +0100</pubDate>
		  <guid>http://news.sky.com/story/covid-19-wales-eases-restrictions-again-but-full-unlocking-still-weeks-away-despite-lowest-case-rate-in-uk-12357520</guid>
		  <enclosure url="http://e3.365dm.com/20/09/70x70/skynews-wales-lockdown-coronavirus_5113460.jpg?20200930085559" length="0" type="image/jpeg"/>
		  <media:description type="html"> woman wearing a face covering walks past a nightclub in Cardiff</media:description>
		  <media:thumbnail url="http://e3.365dm.com/20/09/70x70/skynews-wales-lockdown-coronavirus_5113460.jpg?20200930085559" width="70" height="70"/>
		  <media:content type="image/jpeg" url="http://e3.365dm.com/20/09/70x70/skynews-wales-lockdown-coronavirus_5113460.jpg?20200930085559"/>
		</item>
	  </channel>
	</rss>
	`
)

func TestFetchingFeed(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(feedValue))
	}))
	defer server.Close()

	feeds := []Feed{
		{
			URL:      server.URL,
			Category: "sample feed",
		},
	}

	mockedDb := MockIngestorDB{
		feeds: feeds,
		news:  make([]News, 0),
	}

	srv := service{
		db:     &mockedDb,
		cron:   crontab.New(),
		client: &http.Client{},
		parser: rss.Parser{},
	}

	srv.startProcessing()

	time.Sleep(1 * time.Second)

	if len(mockedDb.news) != 3 {
		t.Error("expected 3 news, found ", len(mockedDb.news))
	}

}
