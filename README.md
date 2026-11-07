# Gator

Gator is a command-line RSS feed aggregator written in Go. It stores feeds and
posts in PostgreSQL and lets users follow feeds and browse their posts.

## Requirements

- Go 1.27.1 or newer
- PostgreSQL
- [Goose](https://github.com/pressly/goose) to apply the database migrations

Install Goose with Go:

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Install Gator

Install the CLI from the public repository:

```sh
go install github.com/MaxKQC/Gator@latest
```

Make sure the Go bin directory is on your `PATH`. It is usually `$(go env GOPATH)/bin`.

The Go command's [compile and install documentation](https://pkg.go.dev/cmd/go#hdr-Compile_and_install_packages_and_dependencies)
has more information about `go install`.

## Configure PostgreSQL

Start PostgreSQL and create a database, for example:

```sh
createdb gator
```

Create `~/.gatorconfig.json` with a PostgreSQL connection URL and an empty
current user name. Replace the example credentials with your local PostgreSQL
settings:

```json
{
  "db_url": "postgres://postgres:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

From the repository root, apply the schema migrations. Use the same database
URL as in the config file:

```sh
goose -dir sql/schema postgres "postgres://postgres:password@localhost:5432/gator?sslmode=disable" up
```

## Run Gator

Register a user to set the current user, then add a feed. Adding a feed also
starts following it automatically:

```sh
gator register alice
gator addfeed "The Boot.dev Blog" https://blog.boot.dev/index.xml
```

Some available commands:

| Command | Description |
| --- | --- |
| `gator login <name>` | Switch to an existing user |
| `gator users` | List users and show the current user |
| `gator addfeed <name> <url>` | Add and follow a feed |
| `gator feeds` | List feeds and their creators |
| `gator follow <url>` | Follow an existing feed |
| `gator following` | List feeds followed by the current user |
| `gator unfollow <url>` | Unfollow a feed |
| `gator browse [limit]` | Browse recent posts; defaults to 2 |
| `gator agg <duration>` | Fetch one feed at a time on a continuous interval, such as `1m` |
| `gator reset` | Delete all users and their dependent feeds, follows, and posts |

Run `gator agg 1m` in a separate terminal to keep collecting posts. Stop the
aggregator with Ctrl+C.
