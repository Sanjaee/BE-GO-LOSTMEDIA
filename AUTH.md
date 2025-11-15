# Authentication Documentation

## Overview

Sistem autentikasi mencakup:
- **Register**: Pendaftaran user baru dengan email verification
- **Login**: Login dengan email dan password
- **Verify Email**: Verifikasi email setelah registrasi
- **Forgot Password**: Request reset password
- **Verify Reset Password**: Verifikasi token reset password
- **Reset Password**: Reset password dengan token

## API Endpoints

### 1. Register

**Endpoint**: `POST /api/v1/auth/register`

**Request Body**:
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response** (201):
```json
{
  "data": {
    "user": {
      "userId": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "isEmailVerified": false,
      "createdAt": "2025-11-15T09:00:00Z"
    },
    "token": "email_verification_token",
    "refreshToken": ""
  }
}
```

**Notes**:
- User akan menerima email verification
- Token yang dikembalikan adalah verification token, bukan JWT
- JWT token akan diberikan setelah email terverifikasi atau saat login

### 2. Login

**Endpoint**: `POST /api/v1/auth/login`

**Request Body**:
```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response** (200):
```json
{
  "data": {
    "user": {
      "userId": "uuid",
      "username": "john_doe",
      "email": "john@example.com",
      "isEmailVerified": true
    },
    "token": "jwt_access_token",
    "refreshToken": "jwt_refresh_token"
  }
}
```

**Notes**:
- Token expired dalam 24 jam
- Refresh token expired dalam 7 hari
- Gunakan Bearer token di Authorization header untuk protected endpoints

### 3. Verify Email

**Endpoint**: `POST /api/v1/auth/verify-email`

**Request Body**:
```json
{
  "token": "email_verification_token"
}
```

**Response** (200):
```json
{
  "data": {
    "message": "Email verified successfully"
  }
}
```

**Notes**:
- Token dari email verification link
- Token expired dalam 24 jam

### 4. Forgot Password

**Endpoint**: `POST /api/v1/auth/forgot-password`

**Request Body**:
```json
{
  "email": "john@example.com"
}
```

**Response** (200):
```json
{
  "data": {
    "message": "If the email exists, a password reset link has been sent"
  }
}
```

**Notes**:
- Always returns success for security (doesn't reveal if email exists)
- User akan menerima email dengan password reset link
- Token expired dalam 1 jam

### 5. Verify Reset Password Token

**Endpoint**: `POST /api/v1/auth/verify-reset-password`

**Request Body**:
```json
{
  "token": "password_reset_token"
}
```

**Response** (200):
```json
{
  "data": {
    "message": "Reset token is valid"
  }
}
```

**Notes**:
- Gunakan sebelum menampilkan form reset password
- Verifikasi token masih valid

### 6. Reset Password

**Endpoint**: `POST /api/v1/auth/reset-password`

**Request Body**:
```json
{
  "token": "password_reset_token",
  "newPassword": "newpassword123"
}
```

**Response** (200):
```json
{
  "data": {
    "message": "Password reset successfully"
  }
}
```

**Notes**:
- Password minimal 8 karakter
- Token akan dihapus setelah reset berhasil
- Token expired dalam 1 jam

## Flow Diagram

### Registration Flow
```
User registers
    ↓
Create user with email verification token
    ↓
Send verification email
    ↓
User clicks verification link
    ↓
Verify email endpoint
    ↓
Email verified, user can login
```

### Password Reset Flow
```
User requests password reset
    ↓
Generate reset token (1 hour expiry)
    ↓
Send reset email
    ↓
User clicks reset link
    ↓
Frontend verifies token
    ↓
User enters new password
    ↓
Reset password endpoint
    ↓
Password updated, token cleared
```

## Security Features

1. **Password Hashing**: Menggunakan bcrypt dengan cost 12
2. **JWT Tokens**: Secure token-based authentication
3. **Token Expiration**: Email verification (24h), password reset (1h)
4. **Email Verification**: Required untuk aktivasi akun
5. **Rate Limiting**: (dapat ditambahkan di middleware)
6. **CORS Protection**: Configured untuk allowed origins

## Error Responses

Semua error mengikuti format:
```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message"
  }
}
```

### Common Error Codes

- `VALIDATION_ERROR` (400): Invalid input
- `UNAUTHORIZED` (401): Invalid credentials or token
- `CONFLICT` (409): Email/username already exists
- `INTERNAL_ERROR` (500): Server error

## Usage Example

### Frontend Integration

```javascript
// Register
const register = async (username, email, password) => {
  const response = await fetch('http://localhost:5000/api/v1/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username, email, password })
  });
  return response.json();
};

// Login
const login = async (email, password) => {
  const response = await fetch('http://localhost:5000/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  const data = await response.json();
  // Store tokens
  localStorage.setItem('token', data.data.token);
  localStorage.setItem('refreshToken', data.data.refreshToken);
  return data;
};

// Protected API call
const fetchProtectedData = async () => {
  const token = localStorage.getItem('token');
  const response = await fetch('http://localhost:5000/api/v1/protected', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
};
```

## Testing

### Using cURL

```bash
# Register
curl -X POST http://localhost:5000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:5000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Verify Email
curl -X POST http://localhost:5000/api/v1/auth/verify-email \
  -H "Content-Type: application/json" \
  -d '{"token":"verification_token"}'

# Forgot Password
curl -X POST http://localhost:5000/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# Reset Password
curl -X POST http://localhost:5000/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token":"reset_token","newPassword":"newpassword123"}'
```

## Notes

1. **Email Configuration**: Jika SMTP tidak dikonfigurasi, email akan di-log ke console (development mode)
2. **Database Migration**: Pastikan migration sudah dijalankan sebelum menggunakan auth
3. **Environment Variables**: Konfigurasi JWT_SECRET dan database harus di-set
4. **Frontend URL**: Pastikan FRONTEND_URL di-set untuk email links

