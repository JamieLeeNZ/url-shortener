# Issues Encountered

This document captures any issues I encountered while working on the project, along with their solutions or workarounds. This is useful for future reference and helps in debugging similar issues later.

## 1. PostgreSQL Prepared Statements Error

### Issue

When running the application, I encountered an error related to PostgreSQL prepared statements in the Set method after I added the redis cache layer. This error, however, occured in the PostgreSQL store. The error message was:

`prepared statement "stmtcache_..." already exists`

ChatGPT suggested that this could be due to the way prepared statements are cached in PostgreSQL. The error occurs when a prepared statement with the same name already exists in the session. It suggested me to:

1. Use a mutex on the Set method to ensure that only one goroutine can execute it at a time. However, this is a temporary workaround as it would block other goroutines and limit concurrency.
2. Use PreferSimpleProtocol: true, but it only works with pgx.ConnConfig (not supported directly in pgxpool anymore as of v5).
3. Use pgx.ConnConfig.StatementCacheCap = 0 (which also didn't work).
4. Inject context as a method parameter instead of using context.Background() to ensure that the context is not shared across multiple goroutines (also didn't work).

### Solution

The working solution involved disabling the statement cache and explicitly instructing pgx to use the simple query protocol for all queries:

```go
  config.ConnConfig.StatementCacheCapacity = 0
  config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
```

- `StatementCacheCapacity = 0` disables pgx’s internal prepared statement cache. This prevents pgx from automatically preparing and caching statements on the PostgreSQL server.
- However, simply disabling the cache is not enough. By default, pgx attempts to use cached prepared statements to execute queries (`QueryExecModeCacheStatement`). When the cache is disabled but the execution mode is left unchanged, pgx tries to cache statements but fails, leading to the error:
  `cannot use QueryExecModeCacheStatement with disabled statement cache`
- To resolve this, `DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol` tells pgx to always use the simple protocol for executing queries. The simple protocol sends queries directly to the server without preparing them. This means:
  - No prepared statement names are generated or reused.
  - There are no concurrency conflicts related to prepared statements.
  - The server simply executes each query as it receives it.

### Benefits & Trade-offs

- **Benefits**:
  - Eliminates the risk of prepared statement conflicts in concurrent environments.
  - Simplifies the execution model by avoiding prepared statements altogether.
  - Allows for concurrent execution of queries without locking or mutexes.
- **Trade-offs**:
  - Slightly reduced performance for repeated queries, as they are not prepared and cached on the server.
  - Loss of some benefits of prepared statements, such as query plan caching and reduced parsing overhead for repeated queries.
