# Messanger project

## Rest API
- `POST /register`
- `GET /login`
- `PATCH /login/{type}`
- `DELETE /login`
- `GET /users`
- `POST /message/{id}`
- `UPDATE /message/{id}`
- `DELETE /message/{id}`
- `GET /notifications`

## Database Structure

```sql
Table Users {
  user_id SERIAL [pk]
  username VARCHAR(50) [not null]
  password VARCHAR(100) [not null]
}

Table Messages {
  message_id SERIAL [pk]
  sender_id INTEGER [not null, ref: > Users.user_id]
  receiver_id INTEGER [not null, ref: > Users.user_id]
  message_text TEXT [not null]
  sent_at TIMESTAMP [default: `CURRENT_TIMESTAMP`]
}

