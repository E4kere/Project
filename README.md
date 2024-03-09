# Messanger project

## Rest API
- `POST /register`
- `GET /login`
- `PATCH /login/{type}`
- `DELETE /login`
- `GET /users`
- `POST /message/{id}`
- `DELETE /message/{id}`


## Database Structure

```sql
CREATE DATABASE cs2_pro_players;

USE cs2_pro_players;

CREATE TABLE Players (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    team VARCHAR(255),
    country VARCHAR(255),
    date_of_birth DATE
);

CREATE TABLE PlayerStats (
    id INT AUTO_INCREMENT PRIMARY KEY,
    player_id INT,
    date DATE,
    kills INT,
    deaths INT,
    kd_ratio DECIMAL(5,2),
    maps_played INT,
    total_rounds INT,
    rounds_per_map DECIMAL(5,2),
    damage_per_round DECIMAL(8,2),
    grenade_dmg_per_round DECIMAL(8,2),
    FOREIGN KEY (player_id) REFERENCES Players(id)
);


