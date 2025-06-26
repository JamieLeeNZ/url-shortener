# Development Workflow Notes

These are my notes on new development workflow concepts I encountered and learnt while building this project. This includes notes on version control, CI/CD, testing, and deployment.

## 1. SQL Migrations

- SQL migrations are version-controlled, incremental changes to the database schema written as SQL scripts.
- Each migration file describes how to apply (Up) and undo (Down) a discrete schema change.

### Why Use SQL Migrations?

- Track schema changes alongside application code.
- Allow rollback of schema changes if needed.
- Support automation in CI/CD pipelines.

### Migration File Structure

- Migration files are named with a version prefix and description, for example:
  ```
  V01__create_url_mappings.sql
  V02__add_user_id_to_url_mappings.sql
  ```
- The files are applied in lexicographical order based on the version prefix so use zero-padded numbers for consistency (e.g., `V01`, `V02`, etc.).
- An example migration file might look like this:

  ```sql
  -- +migrate Up
  CREATE TABLE url_mappings (
    key TEXT PRIMARY KEY,
    original_url TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
  );

  -- +migrate Down
  DROP TABLE url_mappings;
  ```

### Migration Workflow

1. Write migration files describing incremental schema changes.
2. Commit migration files to repository.
3. Use a migration tool such as sql-migrate to apply migrations.
4. Migrations are applied in order, tracked in a database table to avoid reapplying.
5. Use environment variables (e.g., DATABASE_URL) to securely configure database connections.
6. Automate migrations with scripts.
