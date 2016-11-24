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
					wd.innerText = 'Сейчас '+ message.value.day + ' день ' + message.value.ten_day +
        ' декады ' + message.value.month + ' месяца ' + message.value.year +
        ' года, ' + setDisplayTime(+message.value.hour) + ':' +
        setDisplayTime(+message.value.minute);
					break;
				case "worldMap":
					console.log(message.value);
					break;
			}

			//console.log(evt.data)
			/*
			var messages = evt.data.split('\n');
            for (var i = 0; i < messages.length; i++) {
                var item = document.createElement("div");
                item.innerText = messages[i];
            }
			*/
        };
    } else {
        var item = document.createElement("div");
        item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
    }
};
