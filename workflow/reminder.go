package workflow

import (
	"bytes"
	"fmt"
	"html/template"

	sendgrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

func MorningRemind(user, appkey, authtoken, sendgridAPIKey, userEmail string) error {
	summary, err := GetSummaryForWeek(user, appkey, authtoken)
	if err != nil {
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
	err = t.Execute(&content, summary)
	if err != nil {
		return err
	}

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
