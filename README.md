# Guitar Learning Backend (Go)

Backend for the Guitar Learning platform, rewritten from Spring Boot to **Golang** to reduce infrastructure costs.
All API endpoints are kept identical to the original Spring Boot project.

## Directory Structure

```
cmd/server/          → Entry point (main.go)
internal/
  common/            → Middleware, response handling, errors, context
  config/            → Environment variable loading
  auth/              → Authentication (register, login, google/firebase, refresh, logout)
  user/              → User profile
  content/           → Lesson content blocks
  music/             → Chords, strumming patterns
  learning/          → Courses, modules, lessons, learning methods
  progress/          → Learning progress
  song/              → Song library
  assessment/        → Quizzes, attempts
  onboarding/        → Onboarding flow
db/migration/        → Migration files
```

## Requirements

- Go 1.22+
- Environment variables according to `.env.example`

## Installation & Running

```bash
# 1. Enter the project directory
cd backend-go

# 2. Copy the example env file and edit it
cp .env.example .env

# 3. Download dependencies
go mod tidy

# 4. Run the server
go run ./cmd/server
```

The server listens at `http://localhost:5501` by default.

## API Endpoints (unchanged from original)

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/google` | Sign in with Firebase/Google |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| POST | `/api/v1/auth/logout` | Logout |
| GET | `/api/v1/users/me` | Get profile (requires JWT) |
| PUT | `/api/v1/users/me` | Update profile (requires JWT) |
| GET | `/api/v1/users/me/progress` | Get learning progress (requires JWT) |
| POST | `/api/v1/lessons/{id}/complete` | Mark lesson complete (requires JWT) |
| GET | `/api/v1/lessons/{lessonId}/blocks` | Get lesson content blocks |
| GET | `/api/v1/chords` | List chords |
| GET | `/api/v1/chords/{id}` | Get chord detail |
| GET | `/api/v1/chords/{id}/positions` | Get chord finger positions |
| GET | `/api/v1/strumming-patterns` | List strumming patterns |
| POST | `/api/v1/strumming-patterns` | Create strumming pattern (admin) |
| GET | `/api/v1/strumming-patterns/{id}` | Get strumming pattern detail |
| GET | `/api/v1/strumming-patterns/{id}/instructions` | Get step-by-step instructions |
| GET | `/api/v1/courses` | List courses |
| GET | `/api/v1/courses/{id}` | Get course detail |
| GET | `/api/v1/courses/{id}/modules` | Get course modules |
| GET | `/api/v1/modules/{id}` | Get module detail |
| GET | `/api/v1/modules/{moduleId}/lessons` | Get lessons in module |
| GET | `/api/v1/lessons/{id}` | Get lesson detail |
| GET | `/api/v1/learning-methods` | List learning methods |
| GET | `/api/v1/learning-methods/{id}` | Get learning method detail |
| GET | `/api/v1/songs` | List songs |
| POST | `/api/v1/songs` | Create song (admin) |
| GET | `/api/v1/songs/{id}` | Get song detail |
| GET | `/api/v1/lessons/{lessonId}/quiz` | Get quiz for lesson |
| POST | `/api/v1/quizzes/{quizId}/attempts` | Start quiz attempt (requires JWT) |
| POST | `/api/v1/quiz-attempts/{attemptId}/submit` | Submit quiz attempt (requires JWT) |
| GET | `/api/v1/quiz-attempts/{attemptId}` | Get attempt result (requires JWT) |
| GET | `/api/v1/onboarding/status` | Get onboarding status (requires JWT) |
| POST | `/api/v1/onboarding/complete` | Complete onboarding (requires JWT) |

## CORS

Configured in `internal/config/config.go` and the `CORS_ALLOWED_ORIGINS` environment variable.
Default allowed origins: `http://localhost:5500` and `http://127.0.0.1:5500`.

## Docker

```bash
docker build -t hguitar-backend-go .
docker run -p 5501:5501 hguitar-backend-go
```

## Deploy on Render

1. Push this repository to GitHub.
2. Connect the repository to Render.
3. Create a new **Web Service** using Docker.
4. Set the required environment variables in the Render Dashboard.
5. Deploy the service.

## Important Notes

- Password hashing currently uses simple SHA-256. For production, replace it with **bcrypt** or **argon2** in `internal/auth/infrastructure/security/password.go`.
- Firebase ID token verification currently only decodes the payload (dev/local). For production, use `firebase-admin-sdk-go` to verify the signature properly.
