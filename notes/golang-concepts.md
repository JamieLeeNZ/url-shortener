# Golang Concepts Notes

These are my notes on the Go programming language, focusing on its unique features and how they differ from Java, as I build this project.

## 1. Structs

- Structs in Go are composite types that group fields together, similar to classes in Java but without methods attached directly, or structs in C.
- Example:
  ```go
  type MemoryStore struct {
      mu      sync.RWMutex
      storage map[string]string
  }
  ```

## 2. Constructor Functions

- Go does not have constructors like Java, but uses conventionally named functions (e.g., `NewMemoryStore()`) to initialize structs.
- These functions return pointers to newly allocated and initialized structs.
- The `make` function is used to initialize maps, slices, and channels.
- Some struct fields, such as `sync.RWMutex`, are valid with their zero value and do not require explicit initialization.
- Example:

  ```go
  func NewMemoryStore() *MemoryStore {
      return &MemoryStore{
          storage: make(map[string]string),
      }
  }
  ```

## 3. Pointers (`*` and `&`)

- `&` is the address-of operator: returns the memory address of a variable (pointer).
- `*` is the dereference operator: accesses or sets the value stored at a pointer’s address.
- Example:
  ```go
  x := 10
  p := &x       // p is a pointer to x
  fmt.Println(*p) // prints 10
  ```
- Modifying `x` after `p := &x` changes the value accessed by `*p`, but modifying `x` after `p := x` does not affect `p`.
- Java differs from Go because it uses implicit references for objects, does not allow pointers to primitives, and hides memory addresses entirely.

## 4. Visibility

- Go uses capitalization to determine visibility:
  - Capitalized names (e.g., `PublicFunc`) are exported and accessible outside the package.
  - Uncapitalized names (e.g., `privateFunc`) are unexported and only accessible within the package.
- These apply to functions, struct fields, variables, constans, and methods.

## 5. Functions and Methods

- Declared with `func` and not tied to any type.
- Can be defined at the top level of a package.

```go
func greet(name string) string {
    return "Hello, " + name
}
```

- Methods are also defined with `func`, but with a receiver to bind the method to a type (usually a struct).

```go
type Person struct {
    Name string
}

func (p Person) Greet() string {
    return "Hello, " + p.Name
}
```

- In summary, methods have a receiver type, so they automatically get access to the fields of the receiver type, whereas functions do not.
- In Java, all functions are methods as they are always associated with a class or object, while in Go, functions can exist independently of types.

## 6. Maps

- Maps in Go are built-in data types that store key-value pairs, similar to `HashMap` in Java.
- To declare a (nil) map:
  ```go
  var m map[string]int // nil map, not initialized
  ```
  - To initialize a map, use the `make` function:
  ```go
  m = make(map[string]int)
  ```
- Maps can be initialized with values using a map literal:
  ```go
  m := map[string]int{"key1": 1, "key2": 2}
  ```
- Maps are reference types, meaning they are passed by reference, not by value.
- To check if a key exists in a map, use the two-value assignment:
  ```go
  value, exists := m["key1"]
  if exists {
      fmt.Println("Key exists with value:", value)
  } else {
      fmt.Println("Key does not exist")
  }
  ```
- If the key does not exist, `value` will be the zero value for the map's value type (e.g., `0` for `int`, `""` for `string`).
- To add or update a key-value pair:
  ```go
  m["key3"] = 3 // Add or update key3
  ```
- To delete a key-value pair:
  ```go
  delete(m, "key1") // Deletes key1 from the map
  ```
- Iterating over a map can be done using a `for` loop (though the order is not guaranteed):
  ```go
  for key, value := range m {
      fmt.Println("Key:", key, "Value:", value)
  }
  ```

## Mutexs and RWMutexs

- Go provides **mutexes** (mutual exclusion locks) via the `sync` package to handle safe concurrent access to shared data.
- Two common types:
  - `sync.Mutex` — a standard mutex lock (exclusive lock)
  - `sync.RWMutex` — read-write mutex allowing multiple readers or one writer
- Example:

  ```go
  type MemoryStore struct {
      mu   sync.RWMutex
      data map[string]string
  }

  func (s *MemoryStore) Set(key, value string) {
      s.mu.Lock()         // Acquire exclusive write lock
      defer s.mu.Unlock() // Ensure unlock after function exits
      s.data[key] = value
  }

  func (s *MemoryStore) Get(key string) (string, bool) {
      s.mu.RLock()          // Acquire shared read lock
      defer s.mu.RUnlock()  // Ensure unlock after function exits
      value, ok := s.data[key]
      return value, ok
  }
  ```

