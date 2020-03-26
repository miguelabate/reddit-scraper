#Reddit Scraper

Simple tool to fetch all the Urls from the posts of a given SubReddit

#Build
```
go build reddit-scraper.go
```

#Usage

```
Options:
  -h, --help            Display usage
  -v, --version         Display version
      --vv              Display version (extended)

Commands:
  urls                  fetch urls from post
    -s, --sub           Subreddit to scrap
    -o, --output        Output file. If not present, send to STDOUT
    -v, --verbose       Show extra information during the process and stats
```
