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
