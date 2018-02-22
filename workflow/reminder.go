package workflow

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/clarsen/trello"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func formatAsDate(t *time.Time) string {
	local, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return "couldn't load timezone"
	}
	return t.In(local).Format("2006-01-02 (Mon)")
}

// MorningRemind generates and sends an email reminding of weekly goals, daily
// todos, overdue cards
func MorningRemind(user, appkey, authtoken, sendgridAPIKey, userEmail string) error {
	cl, err := New(user, appkey, authtoken)
	if err != nil {
		return err
	}

	summary, err := GetSummaryForWeek(user, appkey, authtoken)
	if err != nil {
		return err
	}

	// today cards
	list, err := listFor(cl.member, "Kanban daily/weekly", "Today")
	if err != nil {
		// handle error
		return err
	}

	todayCards, err := list.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return err
	}

	// overdue check
	overdueBoardsAndLists := []boardAndList{
		{"Kanban daily/weekly", "Today"},
		{"Kanban daily/weekly", "Waiting on"},
		{"Backlog (Personal)", "Backlog"},
		{"Backlog (Personal)", "Projects"},
		{"Backlog (work)", "Backlog"},
		{"Periodic board", "Often"},
		{"Periodic board", "Weekly"},
		{"Periodic board", "Bi-weekly to monthly"},
		{"Periodic board", "Quarterly to Yearly"},
	}

	now := time.Now()

	var dueSoon []*trello.Card
	for _, boardlist := range overdueBoardsAndLists {
		list, err2 := listFor(cl.member, boardlist.Board, boardlist.List)
		if err2 != nil {
			// handle error
			return err2
		}

		cards, err2 := list.GetCards(trello.Defaults())
		if err2 != nil {
			// handle error
			return err2
		}
		for _, card := range cards {
			if card.Due != nil && now.After(card.Due.Add(-time.Hour*3*24)) {
				dueSoon = append(dueSoon, card)
			}
		}
	}

	from := mail.NewEmail("Daily reminder", userEmail)
	subject := fmt.Sprintf("Reminder for %d week %d", summary.Year, summary.Week)
	to := mail.NewEmail("Goal seeker", userEmail)

	fmap := template.FuncMap{
		"formatAsDate": formatAsDate,
	}
	t := template.Must(template.New("morning-reminder.txt").Funcs(fmap).ParseFiles("templates/morning-reminder.txt"))

	data := struct {
		Summary *WeeklySummary
		Today   []*trello.Card
		DueSoon []*trello.Card
	}{
		summary,
		todayCards,
		dueSoon,
	}

	var content bytes.Buffer
	err = t.Execute(&content, data)
	if err != nil {
		return err
	}
	// for testing
	// log.Println(content.String())
	// return nil

	mailcontent := mail.NewContent("text/plain", content.String())
	m := mail.NewV3MailInit(from, subject, to, mailcontent)

	request := sendgrid.GetRequest(sendgridAPIKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	response, err := sendgrid.API(request)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}

	return err
}
