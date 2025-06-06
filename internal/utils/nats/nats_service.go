package nats

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	"notification-service/internal/services"
	"sync"
)

func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Failed to marshal json: %v", err)
	}
	return data
}

type Service interface {
	Publish(subject string, data interface{}) error
	Subscribe()
	RetryPending(subject string)
}

type natsService struct {
	natsURL             string
	nc                  *nats.Conn
	notificationService services.NotificationService
	once                sync.Once
}

func NewNatsService(natsURL string, service services.NotificationService) Service {
	svc := &natsService{natsURL: natsURL, notificationService: service}
	svc.connect()
	svc.Subscribe()
	return svc
}

func (s *natsService) connect() {
	s.once.Do(func() {
		nc, err := nats.Connect(s.natsURL)
		if err != nil {
			log.Fatalf("NATS connection failed: %v", err)
		}
		s.nc = nc
	})
}

func (s *natsService) Publish(subject string, data interface{}) error {
	msg, _ := json.Marshal(data)
	return s.nc.Publish(subject, msg)
}

func (s *natsService) Subscribe() {
	subjects := []string{"authentication", "forgot_password", "asset"}

	for _, subject := range subjects {
		sub := subject
		_, err := s.nc.Subscribe(sub, func(m *nats.Msg) {
			log.Printf("Received message on %s: %s", sub, string(m.Data))
			switch sub {
			case "authentication":
				if err := s.notificationService.SendNotificationAuthentication(m.Data); err != nil {
					log.Printf("Error processing 'authentication': %v", err)
				} else {
					log.Printf("Processed 'authentication' successfully")
				}
			case "forgot_password":
				// Handle 'forgot_password' message
				if err := s.notificationService.SendNotificationEmail(m.Data); err != nil {
					log.Printf("Error processing 'authentication': %v", err)
				} else {
					log.Printf("Processed 'authentication' successfully")
				}
			case "asset":
				// Handle 'asset' message
				log.Println("Handling asset message")
			}
		})
		if err != nil {
			log.Fatalf("Failed to subscribe to %s: %v", sub, err)
		}
	}
	select {} // Keep the subscriber running indefinitely
}

func (s *natsService) RetryPending(subject string) {

}
