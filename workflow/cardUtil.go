package workflow

import (
	"log"
	"regexp"

	"github.com/clarsen/trello"
)

type byDue []*trello.Card

func (c byDue) Len() int {
	return len(c)
}

func (c byDue) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c byDue) Less(i, j int) bool {
	if c[i].Due != nil && c[j].Due != nil {
		return c[i].Due.Before(*c[j].Due)
	} else if c[i].Due != nil && c[j].Due == nil {
		return true
	} else if c[i].Due == nil && c[j].Due != nil {
		return false
	}
	return false
}

func hasDate(card *trello.Card) bool {
	re := regexp.MustCompile("\\(\\d{4}-\\d{2}-\\d{2}\\)")
	date := re.FindString(card.Name)
	return date != ""
}

func GetTitleAndAttributes(card *trello.Card) (title string, created *string, period *string) {
	re := regexp.MustCompile(`(.*)\s+?(\((\d{4}-\d{2}-\d{2})\))?\s*(\((po|p1w|p2w|p4w|p2m|p3m|p6m|p12m)\))?`)
	groups := re.FindSubmatch([]byte(card.Name))
	if groups == nil {
		log.Printf("match %+v didn't match\n", card.Name)
		return
	}
	title = string(groups[1])
	if len(groups[3]) > 0 {
		c := string(groups[3])
		created = &c
	}
	if len(groups[5]) > 0 {
		p := string(groups[5])
		period = &p
	}
	log.Printf("match %+v -> %s, %s, %s\n", card.Name, title, string(groups[3]), string(groups[5]))
	return title, created, period
}

func isPeriodic(card *trello.Card) bool {
	re := regexp.MustCompile("\\((po|p1w|p2w|p4w|p2m|p3m|p6m|p12m)\\)")
	period := re.FindString(card.Name)
	return period != ""
}
