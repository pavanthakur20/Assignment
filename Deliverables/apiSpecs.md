# API Specifications

## Base URL
```
http://localhost:8080/api/v1
```

## 1. Reward User Endpoint

### POST `/reward`

**Purpose:** Creates a stock reward for a user with automatic ledger entries and charge calculations.

**Request Payload:**
```json
{
  "id": "reward_reliance_001",
  "user_id": "user_123",
  "stock_symbol": "RELIANCE",
  "quantity": 10.5,
  "reward_timestamp": "2025-11-17T10:30:00Z"
}
```

**Request Fields:**
- `id` (string, required): Unique reward identifier for idempotency
- `user_id` (string, required): User identifier
- `stock_symbol` (string, required): Stock ticker symbol
- `quantity` (number, required): Number of shares (supports fractional)
- `reward_timestamp` (ISO 8601, required): When the reward was granted

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Stock reward processed successfully",
  "reward": {
    "id": "reward_reliance_001",
    "user_id": "user_123",
    "stock_symbol": "RELIANCE",
    "quantity": 10.5,
    "reward_timestamp": "2025-11-17T10:30:00Z",
    "stock_price_at_reward": 2450.75,
    "created_at": "2025-11-17T10:30:05Z"
  },
  "inr_value": 25732.88,
  "company_charges": {
    "stock_cost": 25732.88,
    "brokerage": 7.71,
    "stt": 25.73,
    "gst": 1.38,
    "total_cost": 25767.70
  }
}
```

**Error Responses:**

*400 Bad Request:*
```json
{
  "success": false,
  "message": "Invalid request: quantity must be greater than 0"
}
```

*409 Conflict (Idempotent duplicate):*
```json
{
  "success": false,
  "message": "Reward with ID 'reward_reliance_001' has already been processed"
}
```

*500 Internal Server Error:*
```json
{
  "success": false,
  "message": "Failed to get stock price: database connection error"
}
```

---

## 2. Today's Stocks Endpoint

### GET `/today-stocks/:userId`

**Purpose:** Retrieves all stocks rewarded to a user today.

**Request:** No payload required

**Success Response (200 OK):**
```json
{
  "user_id": "user_123",
  "date": "2025-11-17",
  "stocks": [
    {
      "stock_symbol": "RELIANCE",
      "total_quantity": 10.5,
      "current_price": 2450.75,
      "current_value": 25732.88
    },
    {
      "stock_symbol": "TCS",
      "total_quantity": 5.25,
      "current_price": 3800.50,
      "current_value": 19952.63
    }
  ],
  "total_value": 45685.51
}
```

**Error Response (500):**
```json
{
  "error": "Failed to retrieve today's stocks"
}
```

---

## 3. Historical INR Values Endpoint

### GET `/historical-inr/:userId`

**Purpose:** Returns daily portfolio values for a user.

**Request:** No payload required

**Success Response (200 OK):**
```json
{
  "user_id": "user_123",
  "daily_values": {
    "2025-11-15": 42000.50,
    "2025-11-16": 43500.75,
    "2025-11-17": 45685.51
  }
}
```

**Error Response (500):**
```json
{
  "error": "Failed to calculate historical INR values"
}
```

---

## 4. User Statistics Endpoint

### GET `/stats/:userId`

**Purpose:** Provides summary statistics for a user.

**Request:** No payload required

**Success Response (200 OK):**
```json
{
  "user_id": "user_123",
  "today_rewards": {
    "RELIANCE": 10.5,
    "TCS": 5.25,
    "INFOSYS": 15.75
  },
  "current_portfolio_inr": 45685.51,
  "total_shares_rewarded": 31.5
}
```

**Error Response (500):**
```json
{
  "error": "Failed to retrieve user statistics"
}
```

---

## 5. User Portfolio Endpoint

### GET `/portfolio/:userId`

**Purpose:** Shows complete portfolio with current holdings and values.

**Request:** No payload required

**Success Response (200 OK):**
```json
{
  "user_id": "user_123",
  "holdings": [
    {
      "stock_symbol": "RELIANCE",
      "total_quantity": 10.5,
      "current_price": 2450.75,
      "current_value": 25732.88
    },
    {
      "stock_symbol": "TCS",
      "total_quantity": 5.25,
      "current_price": 3800.50,
      "current_value": 19952.63
    }
  ],
  "total_value": 45685.51,
  "last_updated": "2025-11-17T14:30:00Z"
}
```

**Error Response (500):**
```json
{
  "error": "Failed to retrieve portfolio"
}
```

---

## Common Headers

**All Requests:**
```
Content-Type: application/json
```

**All Responses:**
```
Content-Type: application/json
```

---

## Charge Calculations

The system automatically applies these charges on stock rewards:

- **Brokerage Fee:** 0.03% of stock cost
- **STT (Securities Transaction Tax):** 0.1% of stock cost
- **GST:** 18% of brokerage fee
- **Total Cost:** Stock Cost + Brokerage + STT + GST

All charges are rounded to 2 decimal places.
