calendar = {};

getCalendar = function(callback) {
    pg.connect(CONN_STR, function(error, client) {
        client.query("SELECT value->'periods' FROM config WHERE id = 'calendar';",
            function(error, result) {
                if(error) {
                    callback(error);
                }
                var periods = {};
                for (var key in result.rows[0]) {
                    periods = result.rows[0][key];
                }
                periods.sort((elem1, elem2) => {
                    return elem2.timeInSeconds - elem1.timeInSeconds;
                });
                callback(null, periods);
            }
        )
    });
};

Meteor.methods({
    'getWCTString': function(worldSeconds) {
        var worldTime = {};
        calendar.forEach((period) => {
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
});