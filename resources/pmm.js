// pmm.js
// Pretty Map Machine

/*
	I will leave this section mostly untouched by now, so it won't break now.	
*/
window.addEventListener("load", function(evt) {

	var output = document.getElementById("output");
	var input = document.getElementById("input");
	var playButton = document.getElementById("play");
	var header = document.getElementById("header");
	var spacer = document.getElementById("spacer");

	var ws;

	var print = function(message) {
		let eventObj = html2event(message);
		// If it's not the actual event, skipping it
		if(eventObj === null) {
			return;
		}
		// xy coordinated on image
		let imgXy = ll2xy(eventObj.lat, eventObj.lon);
		// Foreground and background colors
		let fg = [255, 255, 255, 255];
		let bg = [50, 100, 150, 128];
		let text = eventObj.measurement;
		renderText(text, imgXy[0], imgXy[1], fg, bg );
		// var d = document.createElement("div");
		// d.textContent = message;
		// d.innerHTML = '<p style="line-height:40%">' + message + '</p>';
		// d = d.firstChild;
		// output.appendChild(d);
		// spacer.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'start' })
	};

	var yodelBuffer;
	var resBuffer = new Map();

	var context = new (window.AudioContext || window.webkitAudioContext)();

	function play(audioBuffer) {
		const source = context.createBufferSource();
		source.buffer = audioBuffer;
		source.connect(context.destination);
		source.start();
	}

	(function () {
		'use strict';

		window.fetch("yodel.mp3")
			  .then(response => response.arrayBuffer())
			  .then(arrayBuffer => context.decodeAudioData(arrayBuffer,
					audioBuffer => {
						playButton.disabled = false;
						yodelBuffer = audioBuffer;
					},
					error =>
						console.error(error)
			  ))

		function load(resURL, resID) {
			window.fetch(resURL)
				  .then(response => response.arrayBuffer())
				  .then(arrayBuffer => context.decodeAudioData(arrayBuffer,
					audioBuffer => {
						resBuffer.set(resID, audioBuffer)
					},
					error =>
						console.error(error)
				  ))
		}

		var url_string = window.location.href;
		var url = new URL(url_string);
		var c = url.searchParams.get("collection");
		console.log(c);
		var prefix = ""
		if (c != null) {
			prefix = c.concat("/");
		}

		load(prefix.concat("waiting.mp3"), "waiting");

		// Iterating through all the reading types and all the levels
		// [type][level].mp3
		let mp3Types = ["air", "rad"];
		for(let i = 0; i< mp3Types; i++) {
			let mp3Type = mp3Types[i];
			// Levels - 2 digits (00 - 10)
			for(let j = 0; j <= 10; j++) {
				// Making sure it's two letter number
				let mp3num = j.toString().length === 2 ? "0" + j.toString() : j.toString();
				let mp3file = mp3Type + mp3num + ".mp3";
				load(prefix.concat(mp3file), mp3Type + mp3num);
			}
		}

		playButton.onclick = function(evt) {
			playButton.style.display = "none";
			play(yodelBuffer);
			return false;
		};

	}());

	var enterfn = function(evt) {
		if (ws) {
			return false;
		}

		// Compatibility workaround abomination
		// Getting webpage address
		let httpAddress = window.location.href;
		// Determining if webpage has a secure connection
		let protocol = httpAddress.split(":")[0] === "https" ? "wss" : "ws";
		// Determining if port is specified
		let port = (httpAddress.split(":").length > 2) ? (":" + httpAddress.split(":")[2].split("/")[0]) : "";
		// Determining domain
		// 1. Cutting protocol
		// 2. Cutting query string
		// 3. Cutting triling slash
		// 4. Cutting port
		let domain = httpAddress.split("://")[1].split("?")[0].split("/")[0].split(":")[0];
		// Assembling the address
		let websocketAddress = protocol + "://" + domain + port + "/stream/";


		ws = new WebSocket(websocketAddress);
		ws.onopen = function(evt) {
			header.style.display = "inline";
		}
		ws.onclose = function(evt) {
			print("connection closed");
			ws = null;
		}
		ws.onmessage = function(evt) {
			header.style.display = "none";
			var line = evt.data
			print(line);
			var firstWord = line.substr(0, line.indexOf(" "));
			play(resBuffer.get(firstWord));
			//			document.getElementById(firstWord).click()
		}
		ws.onerror = function(evt) {
			print("error: " + evt.data);
		}
		return false;
	};

	var exitfn = function(evt) {
		if (!ws) {
			return false;
		}
		ws.send(input.value);
		return false;
	};

	var e = document.getElementById("open");
	if (e) { e.onclick = enterfn };
	var e = document.getElementById("send")
	if (e) { e.onclick = exitfn };
	var e = document.getElementById("close")
	if (e) {
		e.onclick = function(evt) {
			if (!ws) {
				return false;
			}
			ws.close();
			return false;
		};
	};

	enterfn();
});

