package main

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// db interface and models

type IngestorDB interface {
	getFeeds() ([]Feed, error)
	storeNews(n News) error
	close()
}

type Feed struct {
	Id       int64
	URL      string
	Category string
}

type News struct {
	Id          int64
	Title       string
	URL         string
	GUID        string
	Description string
	Category    int64
	PubDate     *time.Time
}

// pg implementation

type IngestorPostgresDB struct {
	db *pg.DB
}

func InitPostgresDB(host string, user string, pass string) (dbInstance IngestorPostgresDB, err error) {
	db := pg.Connect(&pg.Options{
		Addr:     host,
		User:     user,
		Password: pass,
	})

	err = createSchema(db)
	if err != nil {
		return
	}

	// HACK :) REMOVE ME LATER
	// adding some feeds for demo purpose
	feed1 := &Feed{
		Id:       1,
		URL:      "http://feeds.bbci.co.uk/news/uk/rss.xml",
		Category: "BBC News",
	}
	_, err = db.Model(feed1).
		OnConflict("(id) DO UPDATE").
		Insert()
	if err != nil {
		return
	}
	feed2 := &Feed{
		Id:       2,
		URL:      "http://feeds.skynews.com/feeds/rss/uk.xml",
		Category: "Sky News",
	}
	_, err = db.Model(feed2).
		OnConflict("(id) DO UPDATE").
		Insert()
	if err != nil {
		return
	}
	// end REMOVE ME

	dbInstance = IngestorPostgresDB{
		db: db,
	}
	return
}

func (idb IngestorPostgresDB) close() {
	idb.db.Close()
}

func (idb IngestorPostgresDB) getFeeds() (feeds []Feed, err error) {
	err = idb.db.Model(&feeds).Select()
	return
}

func (idb IngestorPostgresDB) storeNews(n News) (err error) {
	_, err = idb.db.Model(&n).
		OnConflict("(guid) DO UPDATE").
		Insert()
	return
}

func createSchema(db *pg.DB) (err error) {
	models := []interface{}{
		(*Feed)(nil),
		(*News)(nil),
	}

	for _, model := range models {
		err = db.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists: true,
		})
		if err != nil {
			return
		}
	}
	_, err = db.Model(&News{}).Exec(`
    		CREATE UNIQUE INDEX guid_idx
    		ON ?TableName(guid)
	`)

	return
}
