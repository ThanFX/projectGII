/*
Meteor.publish('newWorldTime', function() {
    let sql = `SELECT 1 AS _id, world_time FROM time WHERE id = 1`;

    // Standard method, works but signals subscription ready too soon
    // return mpg.select('players', sql, function(trig) { return true });

    // Alternative method, produces less flicker on the initial resultset
    mpg.live_select(this, 'time', sql, function(trig) { return true });
});
*/

Meteor.publish('worldTime', function () {
    return liveDb.select('SELECT world_time FROM time WHERE id = 1');
});
