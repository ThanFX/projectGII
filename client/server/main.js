/*
myScore.addEventListener('updated', function(diff, data){
    data.length && console.log(data[0].score);
});

Meteor.publish('allPlayers', function(){
    // No triggers specified, the package will automatically refresh the
    // query on any change to the dependent tables (just players in this case).
    return liveDb.select('SELECT * FROM players ORDER BY score DESC');
});

Meteor.publish('playerScore', function(name){
    // Parameter array used and a manually specified trigger to only refresh
    // the result set when the row changing on the players table matches the
    // name argument passed to the publish function.
    return liveDb.select(
        'SELECT id, score FROM players WHERE name = $1', [ name ],
        {
            'players': function(row) {
                return row.name === name;
            }
        }
    );
});

Meteor.methods({
    'incScore': function(id, amount){
        // Ensure arguments validate
        //check(id, Number);
        //check(amount, Number);

        // Obtain a client from the pool
        pg.connect(CONN_STR, function(error, client, done) {
            if(error) throw error;

            // Perform query
            client.query(
                'UPDATE players SET score = score + $1 WHERE id = $2',
                [ amount, id ],
                function(error, result) {
                    // Release client back into pool
                    done();

                    if(error) throw error;
                }
            )
        });
    }
});
*/