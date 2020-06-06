// stats.js
const xhr = new XMLHttpRequest();
var stats;
var statsLoaded = false;
var windowLoaded = false;
var coopModes = ["Hostage Rescue", "Mission", "Terrorist Hunt"];

// Get stats JSON.
xhr.open("GET", "http://64.225.54.237/stats.json");
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
	container.appendChild(document.createElement("hr"));

	var i;
	for (i = 0; i < stats.length; i++) {
		var data = stats[i];

		// Server Name (0/8)
		h3 = document.createElement("h3");
		h3.appendChild(document.createTextNode(data["server_name"] + ' (' + data["current_players"] + '/' + data["max_players"] + ')'));

		// [Click to copy IP:Port]
		b = document.createElement("button");
		b.id = data["ip_address"] + ':' + data["port"];
		b.onclick = function(mouseEvent) {
			// https://stackoverflow.com/questions/400212/how-do-i-copy-to-the-clipboard-in-javascript/6055620#6055620
			window.prompt("Press Ctrl-C and Enter to copy the IP and port.", mouseEvent["target"]["id"]);
		};
		b.appendChild(document.createTextNode('Click to copy IP:Port'));

		h3.appendChild(b);
		container.appendChild(h3);

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

		// Players (collapsible div)
		b = document.createElement("button");
		b.className = "collapsible";
		b.appendChild(document.createTextNode('Players'));
		b.addEventListener("click", collapseButtonHandler);
		container.appendChild(b);
		container.appendChild(document.createTextNode(' (click to expand)'));
		// Hidden div contains list of players
		div = document.createElement("div");
		div.className = "playerlist";
		table = document.createElement("table");
		tr = document.createElement("tr");
		th = document.createElement("th");
		th.appendChild(document.createTextNode('Name'))
		tr.appendChild(th);
		th = document.createElement("th");
		th.appendChild(document.createTextNode('Kills'))
		tr.appendChild(th);
		table.appendChild(tr);
		var j;
		for (j=0; j<data["players"].length; j++) {
			var p = data["players"][j];
			tr = document.createElement("tr");
			td = document.createElement("td");
			td.appendChild(document.createTextNode(p["name"]))
			tr.appendChild(td);
			td = document.createElement("td");
			td.appendChild(document.createTextNode(p["kills"]))
			tr.appendChild(td);
			table.appendChild(tr);
		}
		div.appendChild(table);
		container.appendChild(div);
		container.appendChild(document.createElement("p"));

		// Map Rotation
		b = document.createElement("button");
		b.className = "collapsible";
		b.appendChild(document.createTextNode('Map Rotation'));
		b.addEventListener("click", collapseButtonHandler);
		container.appendChild(b);
		container.appendChild(document.createTextNode(' (click to expand)'));
		// Hidden div contains list of maps and modes
		div = document.createElement("div");
		div.className = "maplist";
		table = document.createElement("table");
		tr = document.createElement("tr");
		th = document.createElement("th");
		th.appendChild(document.createTextNode('Map'))
		tr.appendChild(th);
		th = document.createElement("th");
		th.appendChild(document.createTextNode('Game Mode'))
		tr.appendChild(th);
		table.appendChild(tr);
		for (j=0; j<data["maps"].length; j++) {
			var m = data["maps"][j];
			tr = document.createElement("tr");
			td = document.createElement("td");
			td.appendChild(document.createTextNode(m["name"]))
			tr.appendChild(td);
			td = document.createElement("td");
			td.appendChild(document.createTextNode(m["mode"]))
			tr.appendChild(td);
			table.appendChild(tr);
		}
		div.appendChild(table);
		container.appendChild(div);
		container.appendChild(document.createElement("p"));

		// Settings
		b = document.createElement("button");
		b.className = "collapsible";
		b.appendChild(document.createTextNode('Settings'));
		b.addEventListener("click", collapseButtonHandler);
		container.appendChild(b);
		container.appendChild(document.createTextNode(' (click to expand)'));
		// Hidden div contains list of server settings
		div = document.createElement("div");
		div.className = "settingslist";
		table = document.createElement("table");
		tr = document.createElement("tr");
		th = document.createElement("th");
		th.appendChild(document.createTextNode('Setting'))
		tr.appendChild(th);
		th = document.createElement("th");
		th.appendChild(document.createTextNode('Value'))
		tr.appendChild(th);
		table.appendChild(tr);
		if (coopModes.indexOf(data["game_mode"]) === -1) { // Game is adversarial mode.
			var s = data["pvp_settings"];
			table.appendChild(addTableRow('Auto Team Balance', s["auto_team_balance"]));
			if (data["game_mode"] === "Bomb") {
				table.appendChild(addTableRow('Bomb Timer', s["bomb_timer"]));
			}
			table.appendChild(addTableRow('Friendly Fire', s["friendly_fire"]));
			table.appendChild(addTableRow('Rounds Per Match', s["rounds_per_match"]));
			table.appendChild(addTableRow('Time Per Round', s["time_per_round"]));
			table.appendChild(addTableRow('Time_Between_Rounds', s["time_between_rounds"]));
		} else {
			var s = data["coop_settings"];
			table.appendChild(addTableRow('AI Backup', s["ai_backup"]));
			table.appendChild(addTableRow('Friendly Fire', s["friendly_fire"]));
			table.appendChild(addTableRow('Terrorist Count', s["terrorist_count"]));
			table.appendChild(addTableRow('Rotate Map on Success', s["rotate_map_on_success"]));
			table.appendChild(addTableRow('Rounds Per Match', s["rounds_per_match"]));
			table.appendChild(addTableRow('Time Per Round', s["time_per_round"]));
			table.appendChild(addTableRow('Time Between Rounds', s["time_between_rounds"]));
		}
		div.appendChild(table);
		container.appendChild(div);
		container.appendChild(document.createElement("hr"));
	}
};

function collapseButtonHandler() {
	// https://www.w3schools.com/howto/tryit.asp?filename=tryhow_js_collapsible
	this.classList.toggle("active");
	var content = this.nextElementSibling;
	if (content.style.display === "block") {
	  content.style.display = "none";
	} else {
	  content.style.display = "block";
	}
}

function addTableRow(label, value) {
	tr = document.createElement("tr");
	td = document.createElement("td");
	td.appendChild(document.createTextNode(label));
	tr.appendChild(td);
	td = document.createElement("td");
	td.appendChild(document.createTextNode(value));
	tr.appendChild(td);
	return tr;
}