/**
 * Created by Than on 31.07.2016.
 */
Template.personList.helpers({
    persons: function() {
        return Persons.reactive();
    }
});