import { Meteor } from 'meteor/meteor';

Meteor.startup(() => {


    /*
    Meteor.publish('cTime', function(){
        return postgres.select('SELECT world_time FROM time WHERE id = 1;');
    });
    */

    //console.log(postgres.select('SELECT * FROM time WHERE id = 1'));
    /*
    var Times = new PG.Table("time");
    var res = Times.where('id', 3).fetch();
    console.log(res);
    */
    //var time = new SQL.Collection('time', 'postgres://postgres:postgres@localhost/postgres');
    /*
    if (Meteor.isServer) {
        time.publish('time', function() {
            return time.select('time.id');
        });
        console.log(time);
    }
    */

});


times = new PgSubscription('cTime');

if (Meteor.isServer) {
    var postDB = new LivePg("postgres://postgres:postgres@localhost:5432/postgres", "postgres_example");

    var closeAndExit = function() {
        postDB.cleanup(process.exit);
    };

    process.on('SIGTERM', closeAndExit);
    process.on('SIGINT', closeAndExit);

    Meteor.publish('cTime', function() {
        return postDB.select('SELECT world_time from time WHERE id = 1');
    });
}