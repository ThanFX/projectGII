var setDisplayTime = function(value){
    if(+value < 10){
        value = '0' + value;
    }
    return value;
};

window.onload = function () {
    var wd = document.querySelector('.worldDate');
	if (window["WebSocket"]) {
        conn = new WebSocket("ws://localhost:8080/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
        };
        conn.onmessage = function (evt) {
			//console.log(evt.data);
			message = JSON.parse(evt.data);
			//console.log(message);

			switch (message.key) {
				case "time":
					setTime(message.value);
					break;
				case "worldMap":
					mapArray = createHTMLMap(message.value);
                    drawMap();
					break;
				case "persons":
					if(persons.length == 0) {
                        persons = message.value;
                        drawPersons(persons);    
                    } else {
                        persons = message.value;
                        updatePersonChr(persons);
                    }
					break;
                /*    
				case "states":
					states = message.value;
					break;
                */
                case "mapInfo":
                    mapInfo = message.value;
                    break;
			}
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    }
};
