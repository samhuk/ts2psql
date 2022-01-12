/**
 * A user that can have some orders
 */
/* ts2psql { "tableName": "users", "toSnakeCase": true } */
export type User = {
  /**
   * The id of the user.
   */
  /* ts2psql { "primaryKey": true, "serial": true } */
  id: number
  /* ts2psql { "unique": true } */
  uuid: string
  /* ts2psql { "columnName": "username", "unique": true, "maxLength": 50 } */
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

/* ts2psql { "tableName": "orders" } */
export type Order = {
  /* ts2psql { "primaryKey": true, "serial": true } */
  id: number
  /* ts2psql { "unique": true } */
  uuid: string
  /* ts2psql { "fk": { "type": User, "property": "id" } } */
  userId: number
  /* ts2psql */
  cancelled: boolean
  /* ts2psql */
  timeCreated: Date
  /* ts2psql */
  timeCancelled?: Date
}
/* ts2psql end */

/* ts2psql { "tableName": "tire_orders" } */
export type TireOrders = {
  /* ts2psql { "primaryKey": true, "serial": true } */
  id: number
  /* ts2psql { "unique": true } */
  uuid: string
  /* ts2psql */
  tireId: number
  /* ts2psql { "fk": { "type": User, "property": "id" } } */
  userId: number
  /* ts2psql */
  cancelled: boolean
  /* ts2psql */
  timeCreated: Date
  /* ts2psql */
  timeCancelled?: Date
}
/* ts2psql end */