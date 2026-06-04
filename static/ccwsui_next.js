const isDevtoolsShortcut = (e) =>
	e.key === "F12" ||
	(e.ctrlKey && e.shiftKey && e.key.toLowerCase() === "i") ||
	(e.metaKey && e.altKey && e.key.toLowerCase() === "i");

window.ccwsui = {
	scale: 3,
	baseline: 8,
	lineheight: 9,
	fontsize: 12,

	socketURL:
		// AIR seems not support WebSocket, so we use the "real" port for it.
		`ws://${window.location.hostname}:${window.location.port == 8081 ? 8080 : window.location.port}` +
		document.getElementById("ccwsui-socketurl").textContent,

	totalRerender(reason) {
		throw new Error("Go wasm not initialized!");
	},
	openDevtools() {
		throw new Error("Go wasm not initialized!");
	},

	_rerenderReason: "Initial render",
	queueRerender(reason) {
		if (!this._rerenderReason) this._rerenderReason = reason;
	},
	_lastRerenderTime: 0,
	_rerenderAnimationFrame() {
		requestAnimationFrame(this._rerenderAnimationFrame.bind(this));

		if (!this._rerenderReason) return;
		this.totalRerender(this._rerenderReason);
		this._rerenderReason = null;

		if (this.mouseScroll.dx || this.mouseScroll.dy)
			this.mouseScroll = { dx: 0, dy: 0 };
		this.mouseDown = false;
	},

	/** @type {HTMLCanvasElement} */
	canvas: document.getElementById("ccwsui-root"),

	mousePos: { x: -1, y: -1 },
	mouseScroll: { dx: 0, dy: 0 },
	mouseDown: false,
	keysDown: new Set(),
	shiftDown: false,
	ctrlDown: false,
	altDown: false,
	async prepare() {
		const font = new FontFace("CCWSUI", 'url("/static/font.ttf")');
		await font.load();
		document.fonts.add(font);

		this.canvas.width = window.innerWidth;
		this.canvas.height = window.innerHeight;

		window.addEventListener("resize", () => {
			this.canvas.width = window.innerWidth;
			this.canvas.height = window.innerHeight;
			this.ctx = this.canvas.getContext("2d");
			this.queueRerender("Resized");
		});
		window.addEventListener("keydown", (e) => {
			this.keysDown.add(e.key);
			this.shiftDown = e.shiftKey;
			this.ctrlDown = e.ctrlKey;
			this.altDown = e.altKey;
			this.queueRerender("Key pressed");
			if (!isDevtoolsShortcut(e)) return;
			e.preventDefault();
			this.openDevtools();
		});
		window.addEventListener("keyup", (e) => {
			this.keysDown.delete(e.key);
			this.shiftDown = e.shiftKey;
			this.ctrlDown = e.ctrlKey;
			this.altDown = e.altKey;
			this.queueRerender("Key released");
		});
		this.canvas.addEventListener("mousemove", (e) => {
			this.mousePos = {
				x: Math.floor(e.clientX / this.scale),
				y: Math.floor(e.clientY / this.scale),
			};
			this.queueRerender("Mouse moved");
		});
		this.canvas.addEventListener("mouseout", () => {
			this.mousePos = { x: -1, y: -1 };
			this.queueRerender("Mouse left");
		});
		this.canvas.addEventListener(
			"wheel",
			(e) => {
				if (e.ctrlKey) {
					e.preventDefault();
					// zooming
					this.scale = Math.min(
						7,
						Math.max(1, this.scale + Math.sign(-e.deltaY)),
					);
					this.queueRerender("Zoomed");
					return;
				}

				this.mouseScroll = { dx: e.deltaX, dy: e.deltaY };
				this.queueRerender("Scrolled");
			},
			{ passive: false },
		);
		this.canvas.addEventListener("mousedown", (e) => {
			this.mouseDown = true;
			this.queueRerender("Mouse down");
		});

		requestAnimationFrame(this._rerenderAnimationFrame.bind(this));
	},

	getDimensions() {
		return {
			w: Math.floor(this.canvas.width / this.scale),
			h: Math.floor(this.canvas.height / this.scale),
		};
	},
	getMousePos() {
		if (this.scissorStack) {
			let valid = true;
			for (const scissor of this.scissorStack)
				if (
					this.mousePos.x < scissor.x ||
					this.mousePos.x > scissor.x + scissor.w ||
					this.mousePos.y < scissor.y ||
					this.mousePos.y > scissor.y + scissor.h
				) {
					valid = false;
					break;
				}
			if (!valid) return { x: -1, y: -1 };
		}
		return this.mousePos;
	},

	scissorStack: [],
	scissor(x, y, w, h) {
		this.canvasContext.save();
		this.canvasContext.beginPath();
		this.canvasContext.rect(
			this.scaled(x),
			this.scaled(y),
			this.scaled(w),
			this.scaled(h),
		);
		this.canvasContext.clip();
		this.scissorStack.push({ x, y, w, h });
	},
	popScissor() {
		if (!this.scissorStack.length) throw new Error("No scissor to pop");
		this.canvasContext.restore();
		this.scissorStack.pop();
	},

	clear() {
		const ctx = this.canvasContext;
		ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
	},
	isRectVisible(x, y, w, h) {
		if (w <= 0 || h <= 0) return false;
		const maxW = Math.floor(this.canvas.width / this.scale);
		const maxH = Math.floor(this.canvas.height / this.scale);
		if (x >= maxW || y >= maxH || x + w <= 0 || y + h <= 0) return false;
		for (const scissor of this.scissorStack)
			if (
				x >= scissor.x + scissor.w ||
				y >= scissor.y + scissor.h ||
				x + w <= scissor.x ||
				y + h <= scissor.y
			)
				return false;
		return true;
	},
	renderText(x, y, text, color) {
		const ctx = this.canvasContext;
		ctx.fillStyle = color;
		ctx.font = this.scaled(this.fontsize) + "px CCWSUI";
		ctx.fillText(text, this.scaled(x), this.scaled(y + this.baseline));
	},
	guessTextWidth(text) {
		const ctx = this.canvasContext;
		ctx.font = this.scaled(this.fontsize) + "px CCWSUI";
		const v = Math.round(ctx.measureText(text).width) / this.scale;
		if (v % 1 !== 0) console.warn("Text width is not integer!", v, text);
		return v;
	},

	/** @type {Record<string, HTMLCanvasElement>} */
	textures: {},
	/** @type {Record<string, string>} */
	userTextures: {},
	/** @type {Record<string, {t: number,  l: number, b: number, r: number}>} */
	textureBorders: {},
	async prepareTextures(...path) {
		let anyAreNew = false;
		await Promise.all(
			path.map((p) => {
				if (this.textures[p]) return;
				anyAreNew = true;

				let flags = {};
				let [path, ...flagList] = p.split(";");
				for (const flag of flagList) {
					const [key, value] = flag.split("=");
					flags[key] = value || true;
				}

				if (path.startsWith("#")) {
					const canvas = document.createElement("canvas");
					canvas.width = 1;
					canvas.height = 1;
					const ctx = canvas.getContext("2d");
					ctx.fillStyle = path;
					ctx.fillRect(0, 0, 1, 1);
					this.textures[p] = canvas;
					return Promise.resolve();
				}

				const img = new Image();
				if (path in this.userTextures)
					img.src = `data:image/png;base64,${this.userTextures[path]}`;
				else if (path.startsWith("@item/")) {
					img.src = `/static/item/${path.slice(6)}.png`;
					flags = {}; // TODO: do we want items to be processed?
				} else if (path.startsWith("@icon/"))
					img.src = `/static/icon/${path.slice(6)}.png`;
				else img.src = `/static/tex/${path}.png`;

				const canvas = document.createElement("canvas");
				this.textures[p] = canvas;
				return new Promise((resolve) => {
					img.onload = () => {
						canvas.width = p.startsWith("@item/")
							? img.width // icons don't need a sidecar texture
							: img.width / 2;
						canvas.height = img.height;

						const proccanvas = document.createElement("canvas");
						proccanvas.width = canvas.width;
						proccanvas.height = canvas.height;
						const ctx = proccanvas.getContext("2d");
						ctx.imageSmoothingEnabled = false;
						ctx.drawImage(img, 0, 0);

						const sccanvas = document.createElement("canvas");
						sccanvas.width = img.width / 2;
						sccanvas.height = img.height;
						const scctx = sccanvas.getContext("2d");
						scctx.imageSmoothingEnabled = false;
						scctx.drawImage(img, -img.width / 2, 0);
						const scdata = scctx.getImageData(
							0,
							0,
							sccanvas.width,
							sccanvas.height,
						).data;

						for (const src in flags) {
							if (!src.startsWith("#")) continue;
							const dst = flags[src];

							const srcColor = this.hexToRgba(src);
							const dstColor = this.hexToRgba(dst);
							const imageData = ctx.getImageData(
								0,
								0,
								canvas.width,
								canvas.height,
							);
							const data = imageData.data;
							for (let i = 0; i < data.length; i += 4) {
								if (
									data[i] === srcColor.r &&
									data[i + 1] === srcColor.g &&
									data[i + 2] === srcColor.b &&
									data[i + 3] === srcColor.a
								) {
									data[i] = dstColor.r;
									data[i + 1] = dstColor.g;
									data[i + 2] = dstColor.b;
									data[i + 3] = dstColor.a;
								}
							}
							ctx.putImageData(imageData, 0, 0);
						}
						if (flags["tint"]) {
							const { r, g, b } = this.hexToRgba(flags["tint"]);
							const imageData = ctx.getImageData(
								0,
								0,
								canvas.width,
								canvas.height,
							);
							const data = imageData.data;
							for (let i = 0; i < data.length; i += 4) {
								// red channel defines tint strength
								// 1-126 -> darker
								// 127 -> exact same color
								// 128-255 -> lighter
								const v = scdata[i] / 127;
								if (v === 0) continue;
								if (v < 1) {
									data[i] = r * v;
									data[i + 1] = g * v;
									data[i + 2] = b * v;
								} else {
									data[i] = r + (255 - r) * (v - 1);
									data[i + 1] = g + (255 - g) * (v - 1);
									data[i + 2] = b + (255 - b) * (v - 1);
								}
							}

							ctx.putImageData(imageData, 0, 0);
						}
						if (flags["shadow"]) {
							const imageData = ctx.getImageData(
								0,
								0,
								canvas.width,
								canvas.height,
							);
							const data = imageData.data;
							for (let i = 0; i < data.length; i += 4) {
								data[i] = data[i] / 4;
								data[i + 1] = data[i + 1] / 4;
								data[i + 2] = data[i + 2] / 4;
							}
							ctx.putImageData(imageData, 0, 0);
						}

						let t = 0,
							l = 0,
							b = 0,
							r = 0;
						const borderpoints = [];
						for (let i = 0; i < scdata.length; i += 4)
							if (scdata[i + 1] & 64)
								borderpoints.push({
									x: (i / 4) % canvas.width,
									y: Math.floor(i / 4 / canvas.width),
								});
						if (borderpoints[0]) {
							t = borderpoints[0].y;
							l = borderpoints[0].x;
							b = canvas.height - borderpoints[0].y - 1;
							r = canvas.width - borderpoints[0].x - 1;
						}
						if (borderpoints[1]) {
							t = Math.min(t, borderpoints[1].y);
							l = Math.min(l, borderpoints[1].x);
							b = Math.min(
								b,
								canvas.height - borderpoints[1].y - 1,
							);
							r = Math.min(
								r,
								canvas.width - borderpoints[1].x - 1,
							);
						}
						if (!path.startsWith("@item/"))
							this.textureBorders[p] = { t, l, r, b };

						const fctx = canvas.getContext("2d");
						fctx.imageSmoothingEnabled = false;
						fctx.save();
						fctx.translate(canvas.width / 2, canvas.height / 2);
						if (flags["rotate"]) {
							const angle = (+flags["rotate"] * Math.PI) / 180;
							fctx.rotate(angle);
						}
						if (flags["flip"] === "x") {
							fctx.scale(-1, 1);
						} else if (flags["flip"] === "y") {
							fctx.scale(1, -1);
						}
						fctx.drawImage(
							proccanvas,
							-proccanvas.width / 2,
							-proccanvas.height / 2,
						);
						fctx.restore();

						resolve();
					};
					img.onerror = () => {
						console.error(`Failed to load texture: ${p}`);
						resolve();
					};
				});
			}),
		);
		if (anyAreNew) this.queueRerender("Textures have loaded");
	},
	renderTex(x, y, w, h, path, sx, sy, sw = w, sh = h, nn = true) {
		if (!this.isRectVisible(x, y, w, h)) return;
		const ctx = this.canvasContext;
		const img = this.textures[path];
		if (!img) return;

		ctx.imageSmoothingEnabled = !nn;
		ctx.drawImage(
			img,
			sx,
			sy,
			sw,
			sh,
			this.scaled(x),
			this.scaled(y),
			this.scaled(w),
			this.scaled(h),
		);
	},
	/** @type {Record<string, CanvasPattern>} */
	renderTexPatternCache: {},
	renderTexPatternCacheSize: 0,
	renderTexPattern(x, y, w, h, path, sx, sy, sw, sh) {
		if (!this.isRectVisible(x, y, w, h)) return;
		const ctx = this.canvasContext;
		const img = this.textures[path];
		if (!img) return;
		ctx.imageSmoothingEnabled = false;

		const cacheKey = `${path}:${sx}:${sy}:${sw}:${sh}:${this.scale}`;
		if (this.renderTexPatternCache[cacheKey]) {
			const pattern = this.renderTexPatternCache[cacheKey];
			pattern.setTransform(
				new DOMMatrix([1, 0, 0, 1, this.scaled(x), this.scaled(y)]),
			);
			ctx.fillStyle = pattern;
			ctx.fillRect(
				this.scaled(x),
				this.scaled(y),
				this.scaled(w),
				this.scaled(h),
			);
			return;
		}

		const subc = document.createElement("canvas");
		subc.width = this.scaled(sw);
		subc.height = this.scaled(sh);
		const subctx = subc.getContext("2d");
		subctx.imageSmoothingEnabled = false;
		subctx.drawImage(img, sx, sy, sw, sh, 0, 0, subc.width, subc.height);

		const scrx = this.scaled(x);
		const scry = this.scaled(y);
		const pattern = ctx.createPattern(subc, "repeat");
		pattern.setTransform(new DOMMatrix([1, 0, 0, 1, scrx, scry]));

		if (this.renderTexPatternCacheSize >= 4096) {
			this.renderTexPatternCache = {};
			this.renderTexPatternCacheSize = 0;
		}
		this.renderTexPatternCache[cacheKey] = pattern;
		this.renderTexPatternCacheSize++;

		ctx.fillStyle = pattern;
		ctx.fillRect(scrx, scry, this.scaled(w), this.scaled(h));
	},
	getTexSize(path) {
		const img = this.textures[path];
		if (!img) return { w: 1, h: 1 };
		return { w: img.width, h: img.height };
	},

	get canvasContext() {
		return this.ctx || this.canvas.getContext("2d");
	},
	scaled(v) {
		return Math.round(v) * this.scale;
	},

	devtoolsWindow: null,
	tookLayoutValues: [],
	tookDrawValues: [],
	renderDevtools(reason, layoutJSON, tookLayout, tookDraw, cctxJSON) {
		const layout = JSON.parse(layoutJSON);
		const cctx = JSON.parse(cctxJSON);
		this.tookLayoutValues.push(tookLayout);
		this.tookDrawValues.push(tookDraw);

		document.getElementById("ccwsui-devtools")?.remove();

		const container = document.createElement("div");
		container.id = "ccwsui-devtools";
		container.style.position = "fixed";
		container.style.top = "0";
		container.style.right = "0";
		container.style.bottom = "0";
		container.style.left = "0";
		container.style.zIndex = "9996";
		document.body.appendChild(container);

		let i = 0;
		const nodes = [layout];
		while (nodes.length > 0) {
			const node = nodes.shift();
			const { X, Y, W, H } = node.Rect;

			const color = `hsl(${(i++ * 137) % 360}deg, 70%, 40%)`;

			const overlay = document.createElement("div");
			overlay.classList.add("ccwsui-layout-tree-overlay");
			overlay.style.left = `${X * this.scale}px`;
			overlay.style.top = `${Y * this.scale}px`;
			overlay.style.width = `${W * this.scale}px`;
			overlay.style.height = `${H * this.scale}px`;
			overlay.style.border = "3px solid " + color;
			overlay.style.backgroundColor = `color-mix(in srgb, ${color} 20%, transparent)`;
			if (node._parent) container.appendChild(overlay);

			const title = document.createElement("p");
			title.classList.add("ccwsui-layout-tree-title");
			title.style.left = `${X * this.scale}px`;
			title.style.top = `${Y * this.scale}px`;
			title.style.backgroundColor = color;
			let titletext = node.Title;
			for (let parent = node._parent; parent; parent = parent._parent)
				titletext += `\nin ${parent.Title}`;
			title.textContent = titletext;
			container.appendChild(title);

			for (const child of node.Children || []) {
				child._parent = node;
				nodes.push(child);
			}
		}

		let treeText = "CCWSUI Layout Tree\n";
		let treeIndent = 1;
		const ind = (s, ...args) =>
			"  ".repeat(treeIndent) +
			s.reduce((a, b, i) => a + args[i - 1] + b);

		const dims = this.getDimensions();
		treeText += ind`Built ${i} nodes at ${new Date().toTimeString().split(" ")[0]} `;
		treeText += `(${reason})\n`;
		treeText += ind`Last layout took ${(tookLayout / 1e6).toFixed(2)}ms `;
		treeText += `(~${(this.tookLayoutValues.reduce((a, b) => a + b) / this.tookLayoutValues.length / 1e6).toFixed(2)}ms)\n`;
		treeText += ind`Last rerender took ${(tookDraw / 1e6).toFixed(2)}ms `;
		treeText += `(~${(this.tookDrawValues.reduce((a, b) => a + b) / this.tookDrawValues.length / 1e6).toFixed(2)}ms)\n`;
		treeText += ind`Rendered at ${dims.w}x${dims.h}px`;
		treeText += ` (${(tookDraw / (dims.w * dims.h)).toFixed(0)}ns/px)\n\n`;
		treeIndent--;

		treeText += "Context " + JSON.stringify(cctx, null, 2) + "\n\n";
		function printTree(node) {
			treeText += ind`${node.Title} {\n`;
			treeIndent++;
			treeText += ind`x=${node.Rect.X}, y=${node.Rect.Y}, `;
			treeText += `w=${node.Rect.W}, h=${node.Rect.H}\n`;
			if (node.Children)
				for (const child of node.Children) {
					treeText += "\n";
					printTree(child);
				}
			treeIndent--;
			treeText += ind`}\n`;
		}
		printTree(layout);

		if (!this.devtoolsWindow || this.devtoolsWindow.closed) {
			this.devtoolsWindow = window.open(
				"",
				"ccwsui-devtools",
				"width=400,height=600",
			);

			this.devtoolsWindow.addEventListener(
				"beforeunload",
				this.openDevtools,
			);
			this.devtoolsWindow.addEventListener("keydown", (e) => {
				if (!isDevtoolsShortcut(e)) return;
				e.preventDefault();
				this.closeDevtools();
			});

			this.devtoolsWindow.document.title = "ccwsui devtools";
		}

		const tree = document.createElement("pre");
		tree.id = "ccwsui-devtools-tree";
		tree.style.color = "#cdd6f4";
		tree.style.fontFamily = "monospace";
		tree.style.fontSize = "12px";
		tree.style.padding = "10px";
		tree.textContent = treeText;
		this.devtoolsWindow.document
			.getElementById("ccwsui-devtools-tree")
			?.remove();
		this.devtoolsWindow.document.body.appendChild(tree);
		this.devtoolsWindow.document.documentElement.style.backgroundColor =
			"#1e1e2e";
	},
	closeDevtools() {
		document.getElementById("ccwsui-devtools")?.remove();
		this.devtoolsWindow?.close();
	},

	hexToRgba(hex) {
		const r = parseInt(hex.slice(1, 3), 16);
		const g = parseInt(hex.slice(3, 5), 16);
		const b = parseInt(hex.slice(5, 7), 16);
		const a = parseInt(hex.slice(7, 9), 16) || 255;
		return { r, g, b, a };
	},
};

const go = new Go();
WebAssembly.instantiateStreaming(
	fetch("/static/ccwsui.wasm"),
	go.importObject,
).then((result) => {
	go.run(result.instance);
});
