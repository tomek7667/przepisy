package mails

import "fmt"

type Agent struct {
	ID          int64  `json:"id"`
	CreatedAt   string `json:"created_at"`
	Email       string `json:"email"`
	Type        string `json:"type"`
	Credentials string `json:"credentials"`
}

type listAgentsResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		Agents []Agent `json:"agents"`
	} `json:"data"`
}

func (c *Client) ListAgents() ([]Agent, error) {
	var res listAgentsResponse
	resp, err := c.http.R().
		SetResult(&res).
		Get("/api/agents")
	if err != nil {
		return nil, fmt.Errorf("failed to request agents: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to list agents: %s", resp.Status())
	}
	if !res.Success {
		return nil, fmt.Errorf("failed to list agents: %s", res.Message)
	}
	return res.Data.Agents, nil
}
