Student info:
- 'Sagataly Aidyn'
- 22B030424




# Messanger project

## Rest API
- `POST /register`
- `GET /login`
- `DELETE /login`
- `GET /users`
- `POST /{id}`
- `UPDATE /{id}`


## Database Structure

```sql
CREATE TABLE IF NOT EXISTS Players (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255),
			team VARCHAR(255),
			country VARCHAR(255),
			date_of_birth DATE
	)
