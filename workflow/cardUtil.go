package workflow

import "github.com/clarsen/trello"

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
