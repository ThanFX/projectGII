/**
 * Created by Than on 20.07.2016.
 */
Meteor.publish('personsCh', function (person_id) {
    return liveDb.select(
        'SELECT hunger, thirst, fatigue, somnolency, state FROM person_health_characteristic WHERE person_id = $1', [ person_id ],
        {
            'person_health_characteristic': function(person_id) {
                return person_id > 0;
            }
        });
});