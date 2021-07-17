package main

type MockIngestorDB struct {
	feeds []Feed
	news  []News
}

func (m *MockIngestorDB) getFeeds() ([]Feed, error) {
	return m.feeds, nil
}

func (m *MockIngestorDB) storeNews(n News) error {
	m.news = append(m.news, n)
	return nil
}
