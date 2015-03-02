package main

// imports
import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

const (
	QSIZE = 10
)

// options to parse with the command
type Opt struct {
	url     string
	target  string
	workers int
}

type Ln struct {
	filename string
	url      string
}

func main() {
	var wg sync.WaitGroup
	opts := params()

	fmt.Printf("Downloading %s\n", opts.url)
	bytes, err := fetch(opts.url)
	if err != nil {
		panic(err)
	}

	links := extract(bytes, opts)
	fmt.Printf("%d files\n", len(links))

	queue := make(chan Ln, QSIZE)
	for i := 0; i < opts.workers; i++ {
		wg.Add(1)
		go worker(i+1, queue, opts, &wg)
	}

	for _, link := range links {
		queue <- link
	}

	fmt.Println("Closing...")
	close(queue)
	wg.Wait()
}

func worker(index int, queue <-chan Ln, opts Opt, wg *sync.WaitGroup) {
	defer wg.Done()
	for link := range queue {
		fmt.Printf("Worker %d, downloading %s\n", index, link.url)
		bytes, err := fetch(link.url)

		if err != nil {
			fmt.Println(err)
			continue
		}

		ioutil.WriteFile(opts.target+link.filename, bytes, 0644)
	}
}

// Fetch function
func fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("Failed to fetch " + url)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Failed to parse " + url)
	}
	return body, nil
}

func extract(html []byte, opts Opt) []Ln {
	fmt.Printf("Detecting subpaths...\n")
	r := regexp.MustCompile("(?i)<td align=top><a href=\"([^\"]+)\">")
	paths := r.FindAllSubmatch(html, -1)

	var links []Ln
	for _, i := range paths {
		path := string(i[1])
		links = append(links, Ln{filename: path, url: opts.url + path})
	}

	return links
}

// params to parse from Opt in the command line
func params() Opt {
	url := flag.String("u", "", "URL")
	target := flag.String("d", "/tmp", "Target directory")
	workers := flag.Int("w", 2, "Number of workers")
	flag.Parse()

	return Opt{url: *url, target: *target, workers: *workers}
}