- `RWMutex` improves performance by allowing concurrent reads while still protecting writes.

## 7. Handlers

- Handlers in Go are functions that process HTTP requests and write HTTP responses.
- Handlers have the following signature:
  ```go
  func(w http.ResponseWriter, r *http.Request)
  ```
  where `http.ResponseWriter` is used to construct the HTTP response, and `*http.Request` contains the incoming request data.
- Handlers can be registered with the HTTP server as:

  - Standalone functions:
    ```go
    http.HandleFunc("/path", handlerFunction)
    ```
  - Methods of a struct:
    ```go
    http.HandleFunc("/path", s.MethodName)
    ```

- Inside handlers:
  - Parse and validate request data (e.g., decode JSON body with json.NewDecoder(r.Body).Decode(&struct))
  - Validate input (e.g., check required fields, URL format)
  - Perform business logic (e.g., check in-memory store, generate short keys)
  - Write responses with status codes and JSON encoding (w.Header().Set, json.NewEncoder(w).Encode(...))

## 8. HTTP Methods

- Go's `net/http` package provides built-in support for handling HTTP methods like GET, POST, PUT, DELETE, etc.
- Along with handlers, we can use:
  - `http.Redirect` to redirect requests
  - `http.Error` to send error responses with specific status codes
  - `http.ListenAndServe` to start the HTTP server
  - `http.ServeFile` to serve static files
  - `http.MethodGET`, `http.MethodPost`, etc., to check request methods in handlers
  - `http.StatusOK`, `http.StatusNotFound`, etc., for standard HTTP status codes
  - `http.SetCookie` to set cookies in responses

## 9. Interfaces

- Interfaces in Go define a set of methods that a type must implement, similar to Java interfaces.
- Example:
  ```go
  type Store interface {
      Set(key, value string)
      Get(key string) (string, bool)
  }
  ```
- In Go, we don't need to explicitly declare that a type implements an interface; it is implicit. If a type has all the methods defined in an interface, it implements that interface.
- IMPORTANT: Interfaces are reference types, meaning they hold a pointer to the underlying data and the type information. We don't need to pass or use pointers (`*MyInterface`) to structs to implement interfaces, as Go automatically handles this (just use `MyInterface`).
- Interfaces let us swap out implementations easily, allowing for flexible code design and testing.
- To enforce that a struct implements an interface at compile time, add this line (usually near the top of the file where the struct is defined):
  ```go
  var _ Store = (*MemoryStore)(nil) // Ensures MemoryStore implements Store interface
  ```

## 10. Integrating PostgreSQL

- Use the `pgxpool` package for PostgreSQL connection pooling.
- Use `godotenv` to load environment variables from a `.env` file like `DATABASE_URL`.
- Use `pgxpool` to create and manage a pool of DB connections:
  ```go
  func NewPostgresStore(connString string) (*PostgresStore, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
      return nil, err
    }
    pool, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
      return nil, err
    }
    if err := pool.Ping(context.Background()); err != nil {
      pool.Close()
      return nil, err
    }
    return &PostgresStore{db: pool}, nil
  }
  ```
- A connection pool is a cache of database connections that can be reused, improving performance by avoiding the overhead of establishing new connections for each request.
- `Context` is used with all DB calls for managing request lifetimes, timeouts, and cancellations.
- CRUD operations can be implemented using `pgxpool` methods like `Query`, `Exec`, and `QueryRow`.
- For example, to insert a new record:
  ```go
  func (s *PostgresStore) Set(key, originalURL string) error {
  _, err := s.db.Exec(context.Background(), `
    INSERT INTO url_mappings (key, original_url) VALUES ($1, $2)
    ON CONFLICT (key) DO UPDATE SET original_url = EXCLUDED.original_url
  `, key, originalURL)
  return err
  }
  ```
- Use placeholders (`$1`, `$2`, etc.) to safely inject parameters and avoid SQL injection.
- Always close the pool with defer pool.Close() when shutting down the app.
- The pool can be configured with options like maximum connections, connection timeouts, and idle timeouts to optimize performance based on the app's needs. For example:

  ```go
  config.MaxConns = 20               // Maximum open connections in the pool
  config.MinConns = 5                // Minimum open connections to keep alive
  config.MaxConnIdleTime = 5 * time.Minute  // Close connections idle longer than this
  config.MaxConnLifetime = 30 * time.Minute // Recycle connections after this duration
  ```
