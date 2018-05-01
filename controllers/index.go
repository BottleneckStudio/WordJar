package controllers

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"

	"github.com/BottleneckStudio/WordJar/models"
)

func IndexController(c *gin.Context) {
	OutputJSON(c, "ok", "Welcome to Index")
}

func WordController(c *gin.Context) {
	word := c.Param("word")
	result := CrawlWord(word)
	data := gin.H{
		"result": result,
	}
	OutputDataAsJSON(c, data, "ok", "1 result")
}

func CrawlWord(word string) models.Word {
	w := models.Word{}

	w.Text = word
	w.Pronunciations = []models.Pronunciation{}

	c := colly.NewCollector(
		// Visit only domains: https://m-w.com
		colly.AllowedDomains("www.merriam-webster.com"),
		colly.Async(true),
	)

	c.OnHTML("div.full-def-box.def-header-box.card-box.def-text.headword-box.show-collapsed", func(h *colly.HTMLElement) {
		pron := models.Pronunciation{}
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
