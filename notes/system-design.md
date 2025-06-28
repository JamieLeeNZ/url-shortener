# System Design Notes

These are my notes on System Design concepts encountered while building this project. This includes notes on databases, caching, scalable systems, etc.

## 1. PostgreSQL

- PostgreSQL is a powerful, open-source relational database management system (RDBMS).
- It has robust SQL support, strong ACID compliance, and extensive features like JSONB support, full-text search, and advanced indexing.
- In cloud environments, it is often used as a managed service (e.g., AWS RDS, Google Cloud SQL).

#### Use Cases:

PostgreSQL typically acts as the **source of truth**. Common system patterns include:

- Backend APIs: PostgreSQL handles user data, relationships, and transactional operations.
- Caching: Frequently accessed data can be cached in memory (e.g., Redis) for fast access.
- Microservices: Each service can have its own PostgreSQL instance or share a common one, depending on the architecture.
- Sharding / Partitioning: For large datasets, PostgreSQL can be partitioned to improve performance and manageability.
- Read Replicas: Read replicas can be used to offload read operations from the primary database without impacting write performance, improving performance and scalability.

#### PostgreSQL vs. MySQL vs. SQLite vs. MongoDB:

PostgreSQL:

- Best all-rounder for modern backend systems with structured data, advanced SQL, and cloud scalability.
- Use it when you need strong data integrity, complex querying, or future-proof cloud deployment.

MySQL:

- Reliable and fast for simpler web apps, especially in traditional LAMP stacks.
- Use it when you want a lightweight relational database with broad hosting support and team familiarity.

SQLite:

- Lightweight and embedded, perfect for mobile apps, prototyping, or single-user tools.
- Use it when your app is local, doesn’t need concurrency, or must run with zero setup.

MongoDB:

- Schema-less and flexible, ideal for fast-changing data models and JSON-heavy APIs.
- Use it when your data is unstructured, evolving rapidly, or you need horizontal scalability and fast iteration.

## 2. Redis

- Redis is an in-memory data structure store, often used as a database, cache, and message broker.
- It supports various data structures like strings, hashes, lists, sets, and sorted sets.
- Redis is fast, supports transactions, and is highly scalable.

#### Use Cases:

- Redis is often used as a cache to reduce database load and improve response times.
- In horizontally scaled apps (multiple instances or containers), each instance has isolated memory, so Redis acts as a shared cache that all instances can access.
- This allows shared, consistent, and fast access to data regardless of where the request is handled.

#### Redis vs. PostresSQL:

Redis:

- In-memory (RAM), optionally persistent
- Data stored as data structures (strings, hashes, lists, etc.)
- Extremely fast for read/write operations
- Ideal for caching, real-time analytics, and pub/sub messaging

PostgreSQL:

- Disk-based, ACID-compliant
- Data stored in tables with SQL support
- Slower than Redis for simple key-value access
- Ideal for structured data, complex queries, and transactional operations
- Robust recovery and durability features

Redis is NOT a replacement for a persistent database like PostgreSQL. It is typically used alongside it to cache frequently accessed data, reducing load on the primary database and improving performance.

#### Caching Strategies:

Write-to-Both:

- Write data to both Redis and PostgreSQL on create/update.
- Keeps DB and cache in sync.
- Slightly slower writes due to two operations.
- Prevents cache misses on reads.

Read-Through:

- On read, check Redis first.
- If cache miss, read from PostgreSQL, store in Redis, and return.
- Fast reads, but initial read may be slower due to cache miss.
- Optimises reads without preloading everything.

## 3. Session Management

A session is a temporary stateful connection between a client (browser) and the server, maintained across multiple HTTP requests after the user logs in. Since HTTP is stateless, session management provides a way to "remember" the authenticated user between requests

### Session Management Workflow:

1. **User Logs In**: User submits credentials (username/password) to the server.
2. **Server Validates Credentials**: The server checks the credentials against the database.
3. **Session Creation**: If valid, the server creates a session:
   - Generates a unique session ID (e.g., UUID).
   - Stores session data (user ID, expiration time, etc.) in a session store (e.g., Redis).
   - Sends the session ID to the client as a cookie.
4. **Client Stores Session ID**: The client stores the session ID in a cookie.
5. **Subsequent Requests**: For each request, the client sends the session ID cookie.
6. **Server Validates Session**: The server retrieves the session data using the session
7. **Session Expiration**: Sessions have an expiration time. After this time, the session is considered invalid (24 hours is common).
8. **User Logs Out**: The user can log out, which invalidates the session

### Cookies

- Cookies store the session_id on the client side so it can be sent automatically with every HTTP request.
- They are set by the server in the HTTP response headers and sent back by the client in subsequent requests:

  `Set-Cookie: session_id=abc123; HttpOnly; Path=/; Expires=...`

- On every request, the browser includes the cookie:

  `Cookie: session_id=abc123`

