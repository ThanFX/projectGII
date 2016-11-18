/**
 * Created by Than on 20.07.2016.
 */
Meteor.publish('personsCh', function() {
    return liveDb.select('SELECT hunger, thirst, fatigue, somnolency, state FROM person_health_characteristic');
});