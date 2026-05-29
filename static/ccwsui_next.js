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

	/** @type {HTMLCanvasElement} */
	canvas: document.getElementById("ccwsui-root"),

	mousePos: { x: -1, y: -1 },
	async prepare() {
		const font = new FontFace("CCWSUI", 'url("/static/font.ttf")');
		await font.load();
		document.fonts.add(font);

		this.canvas.width = window.innerWidth;
		this.canvas.height = window.innerHeight;

		window.addEventListener("resize", () => {
			this.canvas.width = window.innerWidth;
			this.canvas.height = window.innerHeight;
			this.totalRerender("Resized");
		});
		window.addEventListener("keydown", (e) => {
			if (!isDevtoolsShortcut(e)) return;
			e.preventDefault();
			this.openDevtools();
		});
		window.addEventListener("mousemove", (e) => {
			this.mousePos = {
				x: Math.floor(e.clientX / this.scale),
				y: Math.floor(e.clientY / this.scale),
			};
		});
		window.addEventListener(
			"wheel",
			(e) => {
				if (!e.ctrlKey) return;
				e.preventDefault();
				// zooming
				this.scale = Math.min(
					7,
					Math.max(1, this.scale + Math.sign(-e.deltaY)),
				);
				this.totalRerender("Zoomed");
			},
			{ passive: false },
		);
	},

	getDimensions() {
		return {
			w: Math.floor(this.canvas.width / this.scale),
			h: Math.floor(this.canvas.height / this.scale),
		};
	},
	getMousePos() {
		return this.mousePos;
	},

	clear() {
		const ctx = this.canvasContext;
		ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
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

	/** @type {Record<string, HTMLImageElement>} */
	textures: {},
	async prepareTextures(...path) {
		let anyAreNew = false;
		await Promise.all(
			path.map((p) => {
				if (this.textures[p]) return;
				anyAreNew = true;
				const img = new Image();
				if (p.startsWith("@item/"))
					img.src = `/static/item/${p.slice(6)}.png`;
				else img.src = `/static/tex/${p}.png`;
				return new Promise((resolve) => {
					img.onload = () => {
						this.textures[p] = img;
						resolve();
					};
					img.onerror = () => {
						console.error(`Failed to load texture: ${p}`);
						resolve();
					};
				});
			}),
		);
		if (anyAreNew) this.totalRerender("Textures have loaded");
	},
	renderTex(x, y, w, h, path, sx, sy, sw = w, sh = h, nn = true) {
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
	renderTexPattern(x, y, w, h, path, sx, sy, sw, sh) {
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

		if (Object.keys(this.renderTexPatternCache).length >= 4096)
			this.renderTexPatternCache = {};
		this.renderTexPatternCache[cacheKey] = pattern;

		ctx.fillStyle = pattern;
		ctx.fillRect(scrx, scry, this.scaled(w), this.scaled(h));
	},
	getTexSize(path) {
		const img = this.textures[path];
		if (!img) return { w: 0, h: 0 };
		return { w: img.width, h: img.height };
	},

	get canvasContext() {
		return this.canvas.getContext("2d");
	},
	scaled(v) {
		return Math.round(v) * this.scale;
	},

	devtoolsWindow: null,
	renderDevtools(reason, layoutJSON, tookLayout, tookDraw) {
		const layout = JSON.parse(layoutJSON);

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

		const dims = this.getDimensions();
		treeText += ind`Built at ${new Date().toTimeString().split(" ")[0]} `;
		treeText += `(${reason})\n`;
		treeText += ind`Last layout took ${(tookLayout / 1e6).toFixed(2)}ms\n`;
		treeText += ind`Last rerender took ${(tookDraw / 1e6).toFixed(2)}ms\n`;
		treeText += ind`Rendered at ${dims.w}x${dims.h}px`;
		treeText += ` (${(tookDraw / (dims.w * dims.h)).toFixed(0)}ns/px)\n\n`;
		treeIndent--;

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
};

const go = new Go();
WebAssembly.instantiateStreaming(
	fetch("/static/ccwsui.wasm"),
	go.importObject,
).then((result) => {
	go.run(result.instance);
});
