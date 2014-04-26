package main

import (
	"github.com/PuerkitoBio/goquery"
	"regexp"
	"time"
	//log "github.com/dgrijalva/timber"
	log "github.com/kdar/factorlog"
)

type Scraper struct {
	//---- Must be initialized
	Output   chan *story
	Interval int
	Name,
	Icon,
	Url,
	Target string
	//---- May be initialized
	Matcher,
	Excluder,
	Modifier,
	LinkModifier *scraperRegexp
	Unordered bool
	//---- Must not be initialized
	cache  map[string]int64
	last   [2]string
	active bool
	stop   chan struct{}
}

type scraperRegexp struct {
	*regexp.Regexp
}

func (s *scraperRegexp) UnmarshalText(text []byte) error {
	var err error
	s.Regexp, err = regexp.Compile(string(text))
	return err
}

func (s *Scraper) Start() {
	// Starting an active scraper should have no effect
	if s.active {
		log.Warn("Attempt to start an active scraper to ", s.Url)
		return
	}

	s.active = true

	s.stop = make(chan struct{})

	if s.cache == nil && s.Unordered {
		s.cache = make(map[string]int64)
	}

	s.initCache()

	//TODO: Goroutine within goroutine so that Stop always works immediately?
	go func() {
		pullTick := time.NewTicker(time.Duration(s.Interval) * time.Second)
		cacheTick := time.NewTicker(60 * time.Second)
		for {
			select {
			case <-pullTick.C:
				s.scrape()
			case <-cacheTick.C:
				s.cleanCache()
			case <-s.stop:
				break
			}
		}
	}()
}

func (s *Scraper) Stop() {
	// Stopping a stopped scraper should have no effect
	if !s.active {
		log.Warn("Attempt to stop an inactive scraper to ", s.Url)
		return
	}
	s.active = false
	close(s.stop)
}

func (s *Scraper) initCache() {
	log.Trace(`Initializing cache for `, s.Url)
	doc, err := goquery.NewDocument(s.Url)

	if err != nil {
		log.Errorf("Attempt to pull from %s failed: [%s]", s.Url, err.Error())
	} else {
		/*
		   The timestamp is used for unordered sources. Since items can appear in
		   any order, we keep a cache of all of the items we pulled the last time,
		   alerting one new ones. We assign a timestamp to each so that stale
		   entries can be removed from the cache and we don't just eat memory all
		   over the place
		*/
		timestamp := time.Now().Unix()

		items := doc.Find(s.Target)

		/*
		   If the source is ordered, we don't need to waste space on a cache.
		   Keeping track of the first two items we pulled from the last run is
		   enough. We pull two just in case the first item changes (as it sometimes
		   does in the fast-paced world of news.
		*/
		if !s.Unordered {
			s.last[0], _ = items.Slice(0, 1).Attr("href")
			s.last[1], _ = items.Slice(1, 2).Attr("href")
			log.Trace("Cached first two links for ", s.Url, ": [", s.last[0], "] [", s.last[1], "]")
		} else {
			log.Trace("Cacheing unordered source ", s.Url)
			items.Each(func(i int, sel *goquery.Selection) {
				//Extract the text and href of the link
				link, _ := sel.Attr("href")

				//Store the link in the cache
				log.Trace("\t [", link, "]")
				s.cache[link] = timestamp
			})
		}
	}
}

func (s *Scraper) scrape() {
	log.Trace(`Pulling from `, s.Url)
	doc, err := goquery.NewDocument(s.Url)
	if err != nil {
		log.Errorf("Attempt to pull from %s failed: [%s]", s.Url, err.Error())
	} else {
		timestamp := time.Now().Unix()

		items := doc.Find(s.Target)
		var firstTwo [2]string

		/*
		   If the source is unordered, we use this boolean to determine if we found
		   this item in the cache or not
		*/
		var found bool

		if !s.Unordered {
			firstTwo[0], _ = items.Slice(0, 0).First().Attr("href")
			firstTwo[1], _ = items.Slice(1, 1).First().Attr("href")
		}

		items.EachWithBreak(func(i int, sel *goquery.Selection) bool {

			//Extract the text and href of the link
			text := sel.Text()
			link, _ := sel.Attr("href")

			/* Stop processing if this is an ordered source and we hit one of the
			   first two from last time */
			if !s.Unordered {
				if link == s.last[0] || link == s.last[1] {
					return false
				}
			} else {
				// Update the cache for the current link if we're unordered
				_, found = s.cache[link]
				s.cache[link] = timestamp
			}

			//If we want to exclude certain links, run them against the excluder
			//TODO: Add support for href-based exclusion
			if s.Excluder != nil {
				if s.Excluder.MatchString(text) {
					log.Trace("Excluding link ", text)
					return true
				}
			}

			//If we want to modify all links, run them through the modifier
			//TODO: Make this more flexible
			if s.Modifier != nil {
				link = s.Modifier.ReplaceAllString(link, ``)
			}

			if (s.Unordered && !found) || !s.Unordered {
				s.Output <- &story{
					Source: s.Name,
					Icon:   s.Icon,
					Link:   link,
					Text:   text,
				}
			}
			return true
		})

		// Update last with the ones we just pulled
		if !s.Unordered {
			s.last = firstTwo
		}
	}
}

func (s *Scraper) cleanCache() {
	log.Debug("Cleaning the cache")
	now := time.Now().Unix()
	cleaned := 0
	for link, expiration := range s.cache {
		if expiration < now {
			cleaned++
			delete(s.cache, link)
		}
	}
	log.Debug("Removed ", cleaned, " entries")
}
