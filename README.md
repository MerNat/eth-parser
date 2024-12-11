## Problem Statement
Users are unable to receive push notifications for incoming/outgoing transactions. This parser solves the problem by enabling:

- Tracking of Ethereum blockchain transactions for subscribed addresses.
- Querying inbound and outbound transactions for specific addresses.
- Providing an extendable memory-based storage solution.
- Exposing a public interface via an HTTP API.

## Implementation Highlights

### Requirements Addressed

1. **Parser Interface Implementation:**
   - Implements the `Parser` interface with the following methods:
     - `GetCurrentBlock() int`: Retrieves the last parsed Ethereum block.
     - `Subscribe(address string) bool`: Subscribes an address for transaction tracking.
     - `GetTransactions(address string) []Transaction`: Fetches inbound and outbound transactions for a subscribed address.

2. **Memory Storage with Extensibility:**
   - Used in-memory storage for subscriptions and transactions with interfaces (`SubscriptionStorage` and `TransactionStorage`) to ensure future adaptability for other storage mechanisms.

3. **Safe Concurrency:**
   - Used `sync.Mutex` to ensure thread-safe operations, avoiding race conditions when accessing shared data like the current block and storage maps.

4. **Polling for New Blocks:**
   - Continuously polls the Ethereum blockchain for new blocks using the `PollBlocks` function. This ensures that transactions are regularly updated for subscribed addresses.

5. **HTTP API:**
   - Exposed the following endpoints to support the parser functionality:
     - **POST /subscribe:** Subscribe an Ethereum address.
     - **GET /transactions:** Retrieve transactions for a subscribed address.
     - **GET /currentBlock:** Get the latest processed block.

6. **Ethereum JSON-RPC Usage:**
   - Interacts with the Ethereum blockchain via JSON-RPC to fetch block data and transactions.

## How to Run

1. Clone the repository and navigate to the project directory:
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

### 1. Subscribe to an Address
- **Endpoint:** `POST /subscribe`
- **Request Payload:**
  ```json
  {
      "address": "0xYourEthereumAddress"
  }
  ```
- **Example cURL Command:**
  ```bash
  curl -X POST http://localhost:8080/subscribe \
       -H "Content-Type: application/json" \
       -d '{"address": "0xYourEthereumAddress"}'
  ```

### 2. Get Transactions for a Subscribed Address
- **Endpoint:** `GET /transactions?address=0xYourEthereumAddress`
- **Example cURL Command:**
  ```bash
  curl "http://localhost:8080/transactions?address=0xYourEthereumAddress"
  ```

### 3. Get the Latest Processed Block
- **Endpoint:** `GET /currentBlock`
- **Example cURL Command:**
  ```bash
  curl http://localhost:8080/currentBlock
  ```