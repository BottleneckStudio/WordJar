package models

import (
	"fmt"
	"log"
	"time"

	"github.com/gocolly/colly"
)

type Word struct {
	Text           string          `json:"text"`
	Translation    string          `json:"translation"`
	Audio          string          `json:"audio"`
	Definitions    []Definition    `json:"definitions"`
	Pronunciations []Pronunciation `json:"pronunciations"`
	Examples       []string        `json:"examples"`
	Synonyms       []string        `json:"synonyms"`
	Created        time.Time       `json:"created"`
}

type Definition struct {
	Definition   string `json:"definition"`
	PartOfSpeech string `json:"partOfSpeech"`
}

type Pronunciation struct {
	PartOfSpeech string `json:"partOfSpeech"`
	IPA          string `json:"IPA"`
}

func CrawlWord(word string) Word {
	w := Word{}

	w.Text = word
	w.Pronunciations = []Pronunciation{}

	c := colly.NewCollector(
		// Visit only domains: https://m-w.com
		colly.AllowedDomains("www.merriam-webster.com"),
		colly.Async(true),
	)

	c.OnHTML("div.full-def-box.def-header-box.card-box.def-text.headword-box.show-collapsed", func(h *colly.HTMLElement) {
		pron := Pronunciation{}
		pron.PartOfSpeech = h.ChildText("a.important-blue-link")
		pron.IPA = h.ChildText("span.mw")
		w.Pronunciations = append(w.Pronunciations, pron)
		fmt.Println("Word:", w)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://www.merriam-webster.com/dictionary/" + word)

	c.Wait()
	return w
}
