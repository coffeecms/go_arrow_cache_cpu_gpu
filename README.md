Here's the updated `README.md` file that includes detailed instructions for setting and getting keys with CPU and GPU modes, along with instructions on how to run the benchmarks:

---

# High-Performance Cache Management System with Apache Arrow

This project implements a high-performance key-value cache system in Go, utilizing **Apache Arrow** for efficient memory handling and supporting large-scale connections. The system supports caching data with a **Time To Live (TTL)** mechanism to automatically expire stale items. Additionally, the system can operate in **CPU** and **GPU** modes for setting and getting cache entries.

## Features

- **Key-Value Cache**: Stores and retrieves key-value pairs in memory.
- **Apache Arrow Integration**: Uses Apache Arrow to optimize memory usage for binary data storage.
- **TTL Support**: Each cache entry expires after a user-defined period (Time To Live).
- **Concurrency**: Efficient thread-safe management of concurrent cache access.
- **Mode Switching**: Supports both CPU and GPU modes for cache operations.

## Prerequisites

- Go 1.18 or higher
- Apache Arrow Go library (`github.com/apache/arrow/go/v14`)

  Install using:
  ```bash
  go get github.com/apache/arrow/go/v14
  ```

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

## Usage

### 1. Initializing the Cache

Create a new cache instance for either CPU or GPU mode using the `NewCache` function:

```go
// CPU mode
cacheCPU := NewCache(CPU)

// GPU mode
cacheGPU := NewCache(GPU)
```

### 2. Setting a Key-Value Pair with TTL

To store a key-value pair in the cache with a **Time To Live (TTL)**, use the `Set` method. The method automatically uses the appropriate mode (CPU or GPU) based on the cache instance.

#### CPU Mode
```go
cacheCPU.Set("example_key", []byte("This is some data"), 1*time.Minute)
```

#### GPU Mode
```go
cacheGPU.Set("example_key", []byte("This is some data"), 1*time.Minute)
```

- **Key**: A string representing the key for the cache entry.
- **Value**: A byte slice representing the data you want to store.
- **TTL (Time To Live)**: The duration after which the cache entry will expire. In this example, it's set to `1*time.Minute`, but you can adjust it as needed.

### 3. Retrieving a Value from the Cache

To retrieve a cached value by its key, use the `Get` method. It returns the value and a boolean indicating whether the key was found and is still valid (not expired).

#### CPU Mode
```go
value, found := cacheCPU.Get("example_key")
if found {
    fmt.Printf("Value: %s\n", string(value))
} else {
    fmt.Println("Key not found or expired")
}
```

#### GPU Mode
```go
value, found := cacheGPU.Get("example_key")
if found {
    fmt.Printf("Value: %s\n", string(value))
} else {
    fmt.Println("Key not found or expired")
}
```

- **Key**: The key for which you want to retrieve the value.
- **Return Values**: 
  - **Value**: The byte slice stored under the key (if found).
  - **Found**: A boolean (`true` if the key exists and hasn't expired, `false` otherwise).

### 4. Cleaning Up Expired Items

The cache automatically handles TTL, but you can manually clean up expired items using the `CleanExpiredItems` method:

```go
cacheCPU.CleanExpiredItems()
cacheGPU.CleanExpiredItems()
```

This function iterates over all items and removes those that have expired.

## Benchmarking

To benchmark the performance of setting and getting 1,000,000 keys in both CPU and GPU modes, the following commands will run benchmarks for both operations.

### Running Benchmarks

1. **Run the benchmark**:

   ```bash
   go run main.go
   ```

   This will execute the program and display the time taken for setting and getting 1,000,000 keys in both CPU and GPU modes.

### Example Output

```bash
Benchmarking CPU Mode
Time taken to Set 1000000 keys in CPU mode: 2m45s
Time taken to Get 1000000 keys in CPU mode: 1m20s

Benchmarking GPU Mode
Time taken to Set 1000000 keys in GPU mode: 2m30s
Time taken to Get 1000000 keys in GPU mode: 1m15s
```

## Structure

- `CacheItem`: Represents a single cache entry with the data and expiration time.
- `Cache`: Manages cache storage, including setting and getting key-value pairs, with TTL.
- `Apache Arrow`: Used to store binary data efficiently, improving memory management when handling large datasets.
- `Mode Switching`: Handles CPU and GPU modes for cache operations.

## How It Works

1. **Setting Cache Items**: 
   - Data is stored in memory using **Apache Arrow**'s `BinaryBuilder` for CPU mode and a mocked GPU implementation for GPU mode.
   - Each item is associated with a TTL (Time To Live), ensuring that stale data is automatically removed.

2. **Retrieving Cache Items**:
   - The cache checks if the key exists and whether the TTL has expired before returning the value.

3. **Concurrency**:
   - The cache system uses `sync.RWMutex` to handle concurrent reads and writes safely.

## License

This project is licensed under the MIT License. See the `LICENSE` file for more details.

---

This `README.md` provides detailed instructions for setting and getting keys with both CPU and GPU modes, along with guidance on running and interpreting benchmarks. Let me know if you need further adjustments or additional information!