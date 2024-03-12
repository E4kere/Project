# Inventory Management System for a Small Gun Store

## About the Project

This project is a web application for managing CS2 guns. It allows users to perform CRUD (Create, Read, Update, Delete) operations on gun records. Users can add new guns, view existing guns, update gun information, and delete guns from the system.

##Technologies Used

- **Backend: Go (Golang)**
- **Database: PostgreSQL**
- **Web Framework: Gorilla Mux**

### Key Features

- **Catalog Management**: Add, update, delete, and view information about guns.

## Rest API
- **POST /guns** - Add a new gun to the catalog.
- **GET /guns/:id** - Get information about a gun by ID.
- **PUT /guns/:id** - Update gun information by ID.
- **DELETE /guns/:id** - Delete a gun from the catalog.
  


## Database Entity Structure and Relationships

The project uses the following database structure to manage the catalog of guns, orders, and customers.

### Entities:

#### Guns (`guns`)

- `id` (bigserial): Unique identifier for the gun (primary key).
- `model` (text): Model of the gun.
- `brand` (text): Brand of the gun.
- `price` (numeric): Price of the gun.
- `damage` (integer): Damage of the gun.
- `created_at` (timestamp): Date and time the gun was added to the catalog.
- `updated_at` (timestamp): Date and time the gun information was last updated.

### Developer Information:

- Sagataly Aidyn
- Student ID: 22B030424
