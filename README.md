# ts2psql

NOTE: This is WIP. Details in this readme are ideal state. Current usage: `go build && ./ts2psql` (or `go build && ts2psql` if on Windows OS).

A script that converts Typescript type declarations into PostgreSQL CREATE TABLE scripts.

This is mainly useful for saving time when statements have to be created for a large number of types in your various projects.

## Usage

test.ts:

```ts
/* ts2psql { "tableName": "users" } */
export type User = {
  /* ts2psql { "primaryKey": true, "serial": true } */
  id: number
  /* ts2psql { "unique": true } */
  uuid: string
  /* ts2psql { "unique": true, "maxLength": 50 } */
  name: string
}
/* ts2psql end */
```

Add `ts2psqlconfig.json` configuration file:

```
{
  "include": ["*.ts"],
  "outFile": "./out.sql",
}
```

Run `ts2psql`

out.sql:

```postgresql
CREATE TABLE users (
  id serial PRIMARY KEY,
  uuid VARCHAR(50) UNIQUE NOT NULL,
  name VARCHAR(50) NOT NULL
);
```
