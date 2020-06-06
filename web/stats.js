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
	var container = document.getElementById("jscontainer");
	var i;

	var p = document.createElement("p");
	p.appendChild(document.createTextNode('Number of active servers: ' + stats.length));
	container.appendChild(p);
	container.appendChild(document.createElement("hr"));

	for (i = 0; i < stats.length; i++) {
		var data = stats[i];

		// Server Name (0/8)
		var h4 = document.createElement("h4");
		h4.appendChild(document.createTextNode(data["server_name"] + ' (' + data["current_players"] + '/' + data["max_players"] + ') '));

		// [Click to copy IP:Port]
		var b = document.createElement("button");
		b.id = data["ip_address"] + ':' + data["port"];
		b.onclick = function (mouseEvent) {
			// https://stackoverflow.com/questions/400212/how-do-i-copy-to-the-clipboard-in-javascript/6055620#6055620
			window.prompt("Press Ctrl-C and Enter to copy the IP and port.", mouseEvent["target"]["id"]);
		};
		b.appendChild(document.createTextNode('Click to copy IP:Port'));

		h4.appendChild(b);
		container.appendChild(h4);

		// Mode: X | Map: Y
		var p = document.createElement("p");
		p.appendChild(document.createTextNode('Mode: ' + data["game_mode"] + ' | Map: ' + data["current_map"]));
		container.appendChild(p);

		// MOTD: Z
		if (data["motd"] !== "") {
			p = document.createElement("p");
			p.appendChild(document.createTextNode('MOTD: ' + data["motd"]));
			container.appendChild(p);
		}

		// 1x3 table to contain server data
		container.appendChild(createDataTables(data));

		container.appendChild(document.createElement("hr"));
	}// end for loop
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

function createDataTables(data) { // data is json server info
	var table = document.createElement("table");
	table.className = "outertable";

	var tr = document.createElement("tr");

	var td = document.createElement("td");
	td.appendChild(createButton('Active Players'));
	td.appendChild(generatePlayersTable(data));
	tr.appendChild(td);

	td = document.createElement("td");
	td.appendChild(createButton('Map Rotation'));
	td.appendChild(generateMapsTable(data));
	tr.appendChild(td);

	td = document.createElement("td");
	td.appendChild(createButton('Settings'));
	td.appendChild(generateSettingsTable(data));
	tr.appendChild(td);

	table.appendChild(tr);
	return table;
}

function createButton(label) {
	b = document.createElement("button");
	b.appendChild(document.createTextNode(label));
	b.addEventListener("click", collapseButtonHandler);
	return b;
}

function generatePlayersTable(data) { // data is json server info
	var table = document.createElement("table");
	table.className = "innertable";
	table.appendChild(addTableHeaderRow(['Name', 'Kills']));
	var j;
	for (j = 0; j < data["players"].length; j++) {
		var p = data["players"][j];
		table.appendChild(addTableRow([p["name"], p["kills"]]));
	}
	return table;
}

function generateMapsTable(data) { // data is json server info
	var table = document.createElement("table");
	table.className = "innertable";
	table.appendChild(addTableHeaderRow(['Map', 'Game Mode']));
	for (j = 0; j < data["maps"].length; j++) {
		var m = data["maps"][j];
		table.appendChild(addTableRow([m["name"], m["mode"]]));
	}
	return table;
}

function generateSettingsTable(data) { // data is json server info
	var table = document.createElement("table");
	table.className = "innertable";
	table.appendChild(addTableHeaderRow(['Setting', 'Value']));
	if (coopModes.indexOf(data["game_mode"]) === -1) { // Game is adversarial mode.
		var s = data["pvp_settings"];
		table.appendChild(addTableRow(['Auto Team Balance', boolToOnOff(s["auto_team_balance"])]));
		if (data["game_mode"] === "Bomb") {
			table.appendChild(addTableRow(['Bomb Timer', s["bomb_timer"]]));
		}
		table.appendChild(addTableRow(['Friendly Fire', boolToOnOff(s["friendly_fire"])]));
		table.appendChild(addTableRow(['Rounds Per Match', s["rounds_per_match"]]));
		table.appendChild(addTableRow(['Time Per Round', s["time_per_round"]]));
		table.appendChild(addTableRow(['Time_Between_Rounds', s["time_between_rounds"]]));
	} else {
		var s = data["coop_settings"];
		table.appendChild(addTableRow(['AI Backup', boolToOnOff(s["ai_backup"])]));
		table.appendChild(addTableRow(['Friendly Fire', boolToOnOff(s["friendly_fire"])]));
		table.appendChild(addTableRow(['Terrorist Count', s["terrorist_count"]]));
		table.appendChild(addTableRow(['Rotate Map on Success', boolToOnOff(s["rotate_map_on_success"])]));
		table.appendChild(addTableRow(['Rounds Per Match', s["rounds_per_match"]]));
		table.appendChild(addTableRow(['Time Per Round', s["time_per_round"]]));
		table.appendChild(addTableRow(['Time Between Rounds', s["time_between_rounds"]]));
	}
	return table;
}

function addTableHeaderRow(labels) { // labels is []string
	var tr = document.createElement("tr");
	var i;
	for (i = 0; i < labels.length; i++) {
		var th = document.createElement("th");
		th.appendChild(document.createTextNode(labels[i]));
		tr.appendChild(th);
	}
	return tr;
}

function addTableRow(labels) { // labels is []string
	var tr = document.createElement("tr");
	var i;
	for (i = 0; i < labels.length; i++) {
		var td = document.createElement("td");
		td.appendChild(document.createTextNode(labels[i]));
		tr.appendChild(td);
	}
	return tr;
}

function boolToOnOff(b) {
	if (b) {
		return "on";
	} else {
		return "off";
	}
}