# API Documentation

Dokumentasi lengkap untuk semua API endpoints.

## Base URL

- **Local**: `http://localhost:8080`
- **Production**: `https://api.unsri.ac.id`

## Authentication

Kebanyakan endpoint memerlukan authentication menggunakan JWT token.

### Get Token

```bash
POST /api/v1/auth/login
```

Response:
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900
  }
}
```

### Using Token

Include token di header:
```
Authorization: Bearer <access_token>
```

## Endpoints

### Authentication

#### Register
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "student@unsri.ac.id",
  "password": "password123",
  "role": "mahasiswa",
  "name": "John Doe",
  "nim": "1234567890"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "student@unsri.ac.id",
  "password": "password123"
}
```

### Users

#### Get Profile
```http
GET /api/v1/users/profile
Authorization: Bearer <token>
```

#### Update Profile
```http
PUT /api/v1/users/profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "John Doe Updated",
  "phone": "081234567890"
}
```

### Attendance

#### Scan QR Code
```http
POST /api/v1/attendance/scan
Authorization: Bearer <token>
Content-Type: application/json

{
  "qr_code": "<qr_code_data>"
}
```

#### Get Attendance History
```http
GET /api/v1/attendance/history?page=1&per_page=20
Authorization: Bearer <token>
```

### QR Code

#### Generate Class QR
```http
POST /api/v1/qr/class/generate
Authorization: Bearer <token>
Content-Type: application/json

{
  "schedule_id": "<schedule_id>",
  "duration": 15
}
```

#### Generate Access QR
```http
POST /api/v1/qr/access/generate
Authorization: Bearer <token>
Content-Type: application/json

{
  "duration": 525600
}
```

### Location

#### Tap In
```http
POST /api/v1/location/tap-in
Authorization: Bearer <token>
Content-Type: application/json

{
  "latitude": -2.9914,
  "longitude": 104.7565
}
```

#### Tap Out
```http
POST /api/v1/location/tap-out
Authorization: Bearer <token>
Content-Type: application/json

{
  "latitude": -2.9914,
  "longitude": 104.7565
}
```

### Search

#### Search
```http
GET /api/v1/search?q=john&type=users&page=1&per_page=20
Authorization: Bearer <token>
```

#### Global Search
```http
GET /api/v1/search/global?q=algorithm&types=users,courses&limit=10
Authorization: Bearer <token>
```

### Reports

#### Attendance Report
```http
GET /api/v1/reports/attendance?start_date=2024-01-01&end_date=2024-01-31&summary=true
Authorization: Bearer <token>
```

#### Academic Report
```http
GET /api/v1/reports/academic?student_id=<student_id>&semester=2024-1
Authorization: Bearer <token>
```

## Error Responses

Semua error mengikuti format:

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

### Error Codes

- `UNAUTHORIZED` - Token tidak valid atau expired
- `FORBIDDEN` - Tidak memiliki permission
- `NOT_FOUND` - Resource tidak ditemukan
- `BAD_REQUEST` - Request tidak valid
- `VALIDATION_FAILED` - Validasi gagal
- `CONFLICT` - Resource conflict
- `INTERNAL_ERROR` - Server error

## Rate Limiting

API memiliki rate limiting:
- **Default**: 100 requests per minute per IP
- **Authenticated**: 1000 requests per minute per user

## Pagination

Endpoints yang support pagination menggunakan query parameters:
- `page` - Page number (default: 1)
- `per_page` - Items per page (default: 20, max: 100)

Response format:
```json
{
  "success": true,
  "data": [...],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

## Status Codes

- `200 OK` - Request berhasil
- `201 Created` - Resource berhasil dibuat
- `400 Bad Request` - Request tidak valid
- `401 Unauthorized` - Tidak terautentikasi
- `403 Forbidden` - Tidak memiliki permission
- `404 Not Found` - Resource tidak ditemukan
- `409 Conflict` - Resource conflict
- `500 Internal Server Error` - Server error

