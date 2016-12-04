const PATH = '/client/static/img/resources/';
const ROAD_Z_INDEX = 10;
const RIVER_Z_INDEX = 20;
var states, mapInfo;
var mapArray = [];
var persons = [];

var setDisplayTime = function(value){
    if(+value < 10){
        value = '0' + value;
    }
    return value;
};

var rand = function (min, max) {
    return Math.floor(min + Math.random()*(max +1 - min));
};

function setTime(time) {
	var textTime = 'Сейчас '+ time.day + ' день ' + time.ten_day +
        ' декады ' + time.month + ' месяца ' + time.year +
        ' года, ' + setDisplayTime(+time.hour) + ':' +
        setDisplayTime(+time.minute);
	$('.worldDate').text(textTime);
}

function getPersonState(state) {
	//console.log(state);
	if (!states) {
		console.log("States are undefined");
		return false;
	}
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
			class: "person-chunk",
			text: "x: " + person.Chunk.X + ", y: " + person.Chunk.Y
		}))
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

function drawPersonOnMap(person) {
	var personChunk = $('[data-row=' + (mapInfo.startMapY - person.Chunk.Y) + 
	'][data-col=' + (person.Chunk.X - mapInfo.startMapX) + ']');
	personChunk.append($('<div>', {
		class: "person-icon",
		id: person.PersonId,
		css: {
			"top": rand(8, 120) + 'px',
			"left": rand(8, 120) + 'px',
		},
		attr: {
			"data-title": person.Name + ', id: ' + person.PersonId
		}
	}));
}

function drawPersons(persons) {
	persons.forEach(function(person, i){
		drawPersonOnMap(person);
		$('.persons').append(createPerson(person));
	});
}

function updatePersonChr(persons) {
	persons.forEach(function(person, i){
		$('div[data-person-id='+person.PersonId+']').find('.state')
		.text("Состояние: " + getPersonState(person.PersonChr.State));
		$('div[data-person-id='+person.PersonId+']').find('.hunger')
		.text("Голод: " + person.PersonChr.Hunger);
		$('div[data-person-id='+person.PersonId+']').find('.thirst')
		.text("Жажда: " + person.PersonChr.Thirst);
		$('div[data-person-id='+person.PersonId+']').find('.fatigue')
		.text("Усталость: " + person.PersonChr.Fatigue);
		$('div[data-person-id='+person.PersonId+']').find('.somnolency')
		.text("Сонливость: " + person.PersonChr.Somnolency);
	});
}

function updatePersonChunk(persons) {
	persons.forEach(function(person, i) {
		$('div[data-person-id='+person.PersonId+']').find('.person-chunk')
		.text("x: " + person.Chunk.X + ", y: " + person.Chunk.Y);
	});
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

function createMapArray(worldMap) {
	//console.log(mapInfo);
	var mapArray = matrixArray(mapInfo.mapWidth, mapInfo.mapHeight);
	//console.log(mapArray);
	worldMap.forEach((chunk, i)=>{
		//console.log(mapInfo.startMapY - chunk.y);
		//console.log(chunk.x - mapInfo.startMapX);
		mapArray[mapInfo.startMapY - chunk.y][chunk.x - mapInfo.startMapX] = {
			"x": chunk.x,
			"y": chunk.y,
			"isExplored": chunk.isExplored,
			"terrains": chunk.terrains
		};
	});
	return mapArray;
}

function chunkHover() {
	$(this).css({'border': "2px solid #770000"});
}

function chunkUnhover() {
	$(this).css({'border': "none"});
}

function chunkClick(event) {
	$('.chunk').css({'border': "none"});
	$(this).css({'border': "2px solid red"});
	var chunk = mapArray[$(this).attr("data-row")][$(this).attr("data-col")];
	$('.map-info').html(createFormattedChunkInfo(chunk));	
}

function createHTMLMap(worldMap) {
	var mapArray = createMapArray(worldMap);
	for(i = 0; i < mapInfo.mapWidth; i++) {
		var mapRow = $('<div>', {class: "map-row"});
		for(j = 0; j < mapInfo.mapHeight; j++) { 
			mapRow.append($('<div>', {
				class: "chunk",
				attr: {
					"data-row": i,
					"data-col": j
				},
				on: {
					/*
					mouseover: chunkHover,
					mouseleave: chunkUnhover,
					*/
					click: chunkClick
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

function createFormattedChunkInfo(chunk){
    var info = 'Координаты чанка<br>x: ' + chunk.x + ', y: ' + chunk.y + '<hr>';
    var towns = '';
    var rivers = '';
    var roads = '';
    var others = '';
    for(var key in chunk.terrains){
        if (!chunk.terrains.hasOwnProperty(key)) continue;
        if(key == 'urban'){
            towns = 'Населенный пункт:<br>' +
            'Занимаемая площадь: ' + chunk.terrains.urban.percentArea + '%<br>' +
            'Тип: ' + chunk.terrains.urban.type + '<br>' +
            'Название: ' + chunk.terrains.urban.townId + '<hr>';
        } else if(key == 'roads'){
            roads = 'Дороги:<br>';
            for(i = 0; i < chunk.terrains.roads.length; i++){
                roads += (i + 1) + ': Размер: ' + chunk.terrains.roads[i].size +
                ', направление: ' + chunk.terrains.roads[i].direction + '<br>';
            }
            roads += '<hr>';
        } else if(key == 'rivers'){
            rivers = 'Реки:<br>';
            for(var i = 0; i < chunk.terrains.rivers.length; i++){
                rivers += (i + 1) + ': Размер: ' + chunk.terrains.rivers[i].size +
                ', направление: ' + chunk.terrains.rivers[i].direction +
                ', качество: ' + chunk.terrains.rivers[i].quality +
                ', мост: ' + chunk.terrains.rivers[i].bridge + '<br>';
            }
        } else {
            others += key + ':<br>' +
            'Занимаемая площадь: ' + chunk.terrains[key].percentArea + '%<br>' +
            'Проходимость: ' + chunk.terrains[key].passability + '<br>' +
            'Качество: ' + chunk.terrains[key].quality + '<hr>';
        }
    }

    info += towns + roads + others + rivers;
    return info;
}

function drawMap() {
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