/*
	Here is the start of the map code
*/

// DOM elements
let elms = {
	// map image
	mapImg: document.getElementById("mapImg"),
	// canvas overlay
	mapCvs: document.getElementById("mapCvs"),
	// "terminal" canvas
	txtCvs: document.getElementById("txtCvs")
}

// Event - array
// timestamp with ms — when it should happen
// x - "terminal text" x
// y - "terminal text" y,
// fr, fg, fb, fa - RGBA channels of font symbol
// br, bg, bb, ba - RGBA channels of background
// [timestamp with ms, x, y, fr, fg, fb, fa, br, bg, bb, ba]
let terminalDrawingEvents = [

];

// Map canvas context
let mapCvsCtx = elms.mapCvs.getContext("2d");

// mg - map geometry
let mg = {
	// Left/Right edge X
	lx: 66,
	rx: 2017,
	// Top/Bottom edge Y
	ty: 20,
	by: 997
};
// Width/Height
// +1 because it's inclusive range
mg.w = mg.rx - mg.lx + 1;
mg.h = mg.by - mg.ty + 1;
// 0:0 X/Y
// 0 longtitude is dictated by image properties
mg.cx = 980;
mg.cy = mg.ty + (mg.h / 2);
// 1 degree latitude(y)/longtitude(x)
mg.xd = 5.43;
mg.yd = (mg.by - mg.ty) / 180;
// Test output
// console.log(mg);

// latitude/longtitude to x/y
let ll2xy = function(lat, lon) {
	let x = Math.round(mg.cx + (lon * mg.xd));
	let y = Math.round(mg.cy - (lat * mg.yd));
	return [x, y];
}
// Test output
// console.log(ll2xy(34.122, -118.261));

