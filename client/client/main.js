/**
 * Created by Than on 09.07.2016.
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



Template.timer.helpers({
    worldTime: function() {
        wct.time = getWorldTime(cTime);
        return wct.time;
    },
    worldCalendarTime: function() {
        return getWorldCalendarTime(cTime);
    }
});


Template.allPersons.helpers({
    persons: function() {
        return Persons.reactive();
    }
});


/*
// Provide a client side stub for latency compensation
Meteor.methods({
    'incScore': function(id, amount){
        var originalIndex;
        players.forEach(function(player, index){
            if(player.id === id){
                originalIndex = index;
                players[index].score += amount;
                players.changed();
            }
        });

        // Reverse changes if needed (due to resorting) on update
        players.addEventListener('update.incScoreStub', function(index, msg){
            if(originalIndex !== index){
                players[originalIndex].score -= amount;
            }
            players.removeEventListener('update.incScoreStub');
        });
    }
});

Template.leaderboard.helpers({
    players: function () {
        return players.reactive();
    },
    selectedName: function () {
        players.depend();
        var player = players.filter(function(player){
            return player.id === Session.get("selectedPlayer");
        });
        return player.length && player[0].name;
    }
});

Template.leaderboard.events({
    'click .inc': function () {
        Meteor.call('incScore', Session.get("selectedPlayer"), 5);
    }
});

Template.player.helpers({
    selected: function () {
        return Session.equals("selectedPlayer", this.id) ? "selected" : '';
    }
});

Template.player.events({
    'click': function () {
        Session.set("selectedPlayer", this.id);
    }
});
*/
