package workflow

import (
	"regexp"
	"strconv"

	"github.com/clarsen/trello"
)

type byDescWeek []trello.CheckItem

func (c byDescWeek) Len() int {
	return len(c)
}

func (c byDescWeek) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c byDescWeek) Less(i, j int) bool {
	re := regexp.MustCompile(`^week (\d+):.*$`)
	exprI := re.FindStringSubmatch(c[i].Name)
	var weekI, weekJ *int

	if len(exprI) > 0 {
		// log.Printf("for %s got match %+v\n", item.Name, expr)
		// log.Println("got week", expr[1])
		// log.Println("got text", expr[2])
		// log.Println("got status", expr[3])
		week, err := strconv.Atoi(exprI[1])
		if err != nil {
			return false
		}
		weekI = &week
	}
	exprJ := re.FindStringSubmatch(c[j].Name)
	if len(exprJ) > 0 {
		// log.Printf("for %s got match %+v\n", item.Name, expr)
		// log.Println("got week", expr[1])
		// log.Println("got text", expr[2])
		// log.Println("got status", expr[3])
		week, err := strconv.Atoi(exprJ[1])
		if err != nil {
			return false
		}
		weekJ = &week
	}
	if weekI != nil && weekJ != nil {
		return *weekI > *weekJ
	}
	return false
}
