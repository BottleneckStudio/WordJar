package models

import "time"

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
