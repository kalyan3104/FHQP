# ServU Go API

This project runs as two processes from the same Go codebase:

- Booking API process: handles normal HTTP APIs and publishes booking notification events to RabbitMQ.
- Notification worker process: consumes RabbitMQ events, creates notification records, sends FCM push notifications, and owns the websocket endpoint.

## Requirements

- Go
- PostgreSQL/Supabase connection
- RabbitMQ
- Firebase service account JSON

## Environment Variables

Required for database:

```powershell
$env:SUPABASE_DB_URL="postgres://..."
```

Required for RabbitMQ if not using local default:

```powershell
$env:RABBITMQ_URL="amqp://guest:guest@localhost:5672/"
```

Optional ports:

```powershell
$env:PORT="8001"
$env:NOTIFICATION_PORT="8002"
```

Default ports:

- Booking API: `8001`
- Notification worker websocket: `8002`

## Run RabbitMQ Locally

If using Docker:

```powershell
docker run -d --name servu-rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management
```

RabbitMQ dashboard:

```text
http://localhost:15672
```

Default login:

```text
username: guest
password: guest
```

## Run Booking API

Start the normal API process:

```powershell
go run .
```

Booking API runs on:

```text
http://localhost:8001
```

Create booking endpoint:

```text
POST http://localhost:8001/bookings/
```

Example body:

```json
{
  "customer_id": "cust1",
  "provider_id": "prov1",
  "slot_time": "2026-05-05T14:30:00Z"
}
```

When booking is created, the API publishes a `BOOKING_CREATED` event to RabbitMQ. The API does not directly send websocket, FCM, or notification DB records.

## Run Notification Worker

Start the notification worker in a separate terminal:

```powershell
$env:APP_MODE="notification-worker"
go run .
```

Notification worker runs websocket on:

```text
ws://localhost:8002/ws?user_id=<provider_id>
```

Example:

```text
ws://localhost:8002/ws?user_id=prov1
```

Use the provider ID as `user_id`, because booking notifications are sent to the provider.

## Notification Flow

Full flow:

```text
POST /bookings/
  -> Booking API creates booking in DB
  -> Booking API publishes BOOKING_CREATED event to RabbitMQ
  -> Booking API returns response

Notification worker
  -> consumes BOOKING_CREATED from RabbitMQ
  -> creates notification record in DB
  -> sends Firebase push notification if FCM token exists
  -> sends websocket message if provider is connected
```

Websocket payload example:

```json
{
  "type": "NEW_BOOKING",
  "booking_id": "...",
  "title": "New Booking",
  "message": "You received a booking request"
}
```

## RabbitMQ Queue

The integration uses:

```text
exchange: servu.notifications
queue: servu.booking.notifications
routing key: booking.created
```

If the notification worker is running, RabbitMQ queue counts may stay at zero because messages are consumed immediately.

To verify publishing:

1. Stop the notification worker.
2. Create a booking.
3. Check RabbitMQ dashboard.
4. Queue `servu.booking.notifications` should show `Ready: 1`.
5. Start notification worker.
6. Message should be consumed and queue should return to `Ready: 0`.

## Websocket Notes

The websocket connections are stored in memory in the notification worker process.

That means clients must connect to:

```text
ws://localhost:8002/ws?user_id=<provider_id>
```

Do not connect websocket to the booking API:

```text
ws://localhost:8001/ws?user_id=<provider_id>
```

That URL will return `404` because the booking API no longer owns `/ws`.

In production, use `wss://`:

```text
wss://notifications.your-domain.com/ws?user_id=<provider_id>
```

For the current design, run one notification websocket worker. If multiple notification workers are added later, websocket delivery needs Redis pub/sub, RabbitMQ fanout, or a separate websocket server so the message reaches the server where the user is connected.
