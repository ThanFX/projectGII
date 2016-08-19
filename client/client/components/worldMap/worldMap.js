/**
 * Created by Than on 09.08.2016.
 */
mapRows = 5;
mapColumns = 5;
startX = 3;
startY = -5;
var map = new Array(mapRows);

function getChunk(i, j) {
    let c = WorldMap.filter(function(chunk){
        return (chunk.x == (i + startX) && chunk.y == (startY - j));
    });
    return c.length && c[0];
}

Template.worldMap.helpers({

});

Template.worldMap.onCreated(() => {
    for(let i = 0; i < map.length; i++){
        map[i] = new Array(mapColumns);
        for(let j = 0; j < map[0].length; j++) {
            map[i][j] = getChunk(i, j);
        }
    }
    this.$('.chunk').css({
        'backgroundImage': 'url(resources/forest.png)'
    });
    console.log(this.$('.chunk'));
});