# RabbitMQ Integration

## Mengapa Menggunakan RabbitMQ?

### Alasan Utama

1. **Non-Blocking Operations**
   - Email sending tidak blocking HTTP request
   - User tidak perlu menunggu email dikirim sebelum mendapat response
   - Request lebih cepat dan responsif

2. **Scalability**
   - Dapat handle banyak email requests secara bersamaan
   - Queue memastikan tidak ada email yang hilang
   - Dapat scale workers sesuai kebutuhan

3. **Reliability**
   - Jika email service down, jobs akan tersimpan di queue
   - Email akan dikirim setelah service kembali online
   - Messages persistent (tidak hilang saat restart)

4. **Separation of Concerns**
   - API server fokus pada HTTP handling
   - Email processing dilakukan oleh dedicated workers
   - Lebih mudah untuk maintain dan debug

5. **Future Extensibility**
   - Mudah menambahkan job types baru (notifications, feed processing, etc.)
   - Dapat prioritaskan jobs (priority queue)
   - Dapat delay jobs (scheduled emails)

## Architecture

### Flow Diagram

```
HTTP Request (Register/Forgot Password)
    ↓
Auth Service
    ↓
Publish to RabbitMQ Queue (Non-blocking)
    ↓
Return Success Response (Fast response)
    ↓
[In Background]
Email Worker consumes message
    ↓
Process Email (SMTP Send)
    ↓
Email Sent
```

### Queue Structure

1. **email_queue** - Email sending jobs
   - Email verification
   - Password reset emails
   - Future: Welcome emails, notifications, etc.

2. **notification_queue** - User notifications
   - (Untuk implementasi future)

3. **feed_processing_queue** - Feed generation
   - (Untuk implementasi future)

## Implementation Details

### Files

1. **pkg/mq/rabbitmq.go** - RabbitMQ connection & setup
2. **pkg/mq/publisher.go** - Message publishing
3. **internal/workers/email_worker.go** - Email worker consumer
4. **internal/services/auth_service.go** - Updated untuk publish ke queue

### Connection Flow

```go
// 1. Connect to RabbitMQ
mq.Connect()

// 2. Start Email Worker
workers.StartEmailWorker()

// 3. Publish messages
mq.PublishVerificationEmail(email, token)
```

### Message Format

```json
{
  "type": "verification" | "reset_password",
  "to": "user@example.com",
  "token": "verification_token",
  "subject": "Email Subject",
  "body": "Email Body"
}
```

## Fallback Mechanism

Jika RabbitMQ tidak tersedia atau gagal:

1. **Graceful Degradation**
   - System tetap berfungsi
   - Fallback ke direct email sending
   - User tidak merasakan perbedaan

2. **Error Handling**
   ```go
   if err := mq.PublishEmail(message); err != nil {
       // Fallback to direct email
       utils.SendEmail(...)
   }
   ```

## Benefits untuk Auth Flow

### Before (Without RabbitMQ)
```
Register Request
    ↓
Create User (100ms)
    ↓
Send Email (2000ms) ← BLOCKING
    ↓
Return Response (2100ms total)
```

### After (With RabbitMQ)
```
Register Request
    ↓
Create User (100ms)
    ↓
Publish to Queue (5ms) ← NON-BLOCKING
    ↓
Return Response (105ms total) ← 20x FASTER!
    ↓
[Background] Send Email (2000ms)
```

## Monitoring

### RabbitMQ Management UI

Access: http://localhost:15672
- Username: lostmediago
- Password: password123

Dari UI, Anda dapat:
- Monitor queue length
- Check message rates
- Inspect failed messages
- Monitor consumers

### Logs

Email worker akan log:
```
Email worker started, waiting for messages...
Processing email: type=verification, to=user@example.com
Verification email sent to: user@example.com
```

## Configuration

### Environment Variables

```env
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=lostmediago
RABBITMQ_PASSWORD=password123
RABBITMQ_VHOST=/
```

### Docker Compose

RabbitMQ sudah dikonfigurasi di `docker-compose.yml`:
- Port 5672: AMQP connection
- Port 15672: Management UI

## Future Enhancements

1. **Priority Queue**
   - Urgent emails (password reset) di priority tinggi
   - Marketing emails di priority rendah

2. **Scheduled Emails**
   - Welcome email setelah 24 jam
   - Reminder emails

3. **Email Templates**
   - HTML email templates
   - Multi-language support

4. **Retry Mechanism**
   - Auto retry failed emails
   - Exponential backoff

5. **Monitoring & Alerting**
   - Queue length alerts
   - Failed message alerts
   - Email delivery tracking

## Troubleshooting

### Queue Not Consuming

1. Check if worker is running
2. Check RabbitMQ connection
3. Check queue exists
4. Check consumer is registered

### Messages Not Published

1. Check RabbitMQ connection
2. Check queue exists
3. Check permissions
4. Check channel is open

### Emails Not Sending

1. Check SMTP configuration
2. Check email worker logs
3. Check message format
4. Check email service status

## Best Practices

1. **Always have fallback** - System harus tetap bekerja tanpa RabbitMQ
2. **Monitor queue length** - Jangan biarkan queue terlalu panjang
3. **Handle errors gracefully** - Log errors, retry jika perlu
4. **Use persistent messages** - Jangan hilang saat restart
5. **Limit message size** - Keep messages small
6. **Use appropriate queues** - Separate queues untuk different job types