// Terminus Bold 8x16
// TODO: refactor array
let font = [0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,24,24,24,24,24,24,24,0,24,24,0,0,0,0,0,102,102,102,0,0,0,0,0,0,0,0,0,0,0,0,0,0,108,108,108,254,108,108,254,108,108,108,0,0,0,0,0,16,16,124,214,208,208,124,22,22,214,124,16,16,0,0,0,0,102,214,108,12,24,24,48,54,107,102,0,0,0,0,0,0,56,108,108,56,118,220,204,204,220,118,0,0,0,0,0,24,24,24,0,0,0,0,0,0,0,0,0,0,0,0,0,0,12,24,48,48,48,48,48,48,24,12,0,0,0,0,0,0,48,24,12,12,12,12,12,12,24,48,0,0,0,0,0,0,0,0,0,108,56,254,56,108,0,0,0,0,0,0,0,0,0,0,0,24,24,126,24,24,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,24,24,48,0,0,0,0,0,0,0,0,0,0,254,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,24,24,0,0,0,0,0,0,6,6,12,12,24,24,48,48,96,96,0,0,0,0,0,0,124,198,198,206,222,246,230,198,198,124,0,0,0,0,0,0,24,56,120,24,24,24,24,24,24,126,0,0,0,0,0,0,124,198,198,6,12,24,48,96,192,254,0,0,0,0,0,0,124,198,198,6,60,6,6,198,198,124,0,0,0,0,0,0,6,14,30,54,102,198,254,6,6,6,0,0,0,0,0,0,254,192,192,192,252,6,6,6,198,124,0,0,0,0,0,0,60,96,192,192,252,198,198,198,198,124,0,0,0,0,0,0,254,6,6,12,12,24,24,48,48,48,0,0,0,0,0,0,124,198,198,198,124,198,198,198,198,124,0,0,0,0,0,0,124,198,198,198,198,126,6,6,12,120,0,0,0,0,0,0,0,0,0,24,24,0,0,0,24,24,0,0,0,0,0,0,0,0,0,24,24,0,0,0,24,24,48,0,0,0,0,0,0,6,12,24,48,96,48,24,12,6,0,0,0,0,0,0,0,0,0,254,0,0,254,0,0,0,0,0,0,0,0,0,0,96,48,24,12,6,12,24,48,96,0,0,0,0,0,0,124,198,198,198,12,24,24,0,24,24,0,0,0,0,0,0,124,198,206,214,214,214,214,206,192,126,0,0,0,0,0,0,124,198,198,198,198,254,198,198,198,198,0,0,0,0,0,0,252,198,198,198,252,198,198,198,198,252,0,0,0,0,0,0,124,198,198,192,192,192,192,198,198,124,0,0,0,0,0,0,248,204,198,198,198,198,198,198,204,248,0,0,0,0,0,0,254,192,192,192,248,192,192,192,192,254,0,0,0,0,0,0,254,192,192,192,248,192,192,192,192,192,0,0,0,0,0,0,124,198,198,192,192,222,198,198,198,124,0,0,0,0,0,0,198,198,198,198,254,198,198,198,198,198,0,0,0,0,0,0,60,24,24,24,24,24,24,24,24,60,0,0,0,0,0,0,30,12,12,12,12,12,12,204,204,120,0,0,0,0,0,0,198,198,204,216,240,240,216,204,198,198,0,0,0,0,0,0,192,192,192,192,192,192,192,192,192,254,0,0,0,0,0,0,130,198,238,254,214,198,198,198,198,198,0,0,0,0,0,0,198,198,198,230,246,222,206,198,198,198,0,0,0,0,0,0,124,198,198,198,198,198,198,198,198,124,0,0,0,0,0,0,252,198,198,198,198,252,192,192,192,192,0,0,0,0,0,0,124,198,198,198,198,198,198,198,222,124,6,0,0,0,0,0,252,198,198,198,198,252,240,216,204,198,0,0,0,0,0,0,124,198,192,192,124,6,6,198,198,124,0,0,0,0,0,0,255,24,24,24,24,24,24,24,24,24,0,0,0,0,0,0,198,198,198,198,198,198,198,198,198,124,0,0,0,0,0,0,198,198,198,198,198,108,108,108,56,56,0,0,0,0,0,0,198,198,198,198,198,214,254,238,198,130,0,0,0,0,0,0,198,198,108,108,56,56,108,108,198,198,0,0,0,0,0,0,195,195,102,102,60,24,24,24,24,24,0,0,0,0,0,0,254,6,6,12,24,48,96,192,192,254,0,0,0,0,0,0,60,48,48,48,48,48,48,48,48,60,0,0,0,0,0,0,96,96,48,48,24,24,12,12,6,6,0,0,0,0,0,0,60,12,12,12,12,12,12,12,12,60,0,0,0,0,0,24,60,102,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,254,0,0,48,24,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,124,6,126,198,198,198,126,0,0,0,0,0,0,192,192,192,252,198,198,198,198,198,252,0,0,0,0,0,0,0,0,0,124,198,192,192,192,198,124,0,0,0,0,0,0,6,6,6,126,198,198,198,198,198,126,0,0,0,0,0,0,0,0,0,124,198,198,254,192,192,124,0,0,0,0,0,0,30,48,48,252,48,48,48,48,48,48,0,0,0,0,0,0,0,0,0,126,198,198,198,198,198,126,6,6,124,0,0,0,192,192,192,252,198,198,198,198,198,198,0,0,0,0,0,0,24,24,0,56,24,24,24,24,24,60,0,0,0,0,0,0,6,6,0,14,6,6,6,6,6,6,102,102,60,0,0,0,192,192,192,198,204,216,240,216,204,198,0,0,0,0,0,0,56,24,24,24,24,24,24,24,24,60,0,0,0,0,0,0,0,0,0,252,214,214,214,214,214,214,0,0,0,0,0,0,0,0,0,252,198,198,198,198,198,198,0,0,0,0,0,0,0,0,0,124,198,198,198,198,198,124,0,0,0,0,0,0,0,0,0,252,198,198,198,198,198,252,192,192,192,0,0,0,0,0,0,126,198,198,198,198,198,126,6,6,6,0,0,0,0,0,0,222,240,224,192,192,192,192,0,0,0,0,0,0,0,0,0,126,192,192,124,6,6,252,0,0,0,0,0,0,48,48,48,252,48,48,48,48,48,30,0,0,0,0,0,0,0,0,0,198,198,198,198,198,198,126,0,0,0,0,0,0,0,0,0,198,198,198,108,108,56,56,0,0,0,0,0,0,0,0,0,198,198,214,214,214,214,124,0,0,0,0,0,0,0,0,0,198,198,108,56,108,198,198,0,0,0,0,0,0,0,0,0,198,198,198,198,198,198,126,6,6,124,0,0,0,0,0,0,254,12,24,48,96,192,254,0,0,0,0,0,0,28,48,48,48,96,48,48,48,48,28,0,0,0,0,0,0,24,24,24,24,24,24,24,24,24,24,0,0,0,0,0,0,112,24,24,24,12,24,24,24,24,112,0,0,0,0,0,115,219,206,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0];
let fontData = {
	symbolWidth: 8,
	symbolHeight: 16,
	// Lowest codepoint in the array
	// 32 because we are omitting 32 non-printable symbols
	loCode: 32,
	fontLength: font.length / 16
};

