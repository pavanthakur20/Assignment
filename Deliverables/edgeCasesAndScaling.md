# Edge Cases & Scaling Implementation

## Edge Cases Handled

### 1. Idempotency (Duplicate Rewards)
- **Implementation:** Primary key on `stock_rewards.id`
- **Check:** Database query before processing reward
- **Response:** Returns `409 Conflict` if reward ID already exists
- **Benefit:** Same reward never processed twice, even with network retries

### 2. Fractional Shares
- **Implementation:** `NUMERIC(18,6)` database precision
- **Validation:** Accepts any positive decimal (e.g., 0.5, 10.25)
- **Example:** User can receive 0.5 shares of HDFC
- **Benefit:** Supports partial stock rewards accurately

### 3. Zero or Negative Quantities
- **Implementation:** Request validation with `binding:"required,gt=0"`
- **Response:** Returns `400 Bad Request` before database access
- **Benefit:** Invalid data blocked at API layer

### 4. Concurrent Duplicate Requests
- **Implementation:** Database transaction with primary key constraint
- **Behavior:** First request succeeds, others get constraint violation
- **Response:** `409 Conflict` for duplicates
- **Benefit:** Race condition protected

### 5. Data Consistency (Partial Failures)
- **Implementation:** GORM database transaction wrapping all operations
- **Operations:** Reward creation + 5 ledger entries in single transaction
- **Rollback:** Any failure rolls back entire transaction
- **Benefit:** No orphaned records or inconsistent state

### 6. Large Quantity Values
- **Implementation:** `NUMERIC(18,6)` supports up to 12 digits before decimal
- **Example:** Can handle 1,000,000+ shares without overflow
- **Benefit:** No integer overflow issues

### 7. Timezone Handling
- **Implementation:** ISO 8601 timestamp format in API
- **Storage:** All timestamps stored as `TIMESTAMP` (UTC)
- **Conversion:** Go's `time.Time` handles automatic conversion
- **Benefit:** Consistent time handling across timezones

### 8. Database Indexing
- **Indexes Created:**
  - Primary key on `stock_rewards.id`
  - Index on `stock_rewards.user_id`
  - Index on `stock_rewards.reward_timestamp`
  - Foreign key index on `ledger_entries.reward_id`
  - Unique composite index on `stock_prices(stock_symbol, timestamp)`
- **Benefit:** Fast queries for user portfolios and historical data

### 9. Background Price Updates
- **Implementation:** Goroutine with time ticker
- **Frequency:** Runs every 1 hour
- **Execution:** Non-blocking, runs in separate goroutine
- **Benefit:** Price updates don't block API requests

### 10. Thread-Safe Price Cache
- **Implementation:** In-memory map with `sync.RWMutex`
- **Read Lock:** Multiple reads can happen concurrently
- **Write Lock:** Exclusive lock during price updates
- **Benefit:** Fast price lookups without database queries

### 11. Batch Ledger Inserts
- **Implementation:** All 5 ledger entries inserted in single GORM operation
- **Code:** `tx.Create(&ledgerEntries)` with slice
- **Benefit:** Reduces database round trips

### 12. Foreign Key Cascade Delete
- **Implementation:** `constraint:OnDelete:CASCADE` on ledger entries
- **Behavior:** Deleting reward automatically deletes associated ledger entries
- **Benefit:** Maintains referential integrity automatically

### 13. Auto Migrations
- **Implementation:** GORM `AutoMigrate()` at startup
- **Models:** `StockReward`, `LedgerEntry`, `StockPrice`
- **Benefit:** Database schema updates automatically on deployment

### 14. Error Handling & Logging
- **Implementation:** Gin middleware with logger
- **Database Errors:** Captured and returned with proper HTTP status codes
- **Benefit:** Clear error messages for debugging
