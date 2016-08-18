/**
 * Created by Than on 18.08.2016.
 */
Meteor.publish('worldMap', function () {
    return liveDb.select('SELECT * FROM world_map');
});