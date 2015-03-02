# goDownloader
Just a file downloader, what else do you want?

## Run
Just clone this repo and check the example:

```sh
# brainfuck here
GOMAXPROCS=2 go run main.go -u http://gioyik.com/index.html -d . -w 4
```
## Options
* Define what to download: `-u [url]`
* Where will be stored: `-d [dir]`
* How many workers will be running: `-w [number]`