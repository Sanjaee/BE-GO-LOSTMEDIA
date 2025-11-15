# API Documentation

## Base Information

- **Base URL (Development)**: `http://localhost:8080/api/v1`
- **Base URL (Production)**: `https://api.lostmediago.com/api/v1`
- **Content-Type**: `application/json`
- **Authentication**: Bearer Token (JWT)

## Authentication

### Register
```http
POST /api/v1/auth/register
```

**Request Body:**
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "data": {
    "user": {
      "userId": "uuid",
      "username": "john_doe",
      "email": "john@example.com"
    },
    "token": "jwt_token",
    "refreshToken": "refresh_token"
  }
}
```

### Login
```http
POST /api/v1/auth/login
```

**Request Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

### Google OAuth Login
```http
POST /api/v1/auth/login/google
```

**Request Body:**
```json
{
  "idToken": "google_id_token"
}
```

### Refresh Token
```http
POST /api/v1/auth/refresh
```

**Request Body:**
```json
{
  "refreshToken": "refresh_token"
}
```

### Get Current User
```http
GET /api/v1/auth/me
```

**Headers:**
```
Authorization: Bearer {token}
```

## Users

### List Users
```http
GET /api/v1/users?page=1&limit=20&search=john
```

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `search` (optional): Search by username or email

### Get User by ID
```http
GET /api/v1/users/{userId}
```

### Update User Profile
```http
PUT /api/v1/users/{userId}
```

**Request Body:**
```json
{
  "username": "new_username",
  "bio": "Updated bio"
}
```

### Update Avatar
```http
PUT /api/v1/users/{userId}/avatar
```

**Content-Type:** `multipart/form-data`

**Form Data:**
- `avatar`: Image file

### Follow User
```http
POST /api/v1/users/{userId}/follow
```

### Unfollow User
```http
DELETE /api/v1/users/{userId}/follow
```

### Get User Followers
```http
GET /api/v1/users/{userId}/followers?page=1&limit=20
```

### Get User Following
```http
GET /api/v1/users/{userId}/following?page=1&limit=20
```

### Get User Posts
```http
GET /api/v1/users/{userId}/posts?page=1&limit=20
```

## Posts

### List Posts (Feed)
```http
GET /api/v1/posts?page=1&limit=20&category=all&sort=newest
```

**Query Parameters:**
- `page` (optional): Page number
- `limit` (optional): Items per page
- `category` (optional): Filter by category
- `sort` (optional): Sort order (`newest`, `popular`, `trending`)

### Get Post by ID
```http
GET /api/v1/posts/{postId}
```

### Create Post
```http
POST /api/v1/posts
```

**Request Body:**
```json
{
  "title": "Post Title",
  "description": "Post description",
  "content": "Post content",
  "category": "media",
  "mediaUrl": "https://cloudinary.com/image.jpg",
  "blurred": true,
  "scheduledAt": "2025-11-16T10:00:00Z",
  "sections": [
    {
      "type": "text",
      "content": "Section content",
      "order": 1
    },
    {
      "type": "image",
      "src": "https://cloudinary.com/image.jpg",
      "order": 2
    }
  ]
}
```

### Update Post
```http
PUT /api/v1/posts/{postId}
```

### Delete Post
```http
DELETE /api/v1/posts/{postId}
```

### Like Post
```http
POST /api/v1/posts/{postId}/like
```

### Share Post
```http
POST /api/v1/posts/{postId}/share
```

### Get Post Comments
```http
GET /api/v1/posts/{postId}/comments?page=1&limit=20
```

### Increment View Count
```http
POST /api/v1/posts/{postId}/views
```

### Search Posts
```http
GET /api/v1/posts/search?q=search+query&page=1&limit=20
```

**Query Parameters:**
- `q`: Search query
- `page`: Page number
- `limit`: Items per page

### Get Posts by Category
```http
GET /api/v1/posts/category/{category}?page=1&limit=20
```

### Get Scheduled Posts
```http
GET /api/v1/posts/scheduled?page=1&limit=20
```

## Comments

### Create Comment
```http
POST /api/v1/comments
```

**Request Body:**
```json
{
  "postId": "post_uuid",
  "content": "Comment content",
  "parentId": "parent_comment_uuid" // Optional for replies
}
```

### Update Comment
```http
PUT /api/v1/comments/{commentId}
```

**Request Body:**
```json
{
  "content": "Updated comment content"
}
```

### Delete Comment
```http
DELETE /api/v1/comments/{commentId}
```

### Like Comment
```http
POST /api/v1/comments/{commentId}/like
```

### Get Comment Replies
```http
GET /api/v1/comments/{commentId}/replies?page=1&limit=20
```

## Messages

### Get Conversations
```http
GET /api/v1/messages?page=1&limit=20
```

### Get Messages with User
```http
GET /api/v1/messages/{userId}?page=1&limit=50
```

### Send Message
```http
POST /api/v1/messages
```

**Request Body:**
```json
{
  "receiverId": "user_uuid",
  "content": "Message content",
  "mediaUrl": "https://cloudinary.com/image.jpg" // Optional
}
```

### Mark Message as Read
```http
PUT /api/v1/messages/{messageId}/read
```

### Delete Message
```http
DELETE /api/v1/messages/{messageId}
```

## Notifications

### Get Notifications
```http
GET /api/v1/notifications?page=1&limit=20&unreadOnly=false
```

**Query Parameters:**
- `page`: Page number
- `limit`: Items per page
- `unreadOnly`: Filter unread only (boolean)

### Mark Notification as Read
```http
PUT /api/v1/notifications/{notifId}/read
```

### Mark All Notifications as Read
```http
PUT /api/v1/notifications/read-all
```

### Get Unread Count
```http
GET /api/v1/notifications/unread-count
```

**Response:**
```json
{
  "data": {
    "count": 5
  }
}
```

## Roles & Payments

### List Roles
```http
GET /api/v1/roles
```

### Get Role Details
```http
GET /api/v1/roles/{roleName}
```

### Create Payment
```http
POST /api/v1/payments
```

**Request Body:**
```json
{
  "role": "premium",
  "paymentMethod": "credit_card",
  "paymentType": "subscription"
}
```

**Response:**
```json
{
  "data": {
    "paymentId": "uuid",
    "orderId": "ORDER123",
    "snapToken": "snap_token",
    "snapRedirectUrl": "https://app.sandbox.midtrans.com/...",
    "expiryTime": "2025-11-15T10:00:00Z"
  }
}
```

### Midtrans Webhook
```http
POST /api/v1/payments/webhook
```

**Note:** This endpoint is called by Midtrans, not by client.

### Get User Payments
```http
GET /api/v1/payments?page=1&limit=20&status=SUCCESS
```

### Get Payment Details
```http
GET /api/v1/payments/{paymentId}
```

## Media Upload

### Upload Image
```http
POST /api/v1/upload/image
```

**Content-Type:** `multipart/form-data`

**Form Data:**
- `file`: Image file
- `folder` (optional): Cloudinary folder
- `transformation` (optional): Transformation options (JSON string)

**Response:**
```json
{
  "data": {
    "url": "https://res.cloudinary.com/...",
    "publicId": "folder/image_id",
    "format": "jpg",
    "width": 1920,
    "height": 1080
  }
}
```

### Upload Video
```http
POST /api/v1/upload/video
```

**Content-Type:** `multipart/form-data`

**Form Data:**
- `file`: Video file
- `folder` (optional): Cloudinary folder

### Batch Upload
```http
POST /api/v1/upload/batch
```

**Content-Type:** `multipart/form-data`

**Form Data:**
- `files[]`: Multiple files

## Admin Endpoints

### List All Users (Admin)
```http
GET /api/v1/admin/users?page=1&limit=20&isBanned=false
```

### Ban User (Admin)
```http
PUT /api/v1/admin/users/{userId}/ban
```

**Request Body:**
```json
{
  "reason": "Violation of community guidelines"
}
```

### Unban User (Admin)
```http
PUT /api/v1/admin/users/{userId}/unban
```

### List All Posts (Admin)
```http
GET /api/v1/admin/posts?page=1&limit=20&isPublished=false
```

### Publish Post (Admin)
```http
PUT /api/v1/admin/posts/{postId}/publish
```

### Hard Delete Post (Admin)
```http
DELETE /api/v1/admin/posts/{postId}
```

### Get Platform Statistics (Admin)
```http
GET /api/v1/admin/stats
```

**Response:**
```json
{
  "data": {
    "totalUsers": 1000,
    "totalPosts": 5000,
    "activeUsers": 500,
    "totalRevenue": 100000
  }
}
```

## Error Responses

All error responses follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {},
    "timestamp": "2025-11-15T09:00:00Z"
  }
}
```

### Common Error Codes

- `VALIDATION_ERROR` (400): Invalid input data
- `UNAUTHORIZED` (401): Authentication required
- `FORBIDDEN` (403): Insufficient permissions
- `NOT_FOUND` (404): Resource not found
- `CONFLICT` (409): Resource conflict (e.g., duplicate)
- `INTERNAL_ERROR` (500): Internal server error

## Rate Limiting

- **Default**: 100 requests per second per IP
- **Burst**: 200 requests
- **Headers**: Rate limit information in response headers
  - `X-RateLimit-Limit`: Request limit
  - `X-RateLimit-Remaining`: Remaining requests
  - `X-RateLimit-Reset`: Reset time

## Pagination

All list endpoints support pagination:

**Query Parameters:**
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 20, max: 100)

**Response:**
```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "totalPages": 5,
    "hasNext": true,
    "hasPrev": false
  }
}
```

