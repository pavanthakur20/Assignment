# Stock Reward Management System

A Go-based REST API system for managing stock rewards with double-entry ledger accounting.

## Features

- **Stock Reward Management**: Record and track stock rewards for users
- **Double-Entry Ledger**: Properly balanced accounting system tracking stock units, cash flow, and fees
- **Portfolio Tracking**: View user portfolios with current valuations
- **Statistics**: Daily and historical statistics for user rewards
- **Idempotency**: Prevent duplicate reward processing
- **Note**: ## 📊 Supported Stocks

The system supports the following Indian stocks:

- RELIANCE
- TCS
- INFOSYS
- HDFC
- WIPRO
- ITC
- BHARTI
- SBIN
- HDFCBANK
- ICICIBANK

## Project Structure

```
assignment/
├── server.go                 # Main entry point
├── docker-compose.yml        # Database setup
├── .env                      # Environment variables
│
├── controllers/              # Request handlers
│   ├── rewardController.go
│   ├── portfolioController.go
│   ├── statsController.go
│   ├── historicalController.go
│   └── todayStocksController.go
│
├── services/                 # Business logic
│   ├── ledgerService.go
│   └── stockPriceService.go
│
├── models/                   # Data models
│   ├── stockReward.go
│   ├── ledger.go
│   └── stockPrice.go
│
├── initializers/             # Setup & config
│   ├── database.go
│   └── loadEnv.go
│
└── Deliverables/            # Documentation
    ├── apiSpecs.md
    └── databaseSchemaAndRelationships.md
```

## Getting Started

### Prerequisites

- Go 1.25+ installed
- Docker and Docker Compose installed
- Git

### 1. Clone the Repository

```bash
git clone https://github.com/pavanthakur20/Assignment.git
cd assignment
```

### 2. Set Up the Database

Start PostgreSQL using Docker Compose:

```bash
docker-compose up -d
```

This will:
- Start PostgreSQL 15 on port `5432`
- Create database `assignment`
- Username: `user`
- Password: `user`
- Store data in a persistent volume

Check if the database is running:

```bash
docker ps
```

You should see `assignment_postgres` container running.

### 3. Configure Environment Variables

The `.env` file is already configured for local development:

```env
PORT=8080
DB_URL=postgresql://user:user@localhost:5432/assignment?sslmode=disable
```

### 4. Install Dependencies

```bash
go mod download
```

### 5. Run the Application

```bash
go run server.go
```

You should see:

```
Connected to database successfully
Database migrations completed successfully
Stock price update scheduler started
Stock prices updated
Server starting on port 8080
```

The API is now running at `http://localhost:8080`

## API Endpoints

### Base URL
```
http://localhost:8080/api/v1
```

### Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/reward` | Create a stock reward |
| GET | `/today-stocks/:userId` | Get today's rewards for a user |
| GET | `/historical-inr/:userId` | Get historical INR values |
| GET | `/stats/:userId` | Get user statistics |
| GET | `/portfolio/:userId` | Get user portfolio |

## Documentation

- [API Specifications](./Deliverables/apiSpecs.md)
- [Database Schema](./Deliverables/databaseSchemaAndRelationships.md)
- [Postman Collection](./Postman/Assignment.postman_collection.json)
