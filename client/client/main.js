/**
 * Created by Than on 09.07.2016.
 */

var wct = {};

function getWorldTime(cTime) {
    cTime.depend();
    var time = cTime.filter(function(t){
        return t.world_time;
    });
    return time.length && time[0].world_time;
}

function getWorldCalendarTime(worldTime) {
    getData(Meteor.call, 'getWCTString', worldTime).then(
        time => {
            wct = time;
            return new Promise((resolve) => resolve());
        }
    ).catch(
        error => {
            console.log(error);
        }
    );
}



Template.timer.helpers({
    times: function() {
        return cTime.reactive();
    },
    worldTime: function() {
        return getWorldTime(cTime);
    },
    worldCalendarTime: function() {
        getWorldCalendarTime(getWorldTime(cTime));
        if(wct) {
            return wct.hour + ":" + wct.minute;
        }
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
