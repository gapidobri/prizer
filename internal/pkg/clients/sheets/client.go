package sheets

import (
	"context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Client interface {
	AppendRow(sheetId string, tabName string, values []any) error
}

type client struct {
	client *sheets.Service
}

func NewClient(ctx context.Context, serviceAccountKeyPath string) (Client, error) {
	service, err := sheets.NewService(ctx, option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		return nil, err
	}
	return &client{
		client: service,
	}, nil
}

func (c *client) AppendRow(sheetId string, tabName string, values []any) error {
	_, err := c.client.Spreadsheets.Values.
		Append(sheetId, tabName, &sheets.ValueRange{
			MajorDimension: "ROWS",
			Values:         [][]any{values},
		}).
		ValueInputOption("RAW").
		Do()
	return err
}
