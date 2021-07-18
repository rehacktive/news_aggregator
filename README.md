# news-aggregator

The news-aggregator, based on 2 microservices, allows the user to define a list of sources (RSS atom feeds) and fetch them periodically, storing the content in the database, and a way to fetch the content, previously stored, for the original requirement: an app showing news from a list of public feeds.

It uses Postgres (containerized) as database.

The 2 microservices:

* **feedIngestionService**: periodically gets a list of feed URLs from the database, gets the articles from all the sources and stores them in the database;

  The list of feeds is stored in the database, and a category is associated to each different source, allowing later the possibility to filter news by category.

  The fetch interval can be specified as an env variable, default is 1 minute.

  *This  feed list could be managed in a CRUD-way by the endpoint (not implemented)*

* **newsfeedService**: exposes 3 different endpoints, related to the original requirements, to get a list of news stored in the database, filter them by category and get a list of categories (used by the previous filter).

---

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
git clone https://github.com/rehacktive/news_aggregator.git
```

## How to run it

If you have docker and docker-compose, you can just:

```sh
docker-compose up
```

Another option is to install  Postgres locally (port 5432 on localhost) and then run both services with the *run.sh* script present in both folders.



---

## How to use the endpoints

Once both services are up and running, after the defined interval (default is 1 minute) the database should be populated with the articles from the predefined feeds.

When the content is available on the database, the first endpoint returns a list of news, sorted by date:

```sh
$ curl http://localhost:8880/news | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 19934    0 19934    0     0  9733k      0 --:--:-- --:--:-- --:--:-- 9733k
{
  "news": [
    {
      "Id": 21,
      "Title": "Covid-19: PM and chancellor self-isolate after rapid U-turn",
      "URL": "https://www.bbc.co.uk/news/uk-57879730",
      "GUID": "https://www.bbc.co.uk/news/uk-57879730",
      "Description": "It comes after anger over \"VIP testing\" allowing them to work following Sajid Javid's positive test.",
      "Category": 1,
      "PubDate": "2021-07-18T17:27:23Z"
    },
    {
      "Id": 112,
      "Title": "Cavendish wins Tour de France green jersey - but narrowly misses out on record 35 stage wins",
      "URL": "http://news.sky.com/story/tour-de-france-mark-cavendish-wins-green-jersey-but-narrowly-misses-out-on-record-35-stage-wins-12358666",
      "GUID": "http://news.sky.com/story/tour-de-france-mark-cavendish-wins-green-jersey-but-narrowly-misses-out-on-record-35-stage-wins-12358666",
      "Description": "British cyclist Mark Cavendish has won the green jersey in this year's Tour de France.",
      "Category": 2,
      "PubDate": "2021-07-18T16:57:00Z"
    },
    ...
```

To "resolve" the category name another endpoint is available:

```sh
$ curl http://localhost:8880/category | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    70  100    70    0     0  70000      0 --:--:-- --:--:-- --:--:-- 70000
{
  "categories": [
    {
      "Id": 1,
      "Name": "BBC News"
    },
    {
      "Id": 2,
      "Name": "Sky News"
    }
  ]
}

```

it's also now possible use the category ID to filter the news:

```sh
$ curl http://localhost:8880/news/2 | jq .
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100  6453    0  6453    0     0  3150k      0 --:--:-- --:--:-- --:--:-- 3150k
{
  "news": [
    {
      "Id": 112,
      "Title": "Cavendish wins Tour de France green jersey - but narrowly misses out on record 35 stage wins",
      "URL": "http://news.sky.com/story/tour-de-france-mark-cavendish-wins-green-jersey-but-narrowly-misses-out-on-record-35-stage-wins-12358666",
      "GUID": "http://news.sky.com/story/tour-de-france-mark-cavendish-wins-green-jersey-but-narrowly-misses-out-on-record-35-stage-wins-12358666",
      "Description": "British cyclist Mark Cavendish has won the green jersey in this year's Tour de France.",
      "Category": 2,
      "PubDate": "2021-07-18T16:57:00Z"
    },
    {
      "Id": 3,
      "Title": "Hamilton wins British Grand Prix after collision with Verstappen who was taken to hospital",
      "URL": "http://news.sky.com/story/lewis-hamilton-wins-british-grand-prix-after-collision-with-max-verstappen-who-was-taken-to-hospital-12358585",
      "GUID": "http://news.sky.com/story/lewis-hamilton-wins-british-grand-prix-after-collision-with-max-verstappen-who-was-taken-to-hospital-12358585",
      "Description": "Lewis Hamilton has won the British Grand Prix following a race which saw him and F1 title rival Max Verstappen collide on the opening lap.",
      "Category": 2,
      "PubDate": "2021-07-18T15:58:00Z"
    },
    {
      "Id": 5,
      "Title": "Man arrested after teenager riding e-scooter dies following hit-and-run",
      "URL": "http://news.sky.com/story/teenager-riding-e-scooter-dies-after-hit-and-run-12358583",
      "GUID": "http://news.sky.com/story/teenager-riding-e-scooter-dies-after-hit-and-run-12358583",
      "Description": "A man has been arrested after a 16-year-old boy was killed in a hit-and-run incident while riding an e-scooter in south-east London, police say.",
      "Category": 2,
      "PubDate": "2021-07-18T15:58:00Z"
    },
    ...
```

## 





