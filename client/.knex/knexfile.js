// Update with your config settings.

module.exports = {

  development: {
    client: 'postgresql',
    connection: {
      host: 	'localhost',
      port:	'5432',
      database: 'postgres',
      user:     'postgres',
      password: 'postgres'
    }
  },

  staging: {
    client: 'postgresql',
    connection: {
      host: 	'localhost',
      port:	'5432',
      database: 'postgres',
      user:     'postgres',
      password: 'postgres'
    },
    pool: {
      min: 2,
      max: 10
    },
    migrations: {
      tableName: 'knex_migrations'
    }
  },

  production: {
    client: 'postgresql',
    connection: {
      host: 	'localhost',
      port:	'5432',
      database: 'postgres',
      user:     'postgres',
      password: 'postgres'
    },
    pool: {
      min: 2,
      max: 10
    },
    migrations: {
      tableName: 'knex_migrations'
    }
  }

};
