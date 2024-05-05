# Gun Management System

## About the Project

This project is a web application designed to efficiently manage the catalog of guns, sales, and customers in a gun shop. The application allows shop owners to handle their inventory, process sales, and interact with customers.

### Key Features

- **Gun Catalog Management:** Add, update, delete, and view information about guns.

## Rest API

### Gun Endpoints:

- **POST /guns** - Add a new gun to the catalog.
- **GET /guns/{id}** - Retrieve information about a gun by ID.
- **PUT /guns/{id}** - Update information about a gun by ID.
- **DELETE /guns/{id}** - Remove a gun from the catalog.



## Database Structure and Relationships

The project uses the following database structure for managing the gun catalog, sales, and customers.

### Entities:

#### Guns (`guns`)

- `id` (bigserial): Unique identifier for the gun (primary key).
- `name` (text): Name of the gun.
- `price` (numeric): Price of the gun.
- `damage` (integer): Damage level of the gun.
- `created_at` (timestamp): Date and time the gun was added to the catalog.
- `updated_at` (timestamp): Date and time the gun information was last updated.




---

### Sagataly Aidyn 22B030424
