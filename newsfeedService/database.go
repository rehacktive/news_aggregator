package main

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

type FeedDB interface {
	getCategories() ([]Category, error)
	getNews() ([]News, error)
	getNewsByCategory(category string) ([]News, error)
	close()
}

type Category struct {
	Id   int64
	Name string
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
	GUID        string `pg:",unique"`
	Description string
	Category    int64
	PubDate     *time.Time
}

type FeedPostgresDB struct {
	db *pg.DB
}

func InitPostgresDB(host string, user string, pass string) (dbInstance FeedPostgresDB, err error) {
	db := pg.Connect(&pg.Options{
		Addr:     host,
		User:     user,
		Password: pass,
	})

	err = createSchema(db)
	if err != nil {
		return
	}

	dbInstance = FeedPostgresDB{
		db: db,
	}
	return
}

func (idb FeedPostgresDB) close() {
	idb.db.Close()
}

func (idb FeedPostgresDB) getCategories() (ret []Category, err error) {
	var feeds []Feed
	err = idb.db.Model(&feeds).
		Select()
	if err != nil {
		return
	}
	for _, n := range feeds {
		ret = append(ret, Category{
			Id:   n.Id,
			Name: n.Category,
		})
	}
	return
}

func (idb FeedPostgresDB) getNews() (news []News, err error) {
	err = idb.db.Model(&news).
		Order("pub_date DESC").
		Select()
	return
}

func (idb FeedPostgresDB) getNewsByCategory(category string) (news []News, err error) {
	err = idb.db.Model(&news).
		Order("pub_date DESC").
		Where("category = ?", category).
		Select()
	return
}

func createSchema(db *pg.DB) (err error) {
	models := []interface{}{
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
	return
}
