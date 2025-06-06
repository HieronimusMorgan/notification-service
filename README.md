# 📦 Notification Service with Golang, NATS, and Firebase

A scalable backend notification dispatch system written in **Go**, using **Gin**, **GORM**, **NATS**, and **Firebase Cloud Messaging (FCM)**. It supports asynchronous message delivery, persistent storage, and automatic retry mechanisms.

---

## 🚀 Features

- ✅ RESTful API with **Gin**
- ✅ Clean architecture with layered `internal/` structure
- ✅ **GORM** + PostgreSQL for persistent notification storage
- ✅ Asynchronous dispatch using **NATS** pub-sub
- ✅ Integration with **Firebase Admin SDK** (HTTP v1 API)
- ✅ Rich notification support:
  - `AndroidConfig`
  - `WebpushConfig`
  - `Data` payload
- ✅ Retry mechanism for failed or pending notifications

---

## 🧱 Project Structure

```bash
.
├── main.go
├── internal/
│   ├── controller/        # HTTP handlers
│   ├── models/            # GORM models
│   ├── repository/        # DB operations
│   ├── routes/            # HTTP route definitions
│   ├── services/          # Business logic + FCM
│   │   └── fcm.go
│   └── utils/
│       └── nats/          # NATS pub/sub logic
```

---

## 📄 API Example

### Endpoint
```http
POST /notify
```

### Request Body
```json
{
  "title": "System Alert",
  "body": "You have a new message",
  "target_token": "fcm_device_token",
  "platform": "android"
}
```

### Response
```json
{
  "message": "Notification enqueued"
}
```

---

## 🔧 Environment Variables

Create a `.env` file or set these values in your environment:

```env
POSTGRES_DSN=postgres://user:pass@localhost:5432/notification_service?sslmode=disable
NATS_URL=nats://localhost:4222
GOOGLE_APPLICATION_CREDENTIALS=./firebase-service-account.json
FIREBASE_PROJECT_ID=your-firebase-project-id
```

> ⚠️ Keep `firebase-service-account.json` private and excluded from version control.

---

## 🛠️ Running Locally

```bash
git clone https://github.com/your-username/notification-service.git
cd notification-service
go mod tidy
cp .env.example .env
# Make sure PostgreSQL and NATS are running
go run main.go
```

---

## 📬 Firebase Cloud Messaging Setup

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Navigate to **Project Settings > Service Accounts**
3. Click **"Generate new private key"** to download `firebase-service-account.json`
4. Save it to your project and reference it via `GOOGLE_APPLICATION_CREDENTIALS`

---

## 📦 Todo & Enhancements

- [ ] Multicast support
- [ ] Retry with exponential backoff
- [ ] Admin dashboard for notification history
- [ ] Monitoring endpoints (`/healthz`, `/metrics`)

---

## 📜 License

This project is MIT licensed. Feel free to use and adapt.
