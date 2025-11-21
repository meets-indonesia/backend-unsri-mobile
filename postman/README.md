# Postman Collection

Postman collection untuk testing UNSRI Backend API.

## 📦 Files

- `UNSRI_Backend_API.postman_collection.json` - API Collection
- `UNSRI_Backend_Environment.postman_environment.json` - Environment variables

## 🚀 Quick Start

### 1. Import Collection

1. Buka Postman
2. Klik **Import** button
3. Pilih file `UNSRI_Backend_API.postman_collection.json`
4. Klik **Import**

### 2. Import Environment

1. Klik **Environments** di sidebar kiri
2. Klik **Import**
3. Pilih file `UNSRI_Backend_Environment.postman_environment.json`
4. Klik **Import**

### 3. Set Environment Variables

1. Pilih environment **UNSRI Backend - Local**
2. Set `base_url` sesuai environment Anda:
   - Local: `http://localhost:8080`
   - Production: `https://api.unsri.ac.id`

### 4. Get Access Token

1. Buka request **Authentication > Login**
2. Update email dan password di body
3. Klik **Send**
4. Copy `access_token` dari response
5. Set variable `access_token` di environment

### 5. Test API

Sekarang Anda bisa test semua API endpoints. Token akan otomatis digunakan untuk authenticated requests.

## 📋 Collection Structure

- **Authentication**
  - Register
  - Login

- **Users**
  - Get Profile
  - Update Profile

- **Attendance**
  - Scan QR Code
  - Get Attendance History

- **QR Code**
  - Generate Class QR
  - Generate Access QR

- **Location**
  - Tap In
  - Tap Out

- **Search**
  - Search Users
  - Global Search

- **Reports**
  - Attendance Report
  - Academic Report

## 🔧 Environment Variables

- `base_url` - Base URL untuk API (default: http://localhost:8080)
- `access_token` - JWT access token (akan di-set setelah login)
- `refresh_token` - JWT refresh token
- `user_id` - Current user ID

## 💡 Tips

1. **Auto-save Token**: Buat test script di Login request untuk auto-save token:
```javascript
if (pm.response.code === 200) {
    var jsonData = pm.response.json();
    pm.environment.set("access_token", jsonData.data.access_token);
}
```

2. **Pre-request Script**: Untuk auto-include token di semua requests, tambahkan di collection level:
```javascript
pm.request.headers.add({
    key: "Authorization",
    value: "Bearer " + pm.environment.get("access_token")
});
```

## 📚 Documentation

Lihat [API Documentation](../docs/API.md) untuk dokumentasi lengkap semua endpoints.

