/**
 * Created by Than on 09.08.2016.
 */
mapRows = 5;
mapColumns = 5;
startX = 3;
startY = -5;
var map = new Array(mapRows);

function getChunk(x, y) {
    let c = WorldMap.filter(function(chunk){
        return chunk.x == x && chunk.y == y;
    });
    return c.length && c[0];
}

function getMainTerrain(chunk) {
    var majorTerrains = '';
    var m = 0;
    for(var key in chunk.terrains){
        if (!chunk.terrains.hasOwnProperty(key)) continue;
        if(key != 'urban' && key != 'roads' && key != 'rivers'){
            if (m < +chunk.terrains[key].percentArea) {
                m = +chunk.terrains[key].percentArea;
                majorTerrains = key;
            }
        }
    }
    return majorTerrains;
}

Template.chunk.helpers({
    data: function() {
        let x = map[this.row][this.col].x;
        let y = map[this.row][this.col].y;
        return 'x = ' + x + ', y = ' + y;
    }
});

Template.worldMap.onRendered(()=>{
    this.$('.chunk').each((i, item) => {
        let curChunk = map[+$(item).attr("data-row")][+$(item).attr("data-col")];
        let mainTerrain = getMainTerrain(curChunk);
        let mainTerrainFile = 'url(resources/' + mainTerrain + '.png)';

        

        $(item).css({
            'backgroundImage': mainTerrainFile
        });
    });
});

Template.worldMap.onCreated(() => {
    for(let i = 0; i < map.length; i++){
        map[i] = new Array(mapColumns);
        for(let j = 0; j < map[0].length; j++) {
            map[i][j] = getChunk(startX + j, startY - i);
        }
    }
});