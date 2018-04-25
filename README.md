# database-migrator

[![Build Status](https://travis-ci.org/pengsrc/database-migrator.svg?branch=master)](https://travis-ci.org/pengsrc/database-migrator)
[![Go Report Card](https://goreportcard.com/badge/github.com/pengsrc/database-migrator)](https://goreportcard.com/report/github.com/pengsrc/database-migrator)
[![License](http://img.shields.io/badge/license-apache%20v2-blue.svg)](https://github.com/pengsrc/database-migrator/blob/master/LICENSE)

Command-line management tools for managing database schema, inspired by Rails Active Record, works like `rake db:migrate`.

Read [_Active Record Migrations_](http://guides.rubyonrails.org/active_record_migrations.html) in [Rails Guide](http://guides.rubyonrails.org) to know the principle of database migration. The difference have to be aware is that `database-migrator` takes SQL files directly instead of model definitions in Rails.

## Installation

### Install from Source Code

``` bash
$ git clone git@github.com:pengsrc/database-migrator.git
$ cd database-migrator
$ glide install
...
[INFO]	Replacing existing vendor dependencies
$ make install
...
Installing into /data/go/bin/database-migrator...
Done
```

### Download Precompiled Binary

1. Go to [releases tab](https://github.com/pengsrc/database-migrator/releases) and download the binary for your operating system.
2. Unarchive the downloaded file, and put the executable file `database-migrator` into a directory that in the `$PATH` environment variable, for example `/usr/local/bin`.
3. Run `database-migrator --help` to get started.

## Usage

The `database-migrator` contains some sub commands, each of them support `--help` flag to print out all supported parameters.

``` Bash
$ database-migrator --help
Database schema migration tool

Usage:
  database-migrator [flags]
  database-migrator [command]

Available Commands:
  down        Revert one migration
  help        Help about any command
  status      Show database schema status
  sync        Sync database schema
  up          Run one migration

Flags:
  -h, --help   help for database-migrator

Use "database-migrator [command] --help" for more information about a command.
```

``` Bash
$ database-migrator sync --help
Sync database schema

Usage:
  database-migrator sync [flags]

Flags:
  -d, --database string     Specify database name
      --dialect string      Specify database dialect (default "MySQL")
      --help                Show help
  -h, --host string         Specify database host (default "127.0.0.1")
  -m, --migrations string   Specify migration files
  -P, --password string     Specify password
  -p, --port int            Specify database port (default 3306)
  -u, --user string         Specify user (default "root")
```

## Example

### Migration Files

Create a directories to keep migration files, and create a migration file.

``` Bash
$ mkdir migrations

$ database-migrator new -m migrations -n create_objects
New migration created at "migrations/20180425065739_create_objects.sql".

$ cat migrations/20180425065739_create_objects.sql
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied.
CREATE TABLE `examples` (
  `id` bigint(11) UNSIGNED NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
)
  ENGINE = InnoDB
  AUTO_INCREMENT = 10000
  DEFAULT CHARSET = utf8;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back.
DROP TABLE `examples`;
```

The `Up` and `Down` section are pretty much similar to the `up` and `down` method in Rails model, see [_Using the up/down Methods_](http://guides.rubyonrails.org/active_record_migrations.html#using-the-up-down-methods) to get a better understanding.

### Run

Check current status.

``` Bash
$ database-migrator status -h 127.0.0.1 -p 3306 -u root -P root -d test -m migrations
Applied At               Migration
================================================================================
Pending               -  20180425065739_create_objects
```

Sync status.

``` Bash
$ database-migrator sync -h 127.0.0.1 -p 3306 -u root -P root -d test -m migrations
Migrated to the latest database schema.
- 20180425065739_create_objects

$ database-migrator status -h 127.0.0.1 -p 3306 -u root -P root -d test -m migrations
Applied At               Migration
================================================================================
2018-04-25T07:04:27Z  -  20180425065739_create_objects
```

Apply/Revert one by one.

``` Bash
$ database-migrator down -h 127.0.0.1 -p 3306 -u root -P root -d test -m migrations
Revert migration complete, 20180425065739_create_objects.

$ database-migrator up -h 127.0.0.1 -p 3306 -u root -P root -d test -m migrations
Run migration complete, 20180425065739_create_objects.
```

## Notices

- The `database-migrator` creates a table named `schema_migrations` to track schema history.
- Currently, only MySQL direct is supported.

## LICENSE

The Apache License (Version 2.0, January 2004).
