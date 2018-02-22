package workflow

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/clarsen/trello"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

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

	from := mail.NewEmail("Daily reminder", userEmail)
	subject := fmt.Sprintf("Reminder for %d week %d", summary.Year, summary.Week)
	to := mail.NewEmail("Goal seeker", userEmail)
	var content bytes.Buffer
	t, err := template.ParseFiles("templates/morning-reminder.txt")
	if err != nil {
		return err
	}

	data := struct {
		Summary *WeeklySummary
		Today   []*trello.Card
	}{
		summary,
		todayCards,
	}

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
