package services

import (
	"bytes"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"google.golang.org/api/option"
	"html/template"
	"log"
	"net/smtp"
	"notification-service/internal/models"
	"notification-service/internal/repository"
	"time"
)

type NotificationService interface {
	SendNotificationAuthentication(data []byte) error
	SendNotificationEmail(data []byte) error
	SendNotificationAsset(data []byte) error
	SendNotification(notif *models.NotificationRequest) error
}

type notificationService struct {
	repo      repository.NotificationRepository
	ProjectID string
	FilePath  string
	SMTPHost  string
	SMTPPort  string
	Email     string
	Password  string
}

func NewNotificationService(repo repository.NotificationRepository, filePath, id, smtpHost, smtpPort, email, password string) NotificationService {
	return &notificationService{repo: repo, FilePath: filePath, ProjectID: id, SMTPHost: smtpHost, SMTPPort: smtpPort, Email: email, Password: password}
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

	log.Printf("payloadJSON message : %s", payloadJSON)
	log.Printf("tokenDetails message : %s", tokenDetails)

	payload := map[string]string{
		"type":          "system",
		"access_token":  tokenDetails.AccessToken,
		"refresh_token": tokenDetails.RefreshToken,
	}
	log.Printf("payload message : %s", payload)

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

func (s *notificationService) SendNotificationEmail(data []byte) error {
	var email models.Email
	if err := json.Unmarshal(data, &email); err != nil {
		return fmt.Errorf("unmarshal notification: %w", err)
	}

	to := []string{email.To}

	subject := "Subject: Reset Your Password\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	// Define the HTML template
	htmlTemplate := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
		  <meta charset="UTF-8">
		  <title>Forgot Your Password?</title>
		  <style>
			body {
			  font-family: Arial, sans-serif;
			  background-color: #f4f4f4;
			  margin: 0;
			  padding: 20px;
			}
			.container {
			  max-width: 600px;
			  margin: auto;
			  background-color: #ffffff;
			  padding: 30px;
			  border-radius: 8px;
			  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
			}
			.button {
			  display: inline-block;
			  margin-top: 20px;
			  padding: 12px 24px;
			  background-color: #1e88e5;
			  color: #ffffff;
			  text-decoration: none;
			  border-radius: 5px;
			  font-weight: bold;
			}
			.footer {
			  color: #999;
			  font-size: 12px;
			  margin-top: 30px;
			}
		  </style>
		</head>
		<body>
		  <div class="container">
			<h2 style="color: #333;">Hi {{.FullName}},</h2>
			<p style="color: #555;">We received a request to reset the password for your account.</p>
			<p style="color: #555;">To continue, please click the button below. You’ll be redirected to our app to set a new password:</p>
			<a href="{{.URL}}" class="button">Reset Your Password</a>
			<p class="footer">If you didn't request this, you can safely ignore this email. Your password will remain unchanged.</p>
			  </div>
			</body>
			</html>
			`

	// Use html/template to inject dynamic values
	tmpl, err := template.New("email").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, struct {
		FullName string
		URL      string
	}{
		FullName: email.FullName,
		URL:      email.URL,
	}); err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}

	// Final email message
	message := []byte(subject + mime + body.String())

	// Auth
	auth := smtp.PlainAuth("", s.Email, s.Password, s.SMTPHost)

	// Send
	if err := smtp.SendMail(s.SMTPHost+":"+s.SMTPPort, auth, s.Email, to, message); err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	log.Println("✅ Email sent successfully to", email.To)
	return nil
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
