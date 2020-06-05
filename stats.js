var stats;
var statsLoaded = false;
var windowLoaded = false;

// Get stats JSON.
const xhr = new XMLHttpRequest();
xhr.open("GET", "http://127.0.0.1:8081/servers");
xhr.onreadystatechange = function () {
	if (xhr.readyState !== XMLHttpRequest.DONE) {
		return;
	}
	if (xhr.status !== 200) {
		console.log("error: bad response code:", xhr.status);
		return;
	}
	stats = xhr.responseText;

	console.log("stats loaded");
	statsLoaded = true;
	if (windowLoaded) {
		renderPage();
	}
};
xhr.send();

// Make sure DOM has loaded.
window.onload = function () {
	console.log("window loaded");
	windowLoaded = true;
	if (statsLoaded) {
		renderPage();
	}
};

// Render the page.
function renderPage() {
	console.log("stats:", stats);
	
	container = document.getElementById("jscontainer");
	p = document.createElement("p");
	p.appendChild(document.createTextNode("example text"));
	container.appendChild(p);
};