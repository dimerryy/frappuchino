# frappuccino

## Overview

**frappuccino** is a backend service for a coffee shop management system. It allows managing menu items, orders, and inventory with a set of RESTful API endpoints. The system also supports reporting features and bulk order processing using PostgreSQL transactions.

## Features

- **Orders**:
  - Create, update, and close orders.
  - Batch processing of multiple orders with transactional consistency.

- **Menu Management**:
  - Add, list, and delete menu items and their ingredients.

- **Inventory Management**:
  - Track and update inventory quantities.
  - Retrieve leftovers with sorting and pagination.

- **Reports**:
  - Generate daily or monthly reports for ordered items.

## Clone the Repository

```bash
git clone https://github.com/dimerryy/frappuccino.git
cd frappuccino
```

## Run the Server

1. **Start PostgreSQL via Docker Compose**:

```bash
docker-compose up -d
```

2. **Run the Backend**:

```bash
go run cmd/main.go
```

## API Endpoints

### Orders

#### Create Order

```bash
POST /orders
```

#### Close Order

```bash
PUT /orders/{id}/close
```

#### Batch Process Orders

```bash
POST /orders/batch-process
```

### Menu

#### Create Menu Item

```bash
POST /menu
```

#### Delete Menu Item

```bash
DELETE /menu/{id}
```

### Inventory

#### Get Leftovers

```bash
GET /inventory/getLeftOvers?sortBy=quantity&page=1&pageSize=5
```

### Reports

#### Get Ordered Items by Period (Day)

```bash
GET /reports/orderedItemsByPeriod?period=day&month=april
```

#### Get Ordered Items by Period (Month)

```bash
GET /reports/orderedItemsByPeriod?period=month&year=2023
```

## Models

```go
// OrderItem
type OrderItem struct {
    MenuItemID int `json:"menu_item_id"`
    Quantity   int `json:"quantity"`
}

// BatchOrderRequest
type BatchOrderRequest struct {
    Orders []struct {
        CustomerName string       `json:"customer_name"`
        Items        []OrderItem  `json:"items"`
    } `json:"orders"`
}
```

## Usage Tips

- Use [Postman](https://www.postman.com/) to test endpoints with JSON payloads.
- Ensure the `init.sql` runs inside the Docker container to initialize your database schema.


