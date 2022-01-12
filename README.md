# ts2psql

NOTE: This is WIP. Details in this readme are ideal state. Current usage: `go build main.go && ./main` (or main.exe if on Windows OS).

A script that converts Typescript type declarations into PostgreSQL CREATE TABLE scripts.

## Usage

test.ts:

```ts
/**
 * A user of the application
 */
/* ts2psql { "tableName": "users", "toSnakeCase": true } */
export type User = {
  /**
   * The id of the user.
   */
  /* ts2psql { "primaryKey": true, "serial": true } */
  id: number
  /* ts2psql { "unqiue": true } */
  uuid: string
  /* ts2psql { "unqiue": true, "maxLength": 50 } */
  name: string
  /* ts2psql { "default": false } */
  deleted: boolean
  /* ts2psql */
  timeCreated: Date
  /* ts2psql */
  timeLastLogin?: Date
  /* ts2psql */
  timeDeleted?: Date
}
/* ts2psql end */
```

```
ts2psql -f ./test.ts -o ./out.txt
```

out.txt:

```postgresql
CREATE TABLE users (
  id serial PRIMARY KEY,
  uuid VARCHAR ( 50 ) UNIQUE NOT NULL,
  name VARCHAR ( 50 ) NOT NULL,
  ...
);
```