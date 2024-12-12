## Problem Statement

Users are unable to receive push notifications for incoming/outgoing transactions. This parser solves the problem by enabling:

- Tracking of Ethereum blockchain transactions for subscribed addresses.
- Querying inbound and outbound transactions for specific addresses.
- Providing an extendable memory-based storage solution.
- Exposing a public interface via an HTTP API.

## Highlights

1. **Parser Implementation:**
   - Methods: `GetCurrentBlock()`, `Subscribe(address)`, `GetTransactions(address)`, `PollBlocks(intervalSeconds)`, `ProcessBlock(blockNumber)`.
   - Regularly polls the Ethereum blockchain and processes transactions for subscribed addresses.

2. **Memory Storage:**
   - Implements `SubscriptionStorage` and `TransactionStorage` interfaces for extendable storage.

3. **Concurrency-Safe:**
   - Utilizes `sync.Mutex` for thread-safe operations on shared data.

4. **HTTP API:**
   - Endpoints:
     - **POST /subscribe:** Subscribe to an Ethereum address.
     - **GET /transactions:** Fetch transactions for a subscribed address.
     - **GET /currentBlock:** Get the latest processed block.

5. **Unit Testing:**
   - Tests for `SubscriptionStorage`, `TransactionStorage`, and `ParserService`.
   - Run tests with:
     ```bash
     go test ./... -v
     ```

## How to Run

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd <repository-directory>
   ```

2. Run the application:
   ```bash
   go run .
   ```

3. The server will start on `http://localhost:8080`.

## API Endpoints

- **POST /subscribe**
  ```json
  {
      "address": "0xYourEthereumAddress"
  }
  ```
  Example:
  ```bash
  curl -X POST http://localhost:8080/subscribe -H "Content-Type: application/json" -d '{"address": "0xYourEthereumAddress"}'
  ```

- **GET /transactions?address=0xYourEthereumAddress**
  Example:
  ```bash
  curl "http://localhost:8080/transactions?address=0xYourEthereumAddress"
  ```

- **GET /currentBlock**
  Example:
  ```bash
  curl http://localhost:8080/currentBlock
  ```