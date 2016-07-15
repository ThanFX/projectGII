Meteor.publish('worldTime', function () {
    return liveDb.select('SELECT world_time FROM time WHERE id = 1');
});