- Security Flags:
  - HttpOnly: Prevents JavaScript access (XSS protection).
  - Secure: Ensures the cookie is only sent over HTTPS.
  - SameSite: Prevents CSRF (e.g., SameSite=Lax or Strict).
- Expiration: Controlled by Expires or Max-Age attribute. Matches the session TTL (e.g., 24h).

## 4. JWT Authentication

JSON Web Tokens (JWT) are a compact, URL-safe means of representing claims to be transferred between two parties. They are widely used for stateless authentication in web and mobile applications.

- A JWT is a self-contained token that includes encoded JSON data (claims) about the user or session.
- It consists of three parts separated by dots: `header.payload.signature`:
  - **Header**: Contains metadata about the token, such as the signing algorithm used (e.g., HS256).
  - **Payload**: Contains the claims, which can include user information, roles, and expiration time.
  - **Signature**: Used to verify that the sender of the JWT is who it claims to be and to ensure that the message wasn't changed along the way.

### JWT Authentication Workflow:

1. **User Logs In**: User submits credentials (username/password) to the server.
2. **Server Validates Credentials**: The server checks the credentials against the database.
3. **JWT Creation**: If valid, the server creates a JWT:
   - Encodes user information and claims in the payload.
   - Signs the token using a secret key or private key.
4. **JWT Sent to Client**: The server sends the JWT back to the client, typically in the response body or as an HTTP-only cookie.
5. **Client Stores JWT**: The client stores the JWT in its local storage or as an HTTP-only cookie.
6. **Subsequent Requests**: For each request, the client sends the JWT in the `Authorization` header:
   `Authorization: Bearer <token>`
7. **Server Validates JWT**: The server verifies the JWT signature and checks claims (e.g., expiration time, issuer).
8. **Access Granted**: If valid, the server processes the request and grants access to protected resources.
9. **Token Expiration**: JWTs typically have an expiration time (e.g., 1 hour). After expiration, the client must reauthenticate to obtain a new token.
10. **Token Refresh**: Optionally, a refresh token can be used to obtain a new JWT without requiring the user to log in again. The refresh token is usually stored securely and has a longer expiration time than the access token.

### Redis Sessions vs. JWT Authentication:

| Feature     | Redis Authentication                                                                              | JWT Authentication                                                                      |
| ----------- | ------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| Storage     | Session data stored in Redis (server-side)                                                        | Token stored on client (stateless)                                                      |
| Scalability | Requires shared or replicated Redis                                                               | Easy to scale horizontally since no server-side storage needed                          |
| Performance | Fast access to session data, but requires Redis connection and Redis query for every request      | Fast verification (locally), no database access needed for each request                 |
| Expiration  | Managed by Redis (TTL)                                                                            | Managed by JWT claims (e.g., `exp` field)                                               |
| Token Size  | Small cookie with session ID                                                                      | Larger token carrying user claims                                                       |
| Security    | Session ID is opaque and stored server-side                                                       | JWT payload is visible (not encrypted by default), must verify signatures               |
| Complexity  | More complex setup with Redis and session management                                              | Simpler token-based authentication, no session store required                           |
| Revocation  | Can invalidate sessions by deleting from Redis                                                    | Harder to revoke; requires token blacklist or short expiration                          |
| Use Cases   | Best for applications with frequent user interactions and need for server-side session management | Best for stateless APIs, mobile apps, and microservices where scalability is a priority |

Redis sessions are suitable when easy and immediate session revocation is required, session state needs to be stored on the server, maintaining server-side state is acceptable, and tighter control over security is desired.

JWT is appropriate when stateless authentication is preferred, scalability across multiple services is needed, token revocation before expiry is not a primary concern, and security relies heavily on token signing and client-side handling.

## 5. Middleware

Middleware is a layer in the web server architecture that sits between incoming HTTP requests and the final request handlers. It acts as a gatekeeper (by protecting routes), intercepting requests, performing tasks like authentication, logging, and error handling, and then passing the request to the appropriate handler.

### Authentication Middleware:

Authentication middleware is a type of middleware that handles user authentication.

- **Check Authentication**: It checks if the user is authenticated by verifying the session or JWT token.
- **Protect Routes**: It protects specific routes by allowing access only to authenticated users.
- **Redirect or Respond**: If the user is not authenticated, it can redirect them to the login page or respond with an error (e.g., 401 Unauthorized).
- **Inject User Context**: If authenticated, it can inject user information into the request context for use by subsequent handlers.

### Example Middleware Implementation:

```go
func AuthMiddleware(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Check if user is authenticated (e.g., check session or JWT)
    user, err := GetAuthenticatedUser(r)
    if err != nil {
      // If not authenticated, respond with 401 Unauthorized
      http.Error(w, "Unauthorized", http.StatusUnauthorized)
      return
    }

    // Inject user into request context
    ctx := context.WithValue(r.Context(), "user", user)
    r = r.WithContext(ctx)

    // Call the next handler
    next.ServeHTTP(w, r)
  })
}
```
