package models

import (
	"time"
)

type Notification struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	TargetToken   string     `gorm:"not null;index" json:"target_token"`
	Title         string     `gorm:"not null" json:"title"`
	Body          string     `gorm:"not null" json:"body"`
	Platform      string     `gorm:"not null;index" json:"platform"`        // "android", "web"
	Priority      string     `gorm:"default:'high'" json:"priority"`        // "high", "normal"
	Status        string     `gorm:"default:'pending';index" json:"status"` // "pending", "sent", "failed"
	ServiceSource string     `gorm:"not null;index" json:"service_source"`  // e.g., "auth"
	EventType     string     `gorm:"not null;index" json:"event_type"`      // e.g., "asset_updated"
	Payload       string     `gorm:"type:text" json:"payload"`              // raw JSON string
	Color         string     `gorm:"default:'#000000'" json:"color"`
	ClickAction   string     `gorm:"default:'OPEN_APP'" json:"click_action"`
	Icon          string     `gorm:"default:'default'" json:"icon"`
	Sound         string     `gorm:"default:'default'" json:"sound"`
	RetryCount    int        `gorm:"default:0" json:"retry_count"`
	LastError     *string    `gorm:"type:text" json:"last_error,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	SentAt        *time.Time `json:"sent_at,omitempty"`
}

type NotificationResponse struct {
	TargetToken   string            `json:"target_token"`
	Title         string            `json:"title"`
	Body          string            `json:"body"`
	Platform      string            `json:"platform"`
	ServiceSource string            `json:"service_source"`
	EventType     string            `json:"event_type"`
	Payload       map[string]string `json:"payload"`
	Color         string            `json:"color"`
	Priority      string            `json:"priority"`
	ClickAction   string            `json:"click_action"`
}

type NotificationRequest struct {
	TargetToken string            `json:"target_token"`
	Title       string            `json:"title"`
	Body        string            `json:"body"`
	Payload     map[string]string `json:"payload"`
	Color       string            `json:"color"`
	Priority    string            `json:"priority"`
	ClickAction string            `json:"click_action"`
}
