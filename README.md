# Guitar Learning Backend (Go)

Port của project backend gốc (Spring Boot) sang **Golang** để giảm chi phí vận hành.
Tất cả endpoint API được giữ nguyên như project gốc.

## Cấu trúc thư mục

```
cmd/server/          -> Entry point (main.go)
internal/
  common/            -> Middleware, response, errors, ctx, db helper
  config/            -> Load biến môi trường
  auth/              -> Xác thực (register, login, google/firebase, refresh, logout)
  user/              -> Hồ sơ người dùng
  content/           -> Nội dung bài học (blocks)
  music/             -> Hợp âm, strumming patterns
  learning/          -> Khóa học, module, bài học, phương pháp học
  progress/          -> Tiến độ học tập
  song/              -> Thư viện bài hát
  assessment/        -> Trắc nghiệm, bài thi
  onboarding/        -> Hướng dẫn bắt đầu
db/migration/        -> Schema SQL Server (Flyway)
```

## Yêu cầu

- Go 1.22+
- SQL Server (máy local hoặc Azure)
- Biến môi trường theo file `.env.example`

## Cài đặt & chạy

```bash
# 1. Clone và vào thư mục
cd backend-go

# 2. Tạo file .env từ ví dụ (chỉnh DB_URL, JWT_SECRET...)
copy .env.example .env

# 3. Tải dependencies
go mod tidy

# 4. Chạy migration V1__init.sql lên SQL Server trước (hoặc dùng Flyway)

# 5. Chạy server
go run ./cmd/server
```

Server lắng nghe tại `http://localhost:5501` (mặc định).

## Danh sách endpoint API (giữ nguyên)

| Phương thức | Endpoint | Mô tả |
|---|---|---|
| POST | `/api/v1/auth/register` | Đăng ký |
| POST | `/api/v1/auth/login` | Đăng nhập |
| POST | `/api/v1/auth/google` | Đăng nhập Firebase/Google |
| POST | `/api/v1/auth/refresh` | Làm mới access token |
| POST | `/api/v1/auth/logout` | Đăng xuất |
| GET | `/api/v1/users/me` | Xem hồ sơ (cần JWT) |
| PUT | `/api/v1/users/me` | Cập nhật hồ sơ (cần JWT) |
| GET | `/api/v1/users/me/progress` | Tiến độ học (cần JWT) |
| POST | `/api/v1/lessons/{id}/complete` | Đánh dấu bài học hoàn thành (cần JWT) |
| GET | `/api/v1/lessons/{lessonId}/blocks` | Lấy khối nội dung bài học |
| GET | `/api/v1/chords` | Danh sách hợp âm |
| GET | `/api/v1/chords/{id}` | Chi tiết hợp âm |
| GET | `/api/v1/chords/{id}/positions` | Vị trí ngón tay |
| GET | `/api/v1/strumming-patterns` | Danh sách mẫu đánh |
| POST | `/api/v1/strumming-patterns` | Tạo mẫu đánh (admin) |
| GET | `/api/v1/strumming-patterns/{id}` | Chi tiết mẫu đánh |
| GET | `/api/v1/strumming-patterns/{id}/instructions` | Hướng dẫn từng bước |
| GET | `/api/v1/courses` | Danh sách khóa học |
| GET | `/api/v1/courses/{id}` | Chi tiết khóa học |
| GET | `/api/v1/courses/{id}/modules` | Module của khóa học |
| GET | `/api/v1/modules/{id}` | Chi tiết module |
| GET | `/api/v1/modules/{moduleId}/lessons` | Bài học trong module |
| GET | `/api/v1/lessons/{id}` | Chi tiết bài học |
| GET | `/api/v1/learning-methods` | Phương pháp học |
| GET | `/api/v1/learning-methods/{id}` | Chi tiết phương pháp |
| GET | `/api/v1/songs` | Danh sách bài hát |
| POST | `/api/v1/songs` | Tạo bài hát (admin) |
| GET | `/api/v1/songs/{id}` | Chi tiết bài hát |
| GET | `/api/v1/lessons/{lessonId}/quiz` | Trắc nghiệm bài học |
| POST | `/api/v1/quizzes/{quizId}/attempts` | Bắt đầu làm bài (cần JWT) |
| POST | `/api/v1/quiz-attempts/{attemptId}/submit` | Nộp bài (cần JWT) |
| GET | `/api/v1/quiz-attempts/{attemptId}` | Xem kết quả (cần JWT) |
| GET | `/api/v1/onboarding/status` | Trạng thái mở đầu (cần JWT) |
| POST | `/api/v1/onboarding/complete` | Hoàn thành mở đầu (cần JWT) |

## CORS

Cấu hình trong `internal/config/config.go` và biến môi trường `CORS_ALLOWED_ORIGINS`. Mặc định cho phép `http://localhost:5500` và `http://127.0.0.1:5500`.

## Docker

```bash
docker build -t hguitar-backend-go .
docker run -p 5501:5501 --env-file .env hguitar-backend-go
```

## Lưu ý

- Mã hóa mật khẩu hiện dùng SHA-256 đơn giản. Trong production hãy thay bằng **bcrypt** hoặc **argon2**.
- Xác thực Firebase ID token hiện chỉ decode payload (dev/local). Production hãy dùng `firebase-admin-sdk-go` để verify ký thực tế.
- Các câu lệnh SQL dùng placeholder `@p1, @p2...` tương thích driver `go-mssqldb`.
