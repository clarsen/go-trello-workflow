package workflow

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"
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

// MorningRemindHtml generates HTML for inspection as well as use in email
func MorningRemindHtml(user, appkey, authtoken string) (*WeeklySummary, string, error) {
	cl, err := New(user, appkey, authtoken)
	if err != nil {
		return nil, "", err
	}

	year, week := time.Now().ISOWeek()
	summary, err := GetSummaryForWeek(user, appkey, authtoken, year, week)
	if err != nil {
		return nil, "", err
	}

	// today cards
	list, err := listFor(cl.member, "Kanban daily/weekly", "Today")
	if err != nil {
		// handle error
		return nil, "", err
	}

	todayCards, err := list.GetCards(trello.Defaults())
	if err != nil {
		// handle error
		return nil, "", err
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
			return nil, "", err2
		}

		cards, err2 := list.GetCards(trello.Defaults())
		if err2 != nil {
			// handle error
			return nil, "", err2
		}
		for _, card := range cards {
			if card.Due != nil && now.After(card.Due.Add(-time.Hour*3*24)) {
				dueSoon = append(dueSoon, card)
			}
		}
	}

	sort.Sort(byDue(dueSoon))

	fmap := template.FuncMap{
		"formatAsDate": formatAsDate,
	}
	t := template.Must(template.New("morning-reminder.html").Funcs(fmap).ParseFiles("templates/morning-reminder.html"))

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
		return nil, "", err
	}
	return summary, content.String(), err
}

// MorningRemind generates and sends an email reminding of weekly goals, daily
// todos, overdue cards
func MorningRemind(user, appkey, authtoken, sendgridAPIKey, userEmail string) error {
	summary, content, err := MorningRemindHtml(user, appkey, authtoken)
	if err != nil {
		return err
	}
	from := mail.NewEmail("Daily reminder", userEmail)
	subject := fmt.Sprintf("Reminder for %d week %d", summary.Year, summary.Week)
	to := mail.NewEmail("Goal seeker", userEmail)

	mailcontent := mail.NewContent("text/html", content)
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
