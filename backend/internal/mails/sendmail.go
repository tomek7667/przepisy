package mails

import (
	"fmt"
	"time"
)

type AgentType = string

const (
	AgentGmail  AgentType = "gmail"
	AgentiCloud AgentType = "icloud"
)

func GetAgents() []AgentType {
	return []AgentType{
		AgentGmail,
		AgentiCloud,
	}
}

type Priority = string

const (
	PriorityHigh   Priority = "high"
	PriorityNormal Priority = "normal"
	PriorityLow    Priority = "low"
)

type Attachment struct {
	Filename    string `json:"filename"`
	Content     string `json:"content"`
	ContentType string `json:"contentType"`
}

type Options struct {
	From        string       `json:"from"`
	To          string       `json:"to"`
	Subject     string       `json:"subject"`
	HTML        string       `json:"html"`
	Priority    Priority     `json:"priority"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Mail struct {
	ID        int64      `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	SentAt    *time.Time `json:"sent_at"`
	Subject   string     `json:"subject"`
	Contents  string     `json:"contents"`
	Target    string     `json:"target"`
	Sender    string     `json:"sender"`
	Details   *string    `json:"details"`
}

type sendMailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    Mail   `json:"data"`
}

func (c *Client) SendMail(agentID int64, mailOptions Options) (*Mail, error) {
	body := map[string]any{
		"token":   c.Token,
		"agentId": agentID,
		"options": mailOptions,
	}
	var res sendMailResponse
	resp, err := c.http.R().
		SetBody(body).
		SetResult(&res).
		Post("/api/mails")
	if err != nil {
		return nil, fmt.Errorf("failed to send mail: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to send mail: %s", resp.Status())
	}
	if !res.Success {
		return nil, fmt.Errorf("failed to send mail: %s", res.Message)
	}
	return &res.Data, nil
}
