# Products REST API

This project is a REST API for managing products, built with Go and using PostgreSQL as the database. It includes setup instructions for running locally using Docker Compose.

## Table of Contents

- [Products REST API](#products-rest-api)
  - [Table of Contents](#table-of-contents)
  - [Prerequisites](#prerequisites)
  - [Project Setup](#project-setup)
    - [1. Clone the Repository](#1-clone-the-repository)
    - [2. Environment Variables](#2-environment-variables)
    - [3. Build and Run with Docker Compose](#3-build-and-run-with-docker-compose)
    - [4. Build Binary and Run (Without DB)](#4-build-binary-and-run-without-db)
  - [Usage](#usage)
    - [API Endpoints](#api-endpoints)
  - [Next on the List](#next-on-the-list)

## Prerequisites

- [Docker](https://www.docker.com/get-started) (v20+ recommended)
- [Docker Compose](https://docs.docker.com/compose/install/) (v1.29+ recommended)
- Go and Taskfile (if running locally without Docker)

## Project Setup

### 1. Clone the Repository

Clone this repository to your local machine:

```bash
git clone https://github.com/mariosker/products_rest_api.git
cd products_rest_api
```

### 2. Environment Variables

This project requires database connection details. Ensure your environment variables are correctly set in your docker-compose.yml file or the .env file. Here’s an example:

```bash
DBURL=postgres://postgres:postgres@db:5432/products?sslmode=disable
GIN_MODE=release
```

### 3. Build and Run with Docker Compose

To build and run the API with Docker Compose:

```bash
docker compose up --build
```

This command:

- Starts the db container running PostgreSQL.
- Builds and runs the API in the api container.
- Exposes the API on port 8080 and the database on port 5432.
- You can access the API at <http://localhost:8080>.

### 4. Build Binary and Run (Without DB)

To build and run the API using Go:

```bash
task run
```

## Usage

### API Endpoints

Below are examples of available endpoints:

- `POST /products`: Create a new product.
- `GET /products`: Get all products, with optional query parameters limit and offset.
- `GET /products/:id:` Get a product by ID.
- `PUT /products/:id:` Update a product by ID.
- `DELETE /products/:id:` Delete a product by ID.

The Create, Update commands want a JSON in the form of:

```json
"name": "Example",
"price": 100
```

## Next on the List

- [ ] Implement multiple currencies
- [ ] Add support for filtering products by price and name.
- [ ] Create comprehensive API documentation (e.g., using Swagger).
- [ ] Implement logging and monitoring for the API.
- [ ] Consider adding a caching layer (e.g., Redis) for frequently accessed data.
- [ ] Implement rate limiting to prevent abuse of the API.