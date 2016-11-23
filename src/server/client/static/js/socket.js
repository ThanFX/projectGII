window.onload = function () {
    var wd = document.querySelector('.worldDate');
	if (window["WebSocket"]) {
        conn = new WebSocket("ws://localhost:8080/ws");
        conn.onclose = function (evt) {
            var item = document.createElement("div");
            item.innerHTML = "<b>Connection closed.</b>";
        };
        conn.onmessage = function (evt) {
			wd.innerText = evt.data;
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
