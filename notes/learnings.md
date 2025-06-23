# Notes

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
