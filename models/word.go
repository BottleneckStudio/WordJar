package models

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/BottleneckStudio/WordJar/cache"
	"github.com/BottleneckStudio/WordJar/cache/memcachier"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

type Word struct {
	Text           string          `json:"text"`
	Translation    string          `json:"translation,omitempty"`
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

type WordsAPIQuery struct {
	Text          string            `json:"word"`
	Definitions   []DefinitionQuery `json:"results"`
	Pronunciation map[string]string `json:"pronunciation"`
}

type WordsAPIError struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type DefinitionQuery struct {
	Text         string   `json:"definition"`
	PartOfSpeech string   `json:"partOfSpeech"`
	Synonyms     []string `json:"synonyms"`
	Examples     []string `json:"examples"`
}

type Pron struct {
	Pronunciation string `json:"pronunciation"`
}

type EntryListQuery struct {
	XMLName xml.Name `xml:"entry_list"`
	Entry   []Entry  `xml:"entry"`
}

type Entry struct {
	Sound Sound `xml:"sound"`
}

type Sound struct {
	Wav Wav `xml:"wav"`
}

type Wav struct {
	Content string `xml:",innerxml"`
}

type CrawlWordInput struct {
	Word   string
	Locale string
	Wg     sync.WaitGroup
	Ctx    *gin.Context
}

var memcachierService *memcachier.Memcachier

func init() {
	memcachierService = memcachier.NewMemcachier(memcachier.Config{Server: "mc5.dev.ec2.memcachier.com:11211", Username: "287FB9", Password: "C5054054E634497C354D72002A716CC4"})
}

func CrawlWord(input *CrawlWordInput) Word {
	w := Word{}
	var delta int
	var cacheKey string
	if input.Locale != "" {
		delta = 3
		cacheKey = input.Word + "." + input.Locale
	} else {
		delta = 2
		cacheKey = input.Word
	}
	w.Text = input.Word
	w.Pronunciations = []Pronunciation{}

	cacheData, cacheErr := cache.Get(cacheKey, memcachierService)

	if cacheErr == nil && cacheData != "" {
		json.Unmarshal([]byte(cacheData.(string)), &w)
		return w
	}

	// var wg sync.WaitGroup

	// input.Wg = wg
	input.Wg.Add(delta)
	go GetWord(&w, input)
	go GetAudio(&w, input)
	if input.Locale != "" {
		go GetTranslation(&w, input)
	}
	input.Wg.Wait()

	data, err := json.Marshal(w)

	if err != nil {
		log.Printf("Error marshalling due to: %v\n", err)
	} else {
		isCached, err := cache.Set(cacheKey, string(data), memcachierService)

		if !isCached {
			log.Printf("Error setting cache due to: %v\n", err)
		}
	}
	log.Printf("Successfully saved to cache using key: %s", cacheKey)
	return w
}

func GetWord(word *Word, input *CrawlWordInput) {
	defer input.Wg.Done()

	ctx := appengine.NewContext(input.Ctx.Request)
	client := urlfetch.Client(ctx)

	var apiURL = "https://wordsapiv1.p.mashape.com/words/"
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	req, err := http.NewRequest("GET", apiURL+word.Text, nil)
	if err != nil {
		return
	}
	req.Header.Set("X-Mashape-Key", "te6AX6SnBfmshawA0zj6VToSZO3up1MQySvjsnFmGv0qYDjUV3")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	b, err := ioutil.ReadAll(resp.Body)
	r := bytes.NewBuffer(b)

	var w WordsAPIQuery
	var p Pron

	var error WordsAPIError
	err = json.Unmarshal(b, &error)
	if error.Message == "word not found" {
		return
	}

	err = json.NewDecoder(r).Decode(&w)
	if err != nil {
		r = bytes.NewBuffer(b)
		err = json.NewDecoder(r).Decode(&p)
		if err != nil {
			log.Println("error: " + err.Error())
			return
		}
		word.Pronunciations = append(word.Pronunciations, Pronunciation{PartOfSpeech: "all", IPA: p.Pronunciation})
	}

	word.Created = time.Now()
	for _, v := range w.Definitions {
		for _, b := range v.Examples {
			word.Examples = append(word.Examples, b)
		}
		for _, b := range v.Synonyms {
			word.Synonyms = append(word.Synonyms, b)
		}
		word.Definitions = append(word.Definitions, Definition{PartOfSpeech: v.PartOfSpeech, Definition: v.Text})
	}

	for k, v := range w.Pronunciation {
		// log.Printf("key[%s] value[%s]\n", k, v)
		word.Pronunciations = append(word.Pronunciations, Pronunciation{PartOfSpeech: k, IPA: v})
	}
}

func GetAudio(word *Word, input *CrawlWordInput) {
	defer input.Wg.Done()
	ctx := appengine.NewContext(input.Ctx.Request)
	client := urlfetch.Client(ctx)

	var apiURL = "https://www.dictionaryapi.com/api/v1/references/collegiate/xml/"
	var apiKey = "?key=720750f6-2da7-4612-bb3e-2914b923052e"
	var baseAudioURL = "http://media.merriam-webster.com/soundc11/"

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	req, err := http.NewRequest("GET", apiURL+word.Text+apiKey, nil)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	log.Println("request URL: ", resp.Request.URL)

	var eq EntryListQuery
	err = xml.NewDecoder(resp.Body).Decode(&eq)
	if err != nil {
		return
	}
	if len(eq.Entry) == 0 {
		return
	}
	log.Println(eq.Entry[0].Sound.Wav.Content)
	var fileName = eq.Entry[0].Sound.Wav.Content
	var firstLetter = string(eq.Entry[0].Sound.Wav.Content[0])
	word.Audio = baseAudioURL + firstLetter + "/" + fileName
}

func GetTranslation(word *Word, input *CrawlWordInput) {
	defer input.Wg.Done()
	ctx := appengine.NewContext(input.Ctx.Request)
	client := urlfetch.Client(ctx)

	var apiURL = "https://translate.googleapis.com/translate_a/single?client=gtx&sl=en&tl=" + input.Locale + "&dt=t&q="
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	client = &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	req, err := http.NewRequest("GET", apiURL+word.Text, nil)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	log.Println("response Status:", resp.Status)
	log.Println("response Headers:", resp.Header)
	log.Println("request URL: ", resp.Request.URL)
	// log.Println("response Body: ", resp.Body)

	// var struc [][][]interface{}
	// err = json.NewDecoder(resp.Body).Decode(&struc)
	// if err != nil {
	// 	return
	// }
	var a interface{}
	err = json.NewDecoder(resp.Body).Decode(&a)
	if err != nil {
		return
	}

	log.Println(a)

	b := InterfaceSlice(a)
	log.Println(b[0])
	c := InterfaceSlice(b[0])
	log.Println(c[0])
	word.Translation = InterfaceSlice(c[0])[0].(string)
	// for i := range itemdata {
	// 	fmt.Println(itemdata[i]) // This prints '0', two times
	// 	for _ = range itemdata[i] {
	// 		fmt.Println(itemdata[i][0])
	// 		for _ = range itemdata[i][0] {
	// 			fmt.Println(itemdata[i][0][0])
	// 		}
	// 	}
	// }
	// var translation = itemdata[0]
	// word.Translation = translation.(string)
}

func InterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
