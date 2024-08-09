package addressvalidation

import (
	addressvalidation "cloud.google.com/go/maps/addressvalidation/apiv1"
	"cloud.google.com/go/maps/addressvalidation/apiv1/addressvalidationpb"
	"context"
	"google.golang.org/genproto/googleapis/type/postaladdress"
)

type Client struct {
	client *addressvalidation.Client
	token  string
}

func NewClient(ctx context.Context, token string) (*Client, error) {
	client, err := addressvalidation.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
		token:  token,
	}, nil
}

func (c *Client) NormalizeAddress(ctx context.Context, address string) (string, error) {
	request := &addressvalidationpb.ValidateAddressRequest{
		Address: &postaladdress.PostalAddress{
			AddressLines: []string{address},
		},
	}
	response, err := c.client.ValidateAddress(ctx, request)
	if err != nil {
		return "", err
	}
	return response.GetResult().Address.FormattedAddress, nil
}
