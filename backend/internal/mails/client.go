package mails

import "resty.dev/v3"

type Client struct {
	BaseUrl string
	Token   string
	http    *resty.Client
}

func New(baseUrl, token string) (*Client, error) {
	httpClient := resty.New().
		SetBaseURL(baseUrl).
		SetHeader("Authorization", "Mailer "+token)

	c := &Client{
		BaseUrl: baseUrl,
		Token:   token,
		http:    httpClient,
	}
	if err := c.Ping(); err != nil {
		return nil, err
	}
	return c, nil
}
