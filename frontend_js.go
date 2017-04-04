package main

const (
	frontendJs = "// Globals\nvar SCALE = 20;\n\nvar CANVAS = null;\nvar LOADING = null;\nvar PICKER = null;\n\nvar CONTEXT = null;\nvar CONNECTION = null;\n\nvar MY_COLOR = {r: 255, g: 0, b: 0};\nvar ON_COOLDOWN = false;\n\n// Initialize\nwindow.onload = function () {\n  CANVAS = document.getElementById(\"canvas\");\n  LOADING = document.getElementById(\"loading\");\n  PICKER = document.getElementById(\"picker\");\n\n  setupCanvas();\n  setupWebsocket();\n  setupPicker();\n};\n\nfunction setupCanvas () {\n  CONTEXT = CANVAS.getContext(\"2d\");\n\n  CANVAS.width = SIZE;\n  CANVAS.height = SIZE;\n  CANVAS.style.transform = \"scale(\" + SCALE + \")\";\n\n  CANVAS.onclick = function (e) {\n    var x = Math.floor(e.pageX / SCALE);\n    var y = Math.floor(e.pageY / SCALE);\n    if (!ON_COOLDOWN) {\n      sendPixel(x, y, MY_COLOR.r, MY_COLOR.g, MY_COLOR.b);\n      toggleCooldown();\n    }\n  };\n\n  CANVAS.oncontextmenu = function (e) {\n    e.preventDefault();\n  };\n}\n\nfunction toggleCooldown () {\n  if (ON_COOLDOWN) { return; }\n  if (COOLDOWN <= 0) { return; }\n\n  ON_COOLDOWN = true;\n\n  var el = document.getElementById(\"cooldown\");\n\n  var cooldownMs = COOLDOWN * 1000;\n  var cooldownEnd = Date.now() + cooldownMs;\n\n  el.style.display = \"block\";\n  el.textContent = cooldownMs + \"ms\";\n\n  var intervalId = setInterval(() => {\n    el.textContent = (cooldownEnd - Date.now()) + \"ms\";\n  }, 5);\n\n  setTimeout(() => {\n    ON_COOLDOWN = false;\n    clearInterval(intervalId);\n    el.style.display = \"none\";\n  }, cooldownMs);\n}\n\nfunction setupWebsocket () {\n  function processCmd (data) {\n    var parts = data.split(\" \");\n    if (parts.length === 5) {\n      setPixel(parts[0], parts[1], parts[2], parts[3], parts[4]);\n    }\n  }\n\n  function fillCanvas (data) {\n    var dv = new DataView(data);\n    var image = CONTEXT.createImageData(SIZE, SIZE);\n\n    var i, j, k;\n    for (i = 0; i < (dv.byteLength / 3); i += 1) {\n      j = i*3;\n      k = i*4;\n      image.data[k] = dv.getUint8(j);\n      image.data[k+1] = dv.getUint8(j+1);\n      image.data[k+2] = dv.getUint8(j+2);\n      image.data[k+3] = 255;\n    }\n\n    CONTEXT.putImageData(image, 0, 0);\n\n    LOADING.style.display = \"none\";\n    CANVAS.style.display = \"block\";\n    PICKER.style.display = \"block\";\n  }\n\n  var url = \"ws://\" + window.location.host + \"/ws\";\n\n  CONNECTION = new WebSocket(url);\n  CONNECTION.binaryType = \"arraybuffer\";\n\n  CONNECTION.onmessage = function (e) {\n    var data = e.data;\n\n    if (typeof data === \"string\") {\n      return processCmd(data);\n    }\n\n    if (data instanceof ArrayBuffer) {\n      return fillCanvas(data);\n    }\n  };\n}\n\nfunction setupPicker () {\n  var input = document.getElementById(\"picker-input\");\n  var preview = document.getElementById(\"picker-preview\");\n\n  preview.style.backgroundColor = \"#FF0000\";\n  input.value = \"#FF0000\";\n\n  input.onchange = function (e) {\n    var value = e.target.value;\n    var color = hexToRgb(value);\n\n    if (color) {\n      MY_COLOR = color;\n      preview.style.backgroundColor = value;\n    }\n  };\n}\n\n// Canvas methods.\nfunction sendPixel (x, y, r, g, b) {\n  CONNECTION.send([x, y, r, g, b].join(\" \"));\n}\n\nfunction setPixel (x, y, r, g, b) {\n  if (!this.id) {\n    this.id = CONTEXT.createImageData(1, 1);\n    this.idd = this.id.data;\n    this.idd[3] = 255;\n  }\n\n  this.idd[0] = r;\n  this.idd[1] = g;\n  this.idd[2] = b;\n\n  CONTEXT.putImageData(this.id, x, y);\n}\n\n// Util\nfunction hexToRgb (hex) {\n\tvar shorthandRegex = /^#?([a-f\\d])([a-f\\d])([a-f\\d])$/i;\n\n\thex = hex.replace(shorthandRegex, function(m, r, g, b) {\n\t\treturn r + r + g + g + b + b;\n\t});\n\n\tvar result = /^#?([a-f\\d]{2})([a-f\\d]{2})([a-f\\d]{2})$/i.exec(hex);\n\n\treturn result ? {\n\t\tr: parseInt(result[1], 16),\n\t\tg: parseInt(result[2], 16),\n\t\tb: parseInt(result[3], 16)\n\t} : null;\n}\n"
)