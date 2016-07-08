
exports.up = function(knex, Promise) {
    return Promise.all([
        knex.schema.createTable("test_time", function(table){
            table.increments();
            table.integer("world_time").notNullable();
            table.integer("time_speed").notNullable();
        })
    ]);
};

exports.down = function(knex, Promise) {
    return Promise.all([
        knex.schema.dropTable("test_time")
    ])
};
