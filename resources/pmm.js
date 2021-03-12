// pmm.js
// Pretty Map Machine

let getRandomInt = function(max) {
  return Math.floor(Math.random() * Math.floor(max));
}


/*
	I will leave this section mostly untouched by now, so it won't break now.	
*/
var yodelBuffer;
var resBuffer = new Map();
var context = new (window.AudioContext || window.webkitAudioContext)();
// context.suspend();
let currentRenderingY = 0;

var print = function(message) {
	let eventObj = html2event(message);
	// If it's not the actual event, skipping it
	if(eventObj === null) {
		terminalString("no signal", getRandomInt(elms.showCvs.width) - 80, currentRenderingY, [255,255,255,255], [0,0,255,255]);
		return;
	}
	// xy coordinated on image
	let imgXy = ll2xy(eventObj.lat, eventObj.lon);
	// Foreground and background colors
	let level = parseInt(eventObj.type.substring(3));
	let fg = [255, 255, 255, 255];
	let bg = [5 + (level * 25), 255 - (level * 25), 32, 200];
	let text = eventObj.measurement;

	let patchX = imgXy[0] - (elms.showCvs.width / 2);
	let xOffset = 0;
	if(patchX < 44) {
		xOffset = 44 - patchX;
		patchX = 44;
	}
	let patchW = elms.showCvs.width;
	let patchH = 96 + getRandomInt(32);
	if(patchH % 2 === 1) {
		patchH += 1;
	}
	let patchY = imgXy[1] - patchH / 2;
	let mapPatchImageData = mapCvsCtx.getImageData(patchX, patchY, patchW, patchH);

	let rShift = getRandomInt(100);
	let gShift = getRandomInt(100);
	let bShift = getRandomInt(100);

	let invertedDice = getRandomInt(100);

	let glitchOffset = getRandomInt(patchW * 0.75 * 4);
	let glitchRShift = getRandomInt(50);
	let glitchGShift = getRandomInt(50);
	let glitchBShift = getRandomInt(50);

	for(let i = 0; i < mapPatchImageData.data.length; i+=4) {
		mapPatchImageData.data[i + 0] += rShift;
		mapPatchImageData.data[i + 1] += gShift;
		mapPatchImageData.data[i + 2] += bShift;
		if(mapPatchImageData.data[i + 3] === 0) {
			mapPatchImageData.data[i + 3] = 255;
			mapPatchImageData.data[i + 0] = mapPatchImageData.data[i - glitchOffset] - glitchRShift;
			mapPatchImageData.data[i + 1] = mapPatchImageData.data[i - glitchOffset] - glitchGShift;
			mapPatchImageData.data[i + 2] = mapPatchImageData.data[i - glitchOffset] - glitchBShift;
		}

		if(invertedDice >= 50) {
			mapPatchImageData.data[i + 0] = 255 - mapPatchImageData.data[i + 0];
			mapPatchImageData.data[i + 1] = 255 - mapPatchImageData.data[i + 1];
			mapPatchImageData.data[i + 2] = 255 - mapPatchImageData.data[i + 2];
		}
		
	}
	showCvsCtx.putImageData(mapPatchImageData, 0, currentRenderingY);

	sigX = (patchW / 2) - xOffset;
	sigY = (patchH / 2) + currentRenderingY;

	currentRenderingY += patchH;
	if(currentRenderingY > vh) {
		currentRenderingY = 0;
	}

	
	let measurementTxt = "MEASUREMENT " + measurementNo++;
	terminalString(measurementTxt, sigX - (measurementTxt.length * 8) - 8, sigY, [255,255,255,255], [0,0,0,0,255]);
	let typeTxt = "TYPE: " + eventObj.type;
	terminalString(typeTxt, sigX - (measurementTxt.length * 8) - 8, sigY + 16, [255,255,255,255], [0,0,0,0,255]);

	let latLonTxt = "LAT:" + eventObj.lat + " LON:" + eventObj.lon;
	terminalString(latLonTxt, sigX - (measurementTxt.length * 8) - 16, sigY - 16, [255,255,255,255], [0,0,255,255]);
	
	terminalString(text, sigX, sigY, fg, bg);
	terminalString(eventObj.sensor, sigX, sigY + 16, fg, bg);
	
	// var d = document.createElement("div");
	// d.textContent = message;
	// d.innerHTML = '<p style="line-height:40%">' + message + '</p>';
	// d = d.firstChild;
	// output.appendChild(d);
	// spacer.scrollIntoView({ behavior: 'smooth', block: 'nearest', inline: 'start' })
};

