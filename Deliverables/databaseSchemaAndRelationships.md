# Database Schema

## Overview

PostgreSQL database with 2 core tables: 

---

## Table: `stock_rewards`

**Purpose:** Records all stock rewards given to users.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | VARCHAR(255) | PRIMARY KEY | Unique reward identifier (idempotency key) |
| `user_id` | VARCHAR(255) | NOT NULL, INDEXED | User identifier |
| `stock_symbol` | VARCHAR(50) | NOT NULL | Stock ticker symbol |
| `quantity` | NUMERIC(18,6) | NOT NULL | Number of shares (supports fractional) |
| `reward_timestamp` | TIMESTAMP | NOT NULL, INDEXED | When reward was granted |
| `stock_price_at_reward` | NUMERIC(18,4) | NOT NULL | Stock price at reward time |
| `created_at` | TIMESTAMP | AUTO | Record creation time |

**Indexes:**
- Primary Key on `id`
- Index on `user_id` (for fast user queries)
- Index on `reward_timestamp` (for date-based queries)

---

## Table: `ledger_entries`

**Purpose:** Double-entry accounting ledger for all financial transactions.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| `id` | SERIAL | PRIMARY KEY | Auto-increment ID |
| `reward_id` | VARCHAR(255) | FOREIGN KEY, NOT NULL, INDEXED | References stock_rewards.id |
| `account_type` | VARCHAR(50) | NOT NULL | Type of account entry |
| `stock_symbol` | VARCHAR(50) | NULLABLE | Stock symbol (for asset entries) |
| `debit_amount` | NUMERIC(18,4) | NOT NULL, DEFAULT 0 | Debit amount in INR |
| `credit_amount` | NUMERIC(18,4) | NOT NULL, DEFAULT 0 | Credit amount in INR |
| `quantity` | NUMERIC(18,6) | NULLABLE | Stock quantity (for asset entries) |
| `description` | TEXT | | Human-readable description |
| `created_at` | TIMESTAMP | AUTO | Record creation time |

**Account Types:**
- `STOCK_ASSET`: Stock holdings acquired (debit)
- `CASH_ACCOUNT`: Company cash account (credit when paying)
- `BROKERAGE_EXPENSE`: Brokerage charges (debit)
- `STT_EXPENSE`: Securities Transaction Tax (debit)
- `GST_EXPENSE`: Goods and Services Tax (debit)

**Foreign Keys:**
- `reward_id` references `stock_rewards(id)` with CASCADE delete

**Indexes:**
- Primary Key on `id`
- Foreign Key index on `reward_id`

## Relationships

```
stock_rewards (1) ──→ (N) ledger_entries
     │
     │ One reward generates 5 balanced ledger entries:
     │ 
     │ DEBITS (what we acquired/spent):
     │ - 1x STOCK_ASSET (debit: stock value)
     │ - 1x BROKERAGE_EXPENSE (debit: 0.03% of stock value)
     │ - 1x STT_EXPENSE (debit: 0.1% of stock value)
     │ - 1x GST_EXPENSE (debit: 18% of brokerage)
     │ 
     │ CREDITS (where cash came from):
     │ - 1x CASH_ACCOUNT (credit: total cash outflow)
     │
     └─ Total Debits = Total Credits (balanced accounting)
```

---

## Double-Entry Ledger Example

For a reward of **10.5 shares of RELIANCE at ₹2450.75**:

**Calculations:**
- Stock Cost: 10.5 × ₹2450.75 = ₹25,732.88
- Brokerage (0.03%): ₹25,732.88 × 0.0003 = ₹7.71
- STT (0.1%): ₹25,732.88 × 0.001 = ₹25.73
- GST (18% of brokerage): ₹7.71 × 0.18 = ₹1.38
- **Total Cash Outflow**: ₹25,732.88 + ₹7.71 + ₹25.73 + ₹1.38 = **₹25,767.70**

**Ledger Entries:**

| Account Type | Debit | Credit | Description |
|--------------|-------|--------|-------------|
| STOCK_ASSET | 25,732.88 | 0 | Stock acquired (10.5 shares) |
| BROKERAGE_EXPENSE | 7.71 | 0 | Brokerage fee (0.03%) |
| STT_EXPENSE | 25.73 | 0 | Securities Transaction Tax (0.1%) |
| GST_EXPENSE | 1.38 | 0 | GST on brokerage (18%) |
| CASH_ACCOUNT | 0 | 25,767.70 | Total cash paid |
| **TOTAL** | **25,767.70** | **25,767.70** | **BALANCED** |

**Key Points:**
- **Total Debits = Total Credits** (25,767.70)
- Stock and all expenses are DEBITS
- Cash outflow is the single CREDIT
- Proper double-entry accounting

---
