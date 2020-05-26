package workflow

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

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

func ChecklistTitleFromAttributes(title string, week int, created *time.Time, estDuration *string, status *string) string {
	cstr := ""
	if created != nil {
		cstr = created.Format(" (2006-01-02)")
	}
	e := ""
	if estDuration != nil {
		e = fmt.Sprintf("%s ", *estDuration)
	}
	s := ""
	if status != nil {
		s = fmt.Sprintf(" %s", *status)
	}
	return fmt.Sprintf("week %d: %s%s%s%s", week, e, title, cstr, s)
}

func GetAttributesFromChecklistTitle(name string) (title string, week int, created *time.Time, estDuration *string, status *string) {
	re := regexp.MustCompile(`^[Ww]eek (\d+): ((\(\d+.*\)) )?(.*) (?:\((\d{4}-\d{2}-\d{2})\))(?: (\(.*\)))?\s*$`)
	groups := re.FindSubmatch([]byte(name))
	if groups == nil {
		log.Printf("match '%+v' didn't match\n", name)
		return
	}
	week, err := strconv.Atoi(string(groups[1]))
	if err != nil {
		week = 0
	}
	if len(groups[3]) > 0 {
		d := string(groups[3])
		estDuration = &d
	}
	title = string(groups[4])
	if len(groups[5]) > 0 {
		c, err := time.Parse("2006-01-02", string(groups[5]))
		if err == nil {
			created = &c
		}
	}
	if len(groups[6]) > 0 {
		s := string(groups[6])
		status = &s
	}
	return
}

func parseCardName(name string) (title string, created *string, period *string) {
	// https://regex101.com/r/UeBSfi/1
	re := regexp.MustCompile(`^(.*?)(?:\s+\((\d{4}-\d{2}-\d{2})\))?(?:\s*\((po|p1w|p2w|p4w|p2m|p3m|p6m|p12m)\))?$`)
	groups := re.FindSubmatch([]byte(name))
	if groups == nil {
		log.Printf("match %+v didn't match\n", name)
		return
	}

	title = string(groups[1])
	if len(groups[2]) > 0 {
		c := string(groups[2])
		created = &c
	}
	if len(groups[3]) > 0 {
		p := string(groups[3])
		period = &p
	}
	log.Printf("match %+v -> %s, %s, %s\n", name, title, string(groups[2]), string(groups[3]))
	return title, created, period
}

func GetTitleAndAttributes(card *trello.Card) (title string, created *string, period *string) {
	return parseCardName(card.Name)
}

func isPeriodic(card *trello.Card) bool {
	re := regexp.MustCompile("\\((po|p1w|p2w|p4w|p2m|p3m|p6m|p12m)\\)")
	period := re.FindString(card.Name)
	return period != ""
}

func (cl *Client) GetCard(cardId string) (*trello.Card, error) {
	card, err := cl.Client.GetCard(cardId, trello.Arguments{"fields": "all"})
	if err != nil {
		return nil, err
	}
	return card, nil
}

func (cl *Client) CreateCard(title string, board string, list string) (*trello.Card, error) {
	l, err := ListFor(cl, board, list)
	if err != nil {
		return nil, err
	}
	created := time.Now()
	cstr := created.Format(" (2006-01-02)")
	card := trello.Card{
		Name:   fmt.Sprintf("%s%s", title, cstr),
		Pos:    0.0,
		IDList: l.ID,
	}
	err = cl.Client.CreateCard(&card, trello.Arguments{"pos": "top"})
	if err != nil {
		return nil, err
	}
	return &card, err
}
