/**
 * Created by Than on 19.07.2016.
 */
Meteor.publish('allPersons', function () {
    return liveDb.select('SELECT id, name FROM persons');
});