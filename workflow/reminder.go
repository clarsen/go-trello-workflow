package workflow

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"math"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/clarsen/trello"
	"github.com/gobuffalo/packr"

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

func durationUntilDue(t *time.Time) string {

	d := t.Sub(time.Now())
	var ret string
	if d < 0 { // overdue
		ret = "overdue"
	} else if d < 7*24*time.Hour {
		days := int64(math.Round(float64(d/time.Hour) / 24.0))
		ret = fmt.Sprintf("%dd", days)
	} else {
		weeks := int64(math.Round(float64(d/time.Hour) / (24.0 * 7)))
		ret = fmt.Sprintf("%dw", weeks)
	}
	return ret
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
		"formatAsDate":     formatAsDate,
		"durationUntilDue": durationUntilDue,
	}

	box := packr.NewBox("../templates")
	tmpl, err := box.FindString("morning-reminder.html")
	if err != nil {
		return nil, "", err
	}

	t := template.Must(template.New("morning-reminder.html").Funcs(fmap).Parse(tmpl))

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

func awsEmail(from, to, subject, htmlbody, textbody string) error {
	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	CharSet := "UTF-8"

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(htmlbody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(textbody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				log.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				log.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				log.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				log.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Println(err.Error())
		}
	}
	return err
}

func sendgridEmail(sendgridAPIKey, fromAddr, toAddr, subject, htmlbody, textbody string) error {
	from := mail.NewEmail("Daily reminder", fromAddr)
	to := mail.NewEmail("Goal seeker", toAddr)

	mailcontent := mail.NewContent("text/html", htmlbody)
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

// MorningRemind generates and sends an email reminding of weekly goals, daily
// todos, overdue cards
func MorningRemind(user, appkey, authtoken, sendgridAPIKey, userEmail string) error {
	summary, content, err := MorningRemindHtml(user, appkey, authtoken)
	if err != nil {
		return err
	}
	subject := fmt.Sprintf("Reminder for %d week %d", summary.Year, summary.Week)

	if sendgridAPIKey != "" {
		return sendgridEmail(sendgridAPIKey, userEmail, userEmail, subject, content, "This contains HTML")
	} else {
		return awsEmail(userEmail, userEmail, subject, content, "This contains HTML")
	}
}
