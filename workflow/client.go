package workflow

import (
	"github.com/clarsen/trello"
	"github.com/sirupsen/logrus"
)

// Client wraps logged in member
type Client struct {
	Client *trello.Client
	Member *trello.Member
}

// Test does nothing
func (cl *Client) Test() {

}

// New create new client
func New(user string, appKey string, token string) (c *Client, err error) {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	client := trello.NewClient(appKey, token)
	client.Logger = logger

	member, err := client.GetMember(user, trello.Defaults())
	if err != nil {
		// Handle error
		return nil, err
	}
	client.Logger.Debugf("member %+v", member)
	c = &Client{
		Client: client,
		Member: member,
	}
	return
}
