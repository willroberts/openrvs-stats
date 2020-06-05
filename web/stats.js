// stats.js
const xhr = new XMLHttpRequest();
var stats;
var statsLoaded = false;
var windowLoaded = false;

// Get stats JSON.
xhr.open("GET", "http://127.0.0.1:8081/servers");
xhr.onreadystatechange = function () {
	if (xhr.readyState !== XMLHttpRequest.DONE) {
		return;
	}
	if (xhr.status !== 200) {
		console.log("error: bad response code:", xhr.status);
		return;
	}
	stats = JSON.parse(xhr.responseText);

	statsLoaded = true;
	if (windowLoaded) {
		renderPage();
	}
};
xhr.send();

// Make sure DOM has loaded.
window.onload = function () {
	windowLoaded = true;
	if (statsLoaded) {
		renderPage();
	}
};

// Render the page.
function renderPage() {
	container = document.getElementById("jscontainer");
	h1 = document.createElement("h1");
	h1.appendChild(document.createTextNode("Raven Shield: Who's online now?"));
	container.appendChild(h1);

	var i;
	for (i = 0; i < stats.length; i++) {
		var data = stats[i];

		// Server Name (0/8) [Copy IP:Port]
		p = document.createElement("p");
		p.appendChild(document.createTextNode(data["server_name"] + ' (' + data["current_players"] + '/' + data["max_players"] + ')'));

		b = document.createElement("button");
		b.id = data["ip_address"] + ':' + data["port"];
		b.onclick = function(mouseEvent) {
			// https://stackoverflow.com/questions/400212/how-do-i-copy-to-the-clipboard-in-javascript/6055620#6055620
			window.prompt("Press Ctrl-C and Enter to copy the IP and port.", mouseEvent["target"]["id"]);
		};
		b.appendChild(document.createTextNode('Click to copy IP:Port'));

		p.appendChild(b);
		container.appendChild(p);

		// Mode: X | Map: Y
		p = document.createElement("p");
		p.appendChild(document.createTextNode('Mode: ' + data["game_mode"] + ' | Map: ' + data["current_map"]));
		container.appendChild(p);

		// MOTD: Z
		if (data["motd"] !== "") {
			p = document.createElement("p");
			p.appendChild(document.createTextNode('MOTD: ' + data["motd"]));
			container.appendChild(p);
		}

		// Line before next server
		hr = document.createElement("hr");
		container.appendChild(hr);
	}
};