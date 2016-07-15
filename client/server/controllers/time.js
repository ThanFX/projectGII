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
            if (period.timeInSeconds > worldSeconds) {
                worldTime[period.periodLabel] = 0 + period.minValue;
            } else {
                worldTime[period.periodLabel] = Math.floor(worldSeconds / period.timeInSeconds);
                worldSeconds -= (worldTime[period.periodLabel] * period.timeInSeconds);
                worldTime[period.periodLabel] += +period.minValue;
            }
        });
        return worldTime;
    }
});