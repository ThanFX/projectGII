var setDisplayTime = function(value){
    if(+value < 10){
        value = '0' + value;
    }
    return value;
};

var textTime = 'Сейчас '+ time.day + ' день ' + time.ten_day +
        ' декады ' + time.month + ' месяца ' + time.year +
        ' года, ' + setDisplayTime(+time.hour) + ':' +
        setDisplayTime(+time.minute);
$('.worldDate').text(textTime);

function getPersonState(state) {
	var curState = "", curStateValue = "";
	for(var key in states) {
		if(state == states[key]) {
			curState = key;
		}
	}
	switch(curState) {
		case "sleep":
			curStateValue = "Спит";
			break;
		case "chores":
			curStateValue = "Занимается домашними делами";
			break;
		default:
			curStateValue = "Нужно добавить описание!";
	}
	return curStateValue;
}

function createPerson(person) {
	return $('<div>', {
		class: "person",
		attr: {
			"data-person-id": person.PersonId
		},
		on: {
			click: function(event) {
				if($(this).find('.person-charasteristics').css('display') == 'none') {
	        		$(this).find('.person-charasteristics').animate({height: 'show'}, 200);
	    		} else {
	        		$(this).find('.person-charasteristics').animate({height: 'hide'}, 100);
	    		}
			}
		},
		append: $('<span>', {
			class: "name",
			text: person.Name
		})
		.add($('<div>', {
			class: "person-charasteristics",
			append: $('<div>', {
				class: "state",
				text: "Состояние: " + getPersonState(person.PersonChr.State)
			})
			.add($('<div>', {
				class: "hunger",
				text: "Голод: " + person.PersonChr.Hunger.toFixed(2)
			}))
			.add($('<div>', {
				class: "thirst",
				text: "Жажда: " + person.PersonChr.Thirst.toFixed(2)
			}))
			.add($('<div>', {
				class: "fatigue",
				text: "Усталость: " + person.PersonChr.Fatigue.toFixed(2)
			}))
			.add($('<div>', {
				class: "somnolency",
				text: "Сонливость: " + person.PersonChr.Somnolency.toFixed(2)
			}))
		}))
	});
}

function drawPersons(persons) {
	persons.forEach(function(person, i){
		$('.persons').append(createPerson(person));
	})
}

function updatePersonChr(persons) {
	persons.forEach(function(person, i){
		$('div[data-person-id='+person.PersonId+']').find('.state').text("Состояние: " + getPersonState(person.PersonChr.State));
		$('div[data-person-id='+person.PersonId+']').find('.hunger').text("Голод: " + person.PersonChr.Hunger);
		$('div[data-person-id='+person.PersonId+']').find('.thirst').text("Жажда: " + person.PersonChr.Thirst);
		$('div[data-person-id='+person.PersonId+']').find('.fatigue').text("Усталость: " + person.PersonChr.Fatigue);
		$('div[data-person-id='+person.PersonId+']').find('.somnolency').text("Сонливость: " + person.PersonChr.Somnolency);
	})
}

function matrixArray(rows, columns) {
  	var arr = new Array();
  	for(var i = 0; i < columns; i++){
    	arr[i] = new Array();
    	for(var j = 0; j < rows; j++){
      		arr[i][j] = null;
    	}
  	}
  	return arr;
}

function createMapArray() {
	var mapArray = matrixArray(mapInfo.mapWidth, mapInfo.mapHeight);
	worldMap.forEach((chunk, i)=>{
		mapArray[mapInfo.startMapY - chunk.y][chunk.x - mapInfo.startMapX] = {
			"x": chunk.x,
			"y": chunk.y,
			"isExplored": chunk.isExplored,
			"terrains": chunk.terrains
		};
	})
	return mapArray;
}

function createHTMLMap() {
	var mapArray = createMapArray();
	for(i = 0; i < mapInfo.mapWidth; i++) {
		var mapRow = $('<div>', {class: "map-row"});
		for(j = 0; j < mapInfo.mapHeight; j++) { 
			mapRow.append($('<div>', {
				class: "chunk",
				attr: {
					"data-row": i,
					"data-col": j
				}
			}));
		}
		$('.map').append(mapRow);
	}
	return mapArray;
}

function getMainTerrain(chunk) {
    var majorTerrains = '';
    var m = 0;
    for(let key in chunk.terrains){
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

function getMainTerrain(chunk) {
    var majorTerrains = '';
    var m = 0;
    for(let key in chunk.terrains){
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

function getChunkRivers(chunk) {
    let isRiver = !!('rivers' in chunk.terrains);
    if(!isRiver) {
        return false;
    }
    let rivers = [];

    for(let i = 0; i < chunk.terrains.rivers.length; i++){
        let river = {};
        river.size = chunk.terrains.rivers[i].size;
        river.direction = chunk.terrains.rivers[i].direction;
        river.bridge = chunk.terrains.rivers[i].bridge;
        rivers.push(river);
    }
    return rivers;
}

function getChunkRoads(chunk) {
    let isRoad = !!('roads' in chunk.terrains);
    if(!isRoad) {
        return false;
    }
    let roads = [];

    for(let i = 0; i < chunk.terrains.roads.length; i++){
        let road = {};
        road.size = chunk.terrains.roads[i].size;
        road.direction = chunk.terrains.roads[i].direction;
        roads.push(road);
    }
    return roads;
}

function drawMap() {
	var mapArray = createHTMLMap();
	$('.chunk').each((i, item) => {
        let curChunk = mapArray[+$(item).attr("data-row")][+$(item).attr("data-col")];
        let mainTerrain = getMainTerrain(curChunk);
        let mainTerrainFile = 'url(' + PATH + mainTerrain + '.png)';
        $(item).css({
            'backgroundImage': mainTerrainFile
        });

        var rivers = getChunkRivers(curChunk);
        var roads = getChunkRoads(curChunk);
        if(roads) {
           for(i = 0; i < roads.length; i++) {
                let roadFileName = PATH + 'road_' + roads[i].size +
                       '_' + roads[i].direction + '.png';
                let curImg = $("<img>");
                curImg.prop("src", roadFileName);
                curImg.css({'z-index': ROAD_Z_INDEX + roads[i].size});
				
				/* Грязный хак только под текущую захардкоженную карту
				Потом нужно будет переделывать под нормальное определение
				местоположения относительно рек */
				if(roads[i].direction == 'S-N') {
                	curImg.css({'left': "60px"});
                }
				if(roads[i].direction == 'W-E' || roads[i].direction == 'W-C') {
                	curImg.css({'top': "90px"});
                }

                $(item).append(curImg);
            }
        }

        if(rivers) {
           for(i = 0; i < rivers.length; i++) {
               let riverFileName = PATH + 'river_' + rivers[i].size +
                       '_' + rivers[i].direction + '.png';
               let curImg = $("<img>");
               curImg.prop("src", riverFileName);
               curImg.css({'z-index': RIVER_Z_INDEX + rivers[i].size});
               $(item).append(curImg);
            }
        }
    });
}

drawPersons(persons);
drawMap();