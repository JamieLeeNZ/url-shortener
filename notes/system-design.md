# System Design Notes

These are my notes on System Design concepts encountered while building this project. This includes notes on databases, caching, scalable systems, etc.

## 1. PostgreSQL

- PostgreSQL is a powerful, open-source relational database management system (RDBMS).
- It has robust SQL support, strong ACID compliance, and extensive features like JSONB support, full-text search, and advanced indexing.
- In cloud environments, it is often used as a managed service (e.g., AWS RDS, Google Cloud SQL).

#### System design:

PostgreSQL typically acts as the **source of truth**. Common system patterns include:

- Backend APIs: PostgreSQL handles user data, relationships, and transactional operations.
- Caching: Frequently accessed data can be cached in memory (e.g., Redis) for fast access.
- Microservices: Each service can have its own PostgreSQL instance or share a common one, depending on the architecture.
- Sharding / Partitioning: For large datasets, PostgreSQL can be partitioned to improve performance and manageability.
- Read Replicas: Read replicas can be used to offload read operations from the primary database without impacting write performance, improving performance and scalability.

#### PostgreSQL vs. MySQL vs. SQLite vs. MongoDB:

PostgreSQL

- Best all-rounder for modern backend systems with structured data, advanced SQL, and cloud scalability.
- Use it when you need strong data integrity, complex querying, or future-proof cloud deployment.

MySQL

- Reliable and fast for simpler web apps, especially in traditional LAMP stacks.
- Use it when you want a lightweight relational database with broad hosting support and team familiarity.

SQLite

- Lightweight and embedded, perfect for mobile apps, prototyping, or single-user tools.
- Use it when your app is local, doesn’t need concurrency, or must run with zero setup.

MongoDB

- Schema-less and flexible, ideal for fast-changing data models and JSON-heavy APIs.
- Use it when your data is unstructured, evolving rapidly, or you need horizontal scalability and fast iteration.