// code - character code
// x, y - x, y coordinates
// f, b - foreground and background colors as [r, g, b, a] arrays
let renderSymbol = function(code, x, y, f, b) {
	// Intermediate off-screen canvas
	let iCvs = document.createElement("canvas");
	iCvs.height = fontData.symbolHeight;
	iCvs.width = fontData.symbolWidth;
	// Intermediate canvas context
	iCtx = iCvs.getContext("2d");
	// Intermediate image data
	iID = iCtx.getImageData(0, 0, iCvs.width, iCvs.height);

	// Iterating through pixel rows
	for(let rowIndex = 0; rowIndex < iCvs.height; rowIndex++) {
		// Getting font data for a row
		let rowByte = font[((code - 32) * iCvs.height) + rowIndex];
		// Iterating through pixel columns, 1 bit of font per ixel
		for(let columnIndex = 0; columnIndex < iCvs.width; columnIndex++) {
			// Getting ImageData offset for current pixel
			// Multipling by 4 because of 4 color channels
			let pixelOffset = ((rowIndex * 8) + columnIndex) * 4;
			let andOperand = 0b00000001 << (7 - columnIndex);
			let isFilled = (rowByte & andOperand) >> (7 - columnIndex);
			let color = isFilled ? f : b;

			for(let colorOffset = 0; colorOffset < 4; colorOffset++) {
				iID.data[pixelOffset + colorOffset] = color[colorOffset];
			}
		}
	}
	// Rendering result
	mapCvsCtx.putImageData(iID, x, y);
}
// Test output
for(let i = 0; i < 128; i++) {
	renderSymbol(i, (256 + (8 * (i - 32))), 256, [255, 255, 255, 255], [0, 0, 0, 0, 128]);
}

let renderText = function(text, x, y, f, b) {
	for(let i = 0; i < text.length; i++) {
		let code = text.charCodeAt(i);
		renderSymbol(code, x + (i * fontData.symbolWidth), y, f, b);
	}
}
// Test output
// renderText("It is wednesday my dudes!", 256, 384, [64, 64, 255, 255], [255, 255, 255, 128]);

let clearSymbol = function(x, y) {
	renderSymbol(32, x, y, [0, 0, 0, 0], [0, 0, 0, 0]);
}

let delay = function(x) {
	for(let i = 0; i < x; i++) {
		// Make something great
	}
}

// Conversion of websocket events to objects
let html2event = function(crypticSigils) {
	// Websocket message
	// let crypticSigils = 'rad06 60.00% <a href="https://tt.safecast.org/id/pointcast:10036" target="_blank">POINTCAST #10036</a> EC7128 14.0cpm <a href="https://www.google.com/maps/search/?api=1&query=37.494708,139.926040" target="_blank">7768km</a>';
	crypticSigils = crypticSigils.trim();

	// If there's no event returning nothing
	if(crypticSigils === "waiting for events") {
		return null;
	}

	// Off-screen DOM "rendering"
	let manifestation = document.createElement("div");
	// Rendering off-screen
	manifestation.innerHTML = crypticSigils;
	// Link to device
	deviceLink = manifestation.children[0];
	// Link to google maps
	gmapsLink = manifestation.children[1];

	// Full text
	let sigilText = manifestation.textContent;
	let splitSigil = sigilText.split(" ");

	// Signal type
	let signalType = splitSigil[0];
	// Signal percentage
	// "96.22%" >>> 96.22
	let signalPercentage = parseFloat(splitSigil[1].substring(0, (splitSigil[1].length - 1) ) );

	// Device
	let deviceName = deviceLink.textContent;
	let deviceUrl = deviceLink.href;

	// Sensor type
	let sensorType = splitSigil[splitSigil.length - 3];

	// Location
	let gmapsUrl = gmapsLink.href;
	let splitGmapsUrl = gmapsLink.href.split("=");
	splitGmapsUrl = splitGmapsUrl[splitGmapsUrl.length - 1].split(",");
	let lat = parseFloat(splitGmapsUrl[splitGmapsUrl.length - 2]);
	let lon = parseFloat(splitGmapsUrl[splitGmapsUrl.length - 1]);

	// Measurement
	let measurement = splitSigil[splitSigil.length - 2];

	// Result
	let result = {
		type: signalType,
		percentage: signalPercentage,
		name: deviceName,
		url: deviceUrl,
		sensor: sensorType,
		lat: lat,
		lon: lon,
		gmapsUrl: gmapsUrl,
		measurement: measurement
	};

	return result;
}

// We are having 2048x1024 canvas with 8x16 symbols
// pixel xy to "terminal" text xy
let px2tt = function(px, py) {
	let tx = Math.floor(px / 8);
	let ty = Math.floor(px / 16);

	return [tx, ty];
}
// "terminal" text xy to pixel xy
let tt2px = function(tx, ty) {
	let px = tx * 8;
	let py = ty * 8;

	return [px, py];
}

let terminalString = function(textToOutput, x, y, fg, bg) {
	
}
