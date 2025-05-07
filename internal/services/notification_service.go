package services

import (
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"google.golang.org/api/option"
	"log"
	"notification-service/internal/models"
	"notification-service/internal/repository"
	"time"
)

type NotificationService interface {
	SendNotificationAuthentication(data []byte) error
	SendNotificationAsset(data []byte) error
	SendNotification(notif *models.NotificationRequest) error
}

type notificationService struct {
	repo      repository.NotificationRepository
	ProjectID string
	FilePath  string
}

func NewNotificationService(repo repository.NotificationRepository, filePath, id string) NotificationService {
	return &notificationService{repo: repo, FilePath: filePath, ProjectID: id}
}

func (s *notificationService) SendNotificationAuthentication(data []byte) error {
	var notification models.NotificationResponse
	if err := json.Unmarshal(data, &notification); err != nil {
		return fmt.Errorf("unmarshal notification: %w", err)
	}

	if notification.EventType != "assign_user_resource" && notification.EventType != "remove_user_resource" {
		return fmt.Errorf("unsupported event type: %s", notification.EventType)
	}

	var tokenDetails models.TokenDetails
	payloadJSON, _ := json.Marshal(notification.Payload)
	if err := json.Unmarshal(payloadJSON, &tokenDetails); err != nil {
		return fmt.Errorf("unmarshal token details: %w", err)
	}

	payload := map[string]string{
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}

	now := time.Now()
	notifReq := &models.NotificationRequest{
		TargetToken: notification.TargetToken,
		Title:       notification.Title,
		Body:        notification.Body,
		Priority:    notification.Priority,
		Color:       notification.Color,
		ClickAction: notification.ClickAction,
		Payload:     payload,
	}

	notif := models.Notification{
		TargetToken:   notification.TargetToken,
		Title:         notification.Title,
		Body:          notification.Body,
		Platform:      notification.Platform,
		CreatedAt:     now,
		ServiceSource: notification.ServiceSource,
		EventType:     notification.EventType,
		ClickAction:   notification.ClickAction,
		Priority:      notification.Priority,
		Color:         notification.Color,
		Payload:       toJSONString(payload),
		Status:        "pending",
	}

	if err := s.repo.Save(&notif); err != nil {
		return fmt.Errorf("save notification: %w", err)
	}

	if err := s.SendNotification(notifReq); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}

	notif.Status = "sent"
	notif.SentAt = &now
	return s.repo.Update(&notif)
}

func (s *notificationService) SendNotificationAsset(data []byte) error {
	var notification models.NotificationResponse
	if err := json.Unmarshal(data, &notification); err != nil {
		return fmt.Errorf("unmarshal notification: %w", err)
	}

	if notification.EventType != "assign_user_resource" && notification.EventType != "remove_user_resource" {
		return fmt.Errorf("unsupported event type: %s", notification.EventType)
	}

	now := time.Now()
	notifReq := &models.NotificationRequest{
		TargetToken: notification.TargetToken,
		Title:       notification.Title,
		Body:        notification.Body,
		Priority:    notification.Priority,
		Color:       notification.Color,
		ClickAction: notification.ClickAction,
		Payload:     notification.Payload,
	}

	notif := models.Notification{
		TargetToken:   notification.TargetToken,
		Title:         notification.Title,
		Body:          notification.Body,
		ServiceSource: notification.ServiceSource,
		EventType:     notification.EventType,
		Priority:      notification.Priority,
		Color:         notification.Color,
		Payload:       toJSONString(notification.Payload),
		Status:        "pending",
	}

	if err := s.repo.Save(&notif); err != nil {
		return fmt.Errorf("save notification: %w", err)
	}

	if err := s.SendNotification(notifReq); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}

	notif.Status = "sent"
	notif.SentAt = &now
	return s.repo.Update(&notif)
}

func (s *notificationService) SendNotification(request *models.NotificationRequest) error {
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: s.ProjectID}, option.WithCredentialsFile(s.FilePath))
	if err != nil {
		return fmt.Errorf("firebase init: %w", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("firebase client: %w", err)
	}

	msg := &messaging.Message{
		Token: request.TargetToken,
		Data:  request.Payload,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Title:       request.Title,
				Body:        request.Body,
				Color:       request.Color,
				ClickAction: request.ClickAction,
				Icon:        "default",
				Sound:       "default",
			},
		},
		Notification: &messaging.Notification{
			Title: request.Title,
			Body:  request.Body,
		},
	}

	log.Printf("📤 Sending FCM with payload: %+v", request.Payload)

	resp, err := client.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("FCM send: %w", err)
	}

	log.Printf("✅ FCM sent: %s", resp)
	return nil
}

func toJSONString(data map[string]string) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("⚠️ Failed to marshal payload: %v", err)
		return "{}"
	}
	return string(jsonBytes)
}
