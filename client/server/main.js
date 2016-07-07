import { Meteor } from 'meteor/meteor';

Meteor.startup(() => {
    var time = new SQL.Collection('time', 'postgres://postgres:postgres@localhost/postgres');
    /*
    if (Meteor.isServer) {
        time.publish('time', function() {
            return time.select('time.id');
        });
        console.log(time);
    }
    */
});
