package scraper

import (
	"io/ioutil"
	"log"

	"github.com/PuerkitoBio/goquery"
	"strings"
)

// Scraper struct used for scraping activity.
type Scraper struct {
	Config  *Config
	Results results
}

// results key is tag, value is list of scraped data
type results map[string][]string

type resultsURL struct {
	url     string
	results results
}

// New creates new Scraper and returns pointer to it, with error (if occurred).
func New(configPath string) (*Scraper, error) {
	c, err := newConfig(configPath)
	if err != nil {
		return nil, err
	}

	s := &Scraper{
		Config:  c,
		Results: make(results),
	}

	return s, nil
}

// Start will start scraping in separate goroutine and save results
// when it is done in file that is defined in Config.
func (s *Scraper) Start() {
	resultsURLCh := s.scrapeURLs(s.Config.URLs)

	for range s.Config.URLs {
		resultsURL := <-resultsURLCh
		mergeResults(s.Results, resultsURL.results)
		log.Println("Received results from", resultsURL.url)
	}

	log.Println("Done scraping.")

	err := s.save()
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Saved to", s.Config.outputPath())
}

func (s *Scraper) scrapeURLs(urls []string) <-chan *resultsURL {
	resultsURLCh := make(chan *resultsURL)

	for _, url := range urls {
		go func(url string) {
			results := make(results)

			// Once individual URL scraping is done, send results back through channel
			defer func() {
				resultsURLCh <- &resultsURL{url: url, results: results}
			}()

			// Construct document for manipulation
			doc, err := goquery.NewDocument(url)
			if err != nil {
				log.Println(err)
				return
			}

			// Process targets
			for _, target := range s.Config.Targets {
				err := s.processTarget(doc, target, results)
				if err != nil {
					log.Println("Error processing target:", err)
				}
			}
		}(url)
	}

	return resultsURLCh
}

func (s *Scraper) processTarget(doc *goquery.Document, target *target, results results) error {
	var selector string
	// If there was no selector given, whole document will be used
	if target.Selector == "" {
		selector = "html"
	} else {
		selector = target.Selector
	}

	var retErr error
	doc.Find(selector).Each(func(i int, sel *goquery.Selection) {
		var value string

		// Handling different types
		switch {
		// Sets value to inner HTML of the node
		case target.Type == "html":
			html, err := sel.Html()
			if err != nil {
				retErr = err
				return
			}
			value = html
		// Sets value to text of the node
		case target.Type == "text":
			value = sel.Text()
		// Sets value to attribute of the node, for example attr:href for href value
		// If attribute value doesn't exist, target is skipped
		case strings.HasPrefix(target.Type, "attr:"):
			attrTarget := strings.Split(target.Type, ":")[1]
			if attrv, exists := sel.Attr(attrTarget); exists {
				value = attrv
			} else {
				return
			}
		}

		// Submatch regex
		if target.submatchRe != nil {
			matches := target.submatchRe.FindAllStringSubmatch(value, -1)
			for _, match := range matches {
				results[target.Tag] = append(results[target.Tag], match[0])
			}
			return
		}

		results[target.Tag] = append(results[target.Tag], value)
	})

	return retErr
}

func (s *Scraper) save() error {
	data, err := JSONMarshalUnescaped(s.Results)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(s.Config.outputPath(), data, 0644)
	if err != nil {
		return err
	}

	return nil
}
