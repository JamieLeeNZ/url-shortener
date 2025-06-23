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
