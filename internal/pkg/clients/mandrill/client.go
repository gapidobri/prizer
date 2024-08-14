package mandrill

import (
	"github.com/mattbaird/gochimp"
)

type Client interface {
	SendTemplate(templateName string, message gochimp.Message) error
}

type client struct {
	client *gochimp.MandrillAPI
}

func NewClient(apiKey string) (Client, error) {
	c, err := gochimp.NewMandrill(apiKey)
	if err != nil {
		return nil, err
	}
	return &client{
		client: c,
	}, nil
}

func (c *client) SendTemplate(templateName string, message gochimp.Message) error {
	_, err := c.client.MessageSendTemplate(
		templateName,
		[]gochimp.Var{},
		message,
		true,
	)
	return err
}
