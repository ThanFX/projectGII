/**
 * Created by Than on 31.07.2016.
 */

var wct = {};

function getWorldTime(cTime) {
    cTime.reactive();
    cTime.depend();
    var time = cTime.filter(function(t){
        return t["world_time"];
    });
    return time.length && time[0]["world_time"];
}

function getWorldCalendarTime(cTime) {
    var worldTime = getWorldTime(cTime);
    if (worldTime > 0) {
        //console.log(Session.get("calendar"));
        getData(Meteor.call, 'getWCTString', worldTime).then(
            time => {
                wct.calendar = time;
                return new Promise((resolve) => resolve());
            }
        ).catch(
            error => {
                console.log(error);
            }
        );
    }
    var s = "";
    if (wct.calendar && wct.calendar["year"] > 1) {
        s = wct.calendar["year"] + " год, " + wct.calendar["month"] + " месяц, " +
            wct.calendar["ten_day"] + " декада, " + wct.calendar["day"] + " день, " +
            wct.calendar["hour"] + ":" + wct.calendar["minute"];
    }
    return s;
}

Template.time.helpers({
    worldTime: function() {
        wct.time = getWorldTime(cTime);
        return wct.time;
    },
    worldCalendarTime: function() {
        return getWorldCalendarTime(cTime);
    }
});

/*
 newTime = new Mongo.Collection('time');

 Template.newTime.onCreated(function () {
 this.subscribe('newWorldTime');
 });

 Template.newTime.helpers({
 newTime2: function () {
 // Still need to sort client-side since record order is not preserved
 return newTime.find({_id: 1});
 },
 });
 */