window.addEventListener("load", function(evt) {

	var output = document.getElementById("output");
	var input = document.getElementById("input");
	var playButton = document.getElementById("play");
	var header = document.getElementById("header");
	var spacer = document.getElementById("spacer");

	var ws;




	function play(audioBuffer) {
		const source = context.createBufferSource();
		source.buffer = audioBuffer;
		source.connect(context.destination);
		source.start();
	}(function () {
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
		var prefix = ""
		if (c != null) {
			prefix = c.concat("/");
		}

		load(prefix.concat("waiting.mp3"), "waiting");
		load(prefix.concat("waiting.mp3"), "waiting");
		load(prefix.concat("air00.mp3"), "air00");
		load(prefix.concat("air01.mp3"), "air01");
		load(prefix.concat("air02.mp3"), "air02");
		load(prefix.concat("air03.mp3"), "air03");
		load(prefix.concat("air04.mp3"), "air04");
		load(prefix.concat("air05.mp3"), "air05");
		load(prefix.concat("air06.mp3"), "air06");
		load(prefix.concat("air07.mp3"), "air07");
		load(prefix.concat("air08.mp3"), "air08");
		load(prefix.concat("air09.mp3"), "air09");
		load(prefix.concat("air10.mp3"), "air10");
		load(prefix.concat("rad00.mp3"), "rad00");
		load(prefix.concat("rad01.mp3"), "rad01");
		load(prefix.concat("rad02.mp3"), "rad02");
		load(prefix.concat("rad03.mp3"), "rad03");
		load(prefix.concat("rad04.mp3"), "rad04");
		load(prefix.concat("rad05.mp3"), "rad05");
		load(prefix.concat("rad06.mp3"), "rad06");
		load(prefix.concat("rad07.mp3"), "rad07");
		load(prefix.concat("rad08.mp3"), "rad08");
		load(prefix.concat("rad09.mp3"), "rad09");
		load(prefix.concat("rad10.mp3"), "rad10");

		/*
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
		*/

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
let measurementNo = 0;

// DOM elements
let elms = {
	// map image
	mapImg: document.getElementById("mapImg"),
	// canvas overlay
	mapCvs: document.getElementById("mapCvs"),
	// "terminal" canvas
	txtCvs: document.getElementById("txtCvs"),
	// Branding canvas
	brandingCvs: document.getElementById("brandingCvs"),
	// Show canvas
	showCvs: document.getElementById("showCvs"),
	// Sharing canvas
	sharingCvs: document.getElementById("sharingCvs"),
};

// Map canvas context
let mapCvsCtx = elms.mapCvs.getContext("2d");
// Branding canvascontext
let brandingCvsCtx = elms.brandingCvs.getContext("2d");
// SHOW canvas context
let showCvsCtx = elms.showCvs.getContext("2d");
// Terminal canvas context
let txtCvsCtx = elms.txtCvs.getContext("2d");
// Sharing canvas context
let sharingCvsCtx = elms.sharingCvs.getContext("2d");

let mapCached = false;
let cacheMap = function() {
	// Caching the map
	mapCvsCtx.drawImage(elms.mapImg, 0, 0);
	mapCached = true;

	let mapID = mapCvsCtx.getImageData(44, 0, vw - 44, vh);
}

let applyBranding = function() {
	// Caching map in the intermediate canvas to get image data (if not cached previously)
	if(!mapCached) {
		cacheMap();
	}
	// getting branding strip image data
	let brandingStripImageData = mapCvsCtx.getImageData(0, 0, 44, 1024);
	// Getting screen height / branding height porportions to zoom
	brandingCvs.width = 44;
	brandingCvs.height = vh;
	brandingCvsCtx.putImageData(brandingStripImageData, 0, 0);
	placeContainer();
}

let placeContainer = function() {
	let showCvsImageData = showCvsCtx.getImageData(0, 0, elms.showCvs.width, elms.showCvs.height);
	elms.showCvs.width = vw - 44;
	elms.showCvs.height = vh;
	showCvsCtx.putImageData(showCvsImageData, 0, 0);

	showCvsImageData = showCvsCtx.getImageData(0, 0, elms.showCvs.width, elms.showCvs.height);

	for(let i = 0; i < showCvsImageData.data.length; i+=4) {
		if(showCvsImageData.data[i+3] === 0) {
			let greyscale = getRandomInt(32);
			showCvsImageData.data[i+0] = greyscale;
			showCvsImageData.data[i+1] = greyscale;
			showCvsImageData.data[i+2] = greyscale;
			showCvsImageData.data[i+3] = 255;
		}
	}

	showCvsCtx.putImageData(showCvsImageData, 0, 0);

	let txtCvsImageData = txtCvsCtx.getImageData(0, 0, elms.txtCvs.width, elms.txtCvs.height);
	elms.txtCvs.width = vw - 44;
	elms.txtCvs.height = vh;
	txtCvsCtx.putImageData(txtCvsImageData, 0, 0);

}

// Viewport dimensions
let vw = null;
let vh = null;
let portrait = false;
let proportion = null;

let getViewportSize = function() {
	vw = Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0);
	vh = Math.max(document.documentElement.clientHeight || 0, window.innerHeight || 0);
	proportion = vh/1024;
	portrait = vw < vh;
	applyBranding();
}

elms.mapImg.onload = function() {
	cacheMap();
	applyBranding();
}
getViewportSize();
window.onresize = getViewportSize;

// Current frame timestamp (unixtime with milliseconds) for renering purposes
let ts = null;
// Lock to indicate in-progress rendering
let renderingLock = null;
// Terminal Drawing Sequence
// Event - array
// timestamp with ms — when it should happen
// c - symbol code
// x - "terminal text" x
// y - "terminal text" y,
// fr, fg, fb, fa - RGBA channels of font symbol
// br, bg, bb, ba - RGBA channels of background
// [timestamp with ms, c, x, y, fr, fg, fb, fa, br, bg, bb, ba]
let tds = [
	
];
// Fixed amount of items
for(let i = 0; i < 2048; i++) {
	tds.push(null);
}

let renderTerminalFrame = function() {
	if(renderingLock) {
		return false;
	}

	renderingLock = true;

	let dateObj = new Date();
	ts = dateObj.getTime();

	for(let i = 0; i < 2048; i++) {
		// e - current event
		let e = tds[i];

		// If there's no event, skipping
		if(!e) {
			continue;
		}

		// ets - event timestamp
		let ets = e[0];
		// If the time didn't come yet, skipping
		if(ets > ts) {
			continue;
		}

		// c, xy - symbol code and placement on a screen
		let c = e[1];
		let xy = tt2px(e[2], e[3]);
		// Foreground/background colors
		let fg = [ e[4], e[5], e[6], e[7] ];
		let bg = [ e[8], e[9], e[10], e[11] ];

		renderSymbol(c, xy[0], xy[1], fg, bg);
		tds[i] = null;
	}

	renderingLock = false;
}


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
	txtCvsCtx.putImageData(iID, x, y);
}
// Test output
/*
for(let i = 0; i < 128; i++) {
	renderSymbol(i, (256 + (8 * (i - 32))), 256, [255, 255, 255, 255], [0, 0, 0, 0, 128]);
}
*/

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
	let ty = Math.floor(py / 16);

	return [tx, ty];
}
// "terminal" text xy to pixel xy
let tt2px = function(tx, ty) {
	let px = tx * 8;
	let py = ty * 16;

	return [px, py];
}

let framerate = 10;

let terminalString = function(textToOutput, x, y, fg, bg, dontClear = false) {
	let dateObj = new Date();
	let startTs = dateObj.getTime();
	for(let i = 0; i < textToOutput.length; i++) {
		// Event timestamp 
		let ets = startTs + ((1000 / framerate) * i);
		// Symbol code
		let c = textToOutput.charCodeAt(i);
		// "Terminal" xy
		let xy = px2tt(x, y);
		// Horizontal symbol offset
		xy[0] += i;
		// Rendering event array
		let re = [ets, c, xy[0], xy[1], fg[0], fg[1], fg[2], fg[3], bg[0], bg[1], bg[2], bg[3]];
		// delay 5000 ms before cleanup
		let ce = [ets + 5000, 32, xy[0], xy[1], 0, 0, 0, 0, 0, 0, 0, 0];

		// Flags to see if events were placed into Terminal Drawing Sequence
		let rePlaced = false;
		let cePlaced = false;
		// 
		for(let j = 0; j < 2048; j++) {
			let e = tds[j];

			// if there's no event, placing our own event
			if(e === null && (!rePlaced || (!cePlaced && !dontClear))) {
				if(!rePlaced) {
					tds[j] = re;
					rePlaced = true;
				}
				else if(!cePlaced && !dontClear) {
					tds[j] = ce;
					cePlaced = true;
				}
			}
			
			// Nothing else can be done to the empty event, skipping it
			if(e === null) {
				continue;
			}

			// If it's a cleanup event after rendering our symbol before its cleanup, erasing it
			// Coordinates
			if(e[2] === xy[0] && e[3] === xy[1]) {
				// If it's cleanup
				if(e[7] === 0 && e[11] === 0) {
					// If it prevents our symbol ot be shown for 5 seconds
					if(e[0] >= ets && e[0] <= (ets + 5000)) {
						tds[j] = null;
					}
				}
			}
		}
	}
}
setInterval(renderTerminalFrame, 1000 / framerate);

// Test output
// terminalString("hello", 35, 22, [255, 255, 255, 255], [0, 0, 0, 255]);

let explainer = [
	"TRANSMISSION EXPLAINER:",
	"",
	"With inspiration from Listen To Wikipedia, ",
	"Brian Eno’s ambient works, and Nine Inch Nails Ghosts I-IV - ",
	"safecast.live uses the real time data stream coming in from ",
	"Safecast’s global network of environmental sensors to trigger a ",
	"random and ever evolving anti-pattern of audio samples. ",
	"",
	"Using hand held mobile sensors, in the 10 years since the 3/11 ",
	"earthquake and nuclear meltdown at Fukushima Daiichi the non-profit",
	"organization Safecast has built the largest collection of ",
	"background radiation measurements ever assembled and placed them ",
	"entirely into the public domain. Over the last few years a growing ",
	"network of Safecast designed and maintained sensors have begun to ",
	"send in live radiation readings from around the world. Marking the ",
	"10 year anniversary of Safecast, Blues Wireless has designed a ",
	"fleet of air quality sensors that will do the same, sending open ",
	"data into Safecast’s system to be freely published out to the ",
	"world. ",
	"",
	"You are listening to that data stream right now, radiation ",
	"measurements trigger some samples, air quality measurements trigger ",
	"others. As these measurements are an ever changing random stream, ",
	"so will this audio constantly evolve."
];

let credits = [
	"CREDITS:",
	"",
	"Built by:",
	"    Kether Cortex (Design & Development), ",
	"    Ray Ozzie (Development & Concept), ",
	"    Sean Bonner (Audio Editing & Concept),",
	"    Evgeniy Kaptan (Development, Prototyping)",
	"",
	"Audio by: ",
	"    Nine Inch Nails (Ghosts samples, used under CC), ",
	"    Hainbach (Michelsonne toy piano samples, used under CC), ",
	"    Samples From Mars (Polivox synthesizer samples), ",
	"    Sean Bonner (Buchla music easel samples)"
]

let renderExplainer = function() {
	let x = 0;
	let y = 60;
	let fg = [200, 220, 40, 255];
	let bg = [0, 0, 0, 200];

	for(let i = 0; i < explainer.length; i++) {
		terminalString(explainer[i], x, y + (16 * i), fg, bg, true);
	}

	setTimeout(renderCredits, 10000);
}

let renderCredits = function() {
	let x = 768;
	let y = 768;
	let fg = [200, 220, 40, 255];
	let bg = [0, 0, 0, 200];

	for(let i = 0; i < credits.length; i++) {
		terminalString(credits[i], x, y + (16 * i), fg, bg, true);
	}
}

// document.getElementById("explainer_credits").onclick = renderExplainer;

document.getElementById("play_btn").onclick = function() {
	if(context.state === "running") {
		context.suspend();
	}
	else {
		context.resume();
	}
}

document.getElementById("actualMapContainer").onclick = function() {
	let navigationEl = document.getElementById("navigation");
	navigationEl.style.display = (navigationEl.style.display === "block") ? "none" : "block";
}

document.getElementById("about_btn").onclick = function() {
	let about = document.getElementById("about");
	about.style.display = "block";
	about.style.zIndex = 2048;
}

document.getElementById("close_about").onclick = function() {
	let about = document.getElementById("about");
	about.style.display = "none";
	about.style.zIndex = 0;
}

document.getElementById("credits_btn").onclick = function() {
	let credits = document.getElementById("credits");
	credits.style.display = "block";
	credits.style.zIndex = 2048;
}

document.getElementById("close_credits").onclick = function() {
	let credits = document.getElementById("credits");
	credits.style.display = "none";
	credits.style.zIndex = 0;
}


let fileToShare = null;
let codeToShare = null;
let storageUrl = "https://safecast-live-us-west-2-ith8aefe.s3.us-west-2.amazonaws.com/";

let sharing = function() {
	elms.sharingCvs.style.display = "block";
	elms.sharingCvs.style.zIndex = 2048;
	elms.sharingCvs.width = vw;
	elms.sharingCvs.height = vh;
	// ImageDatas for branding, map and letters
	let bID = brandingCvsCtx.getImageData(0, 0, elms.brandingCvs.width, elms.brandingCvs.height);
	sharingCvsCtx.putImageData(bID, 0, 0);
	let sID = showCvsCtx.getImageData(0, 0, elms.showCvs.width, elms.showCvs.height);
	let tID = txtCvsCtx.getImageData(0, 0, elms.txtCvs.width, elms.txtCvs.height);
	
	for(let i = 0; i < sID.data.length; i+=4) {
		if(sID.data[i+3] === 0) {
			let greyscale = getRandomInt(255);
			sID.data[i+0] = greyscale;
			sID.data[i+1] = greyscale;
			sID.data[i+2] = greyscale;
			sID.data[i+3] = 255;
		}

		if(tID.data[i+3] === 0) {
			continue;
		}



		sID.data[i+0] = tID.data[i+0];
		sID.data[i+1] = tID.data[i+1];
		sID.data[i+2] = tID.data[i+2];
	}

	sharingCvsCtx.putImageData(sID, 44, 0);

	imgBlob = elms.sharingCvs.toBlob(processBlob, "image/png");
}

let processBlob = function(blobFromCanvas) {
	let timestamp = ts.toString(16);
	let randomString = timestamp + "-";

	for(let i = 0; i < 8; i++) {
		let randomHexDigit = getRandomInt(16).toString(16);
		randomString += randomHexDigit;
	}

	codeToShare = randomString;

	fileToShare = new File([ blobFromCanvas ], randomString + ".png", {type: blobFromCanvas.type});

	let shareXhr = new XMLHttpRequest();
	shareXhr.open("PUT", storageUrl + fileToShare.name, true);
	shareXhr.onreadystatechange = function() {
		if(shareXhr.readyState === XMLHttpRequest.DONE) {
			let linkToSharedPic = "/shared.html?code=" + codeToShare;
			let sharingLink = document.getElementById("sharingLink");
			sharingLink.href = linkToSharedPic;
			sharingLink.style.display = "block";
			sharingLink.style.zIndex = 2049;
		}
	}

	shareXhr.send(fileToShare);
}

elms.sharingCvs.onclick = function() {
	elms.sharingCvs.width = 1;
	elms.sharingCvs.height = 1;
	elms.sharingCvs.style.zIndex = 0;
	elms.sharingCvs.style.display = "none";

	let sharingLink = document.getElementById("sharingLink");
	sharingLink.href = "#";
	sharingLink.style.display = "none";
	sharingLink.style.zIndex = 0;
}

document.getElementById("share_btn").onclick = sharing;


