# postgres-query

This Postgres step takes a query that you define, runs it against the databse
of your choice, and serializes the results as an array of JSON objects to a
well-known output.

## Requirements

You have to setup your Postgres server such that there's a user account that
can connect externally. This obviously includes the relevant firewall rules to
make it possible for Relay infrastructure to connect as well.

## Specification

| Setting      | Data type        | Description                           | Default | Required |
|--------------|------------------|---------------------------------------|---------|----------|
| `connection` | Relay Connection | Connection to Postgres defined by URL | None    | True     |
| `statement`  | string           | The query to run.                     | None    | True     |

## Outputs

| Key       | Data type           | Description                                                                                       | 
|-----------|---------------------|---------------------------------------------------------------------------------------------------|
| `results` | string (JSON array) | The results of a successful query represented as a JSON array of objects, serialized to a string. |

## Example  

```yaml
steps:
# ...
- name: postgres-query
  image: relaysh/postgres-query
  spec:
    connection:
      url: postgres://my-postgres-host.my-domain.com:5432/my_database?sslmode=disable
    statement: "SELECT name, email_address, phone_number FROM users WHERE created_at BETWEEN '2020-01-01' AND '2020-06-01'"
```
