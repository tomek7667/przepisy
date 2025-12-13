package mails

import "fmt"

func (c *Client) Ping() error {
	_, err := c.http.R().Get("/ping")
	if err != nil {
		return fmt.Errorf("server unreachable at %s/ping: %w", c.BaseUrl, err)
	}
	return nil
}
