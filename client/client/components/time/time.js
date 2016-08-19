/**
 * Created by Than on 31.07.2016.
 */

var wct = {};

function getWCTString(worldSeconds) {
    var worldTime = {};
    Calendar.sort((elem1, elem2) => {
        return elem2.timeInSeconds - elem1.timeInSeconds;
    });
    Calendar.forEach((period) => {
        var t = Math.floor(worldSeconds / period.timeInSeconds) + period.minValue;
        if ((period.periodLabel == 'minute' || period.periodLabel == 'hour') && t < 10) {
            worldTime[period.periodLabel] = '0' + t;
        } else {
            worldTime[period.periodLabel] = t;
        }
        worldSeconds -= ((t - period.minValue) * period.timeInSeconds);
    });
    return worldTime;
}

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
        wct.calendar = getWCTString(worldTime);
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