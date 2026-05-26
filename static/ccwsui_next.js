window.ccwsui = {
	scale: 3,
	lineheight: 9,
	fontsize: 12,

	socketURL:
		// AIR seems not support WebSocket, so we use the "real" port for it.
		`ws://${window.location.hostname}:${window.location.port == 8081 ? 8080 : window.location.port}` +
		document.getElementById("ccwsui-socketurl").textContent,

	totalRerender() {
		throw new Error("Go wasm not initialized!");
	},
	openDevtools() {
		throw new Error("Go wasm not initialized!");
	},

	/** @type {HTMLCanvasElement} */
	canvas: document.getElementById("ccwsui-root"),
	async prepare() {
		const font = new FontFace("CCWSUI", 'url("/static/font.ttf")');
		await font.load();
		document.fonts.add(font);

		this.canvas.width = window.innerWidth;
		this.canvas.height = window.innerHeight;

		window.addEventListener("resize", () => {
			this.canvas.width = window.innerWidth;
			this.canvas.height = window.innerHeight;
			this.totalRerender();
		});
		window.addEventListener("keydown", (e) => {
			const isDevtoolsShortcut =
				e.key === "F12" ||
				(e.ctrlKey && e.shiftKey && e.key.toLowerCase() === "i") ||
				(e.metaKey && e.altKey && e.key.toLowerCase() === "i");
			if (isDevtoolsShortcut) {
				e.preventDefault();
				this.openDevtools();
			}
		});
	},

	getDimensions() {
		return {
			x: Math.floor(this.canvas.width / this.scale),
			y: Math.floor(this.canvas.height / this.scale),
		};
	},
	clear() {
		const ctx = this.canvas.getContext("2d");
		ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
	},
	renderText(x, y, text) {
		const ctx = this.canvas.getContext("2d");
		ctx.fillStyle = "white";
		ctx.font = this.scaled(this.fontsize) + "px CCWSUI";
		ctx.fillText(text, this.scaled(x), this.scaled(y + this.lineheight));
	},
	guessTextWidth(text) {
		const ctx = this.canvas.getContext("2d");
		ctx.font = this.scaled(this.fontsize) + "px CCWSUI";
		const v = Math.round(ctx.measureText(text).width) / this.scale;
		if (v % 1 !== 0) console.warn("Text width is not integer!", v, text);
		return v;
	},

	textures: {},
	async prepareTextures(...path) {
		await Promise.all(
			path.map((p) => {
				const img = new Image();
				img.src = p;
				return new Promise((resolve) => {
					img.onload = () => {
						this.textures[p] = img;
						resolve();
					};
					img.onerror = resolve;
				});
			}),
		);
	},
	renderTex(x, y, w, h, path, sx, sy) {
		const ctx = this.canvas.getContext("2d");
		const img = this.textures[path];
		if (!img) return;
		ctx.imageSmoothingEnabled = false;
		ctx.drawImage(
			img,
			sx,
			sy,
			w,
			h,
			this.scaled(x),
			this.scaled(y),
			this.scaled(w),
			this.scaled(h),
		);
	},

	scaled(v) {
		return Math.round(v) * this.scale;
	},

	renderDevtools(layoutJSON) {
		this.closeDevtools();
		const layout = JSON.parse(layoutJSON);

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
		const nodes = [...(layout.Children || [])];
		const occupied = new Set();
		while (nodes.length > 0) {
			const node = nodes.pop();
			const { X, Y, W, H } = node.Rect;

			const color = `hsl(${(i++ / 10) * 360}, 70%, 40%)`;

			const overlay = document.createElement("div");
			overlay.style.position = "absolute";
			overlay.style.left = `${X * this.scale}px`;
			overlay.style.top = `${Y * this.scale}px`;
			overlay.style.width = `${W * this.scale}px`;
			overlay.style.height = `${H * this.scale}px`;
			overlay.style.border = "4px solid " + color;
			overlay.style.backgroundColor = `color-mix(in srgb, ${color} 20%, transparent)`;
			overlay.style.boxSizing = "border-box";
			overlay.style.zIndex = "9997";
			container.appendChild(overlay);

			let x = Math.ceil(X / 70) * 70;
			let y = Math.ceil(Y / 8) * 8;
			while (occupied.has(`${x},${y}`)) y += 8;
			occupied.add(`${x},${y}`);

			const title = document.createElement("p");
			title.style.position = "absolute";
			title.style.left = `${x * this.scale}px`;
			title.style.top = `${y * this.scale}px`;
			title.style.color = `white`;
			title.style.backgroundColor = color;
			title.style.padding = "2px";
			title.style.font = "12px monospace";
			title.textContent = node.Title;
			title.style.zIndex = "9998";
			container.appendChild(title);

			if (node.Children) nodes.push(...node.Children);
		}

		const panel = document.createElement("panel");
		panel.style.position = "fixed";
		panel.style.top = "0";
		panel.style.right = "0";
		panel.style.width = "5%";
		panel.style.height = "100%";
		panel.style.backgroundColor = "rgba(0, 0, 0, 0.5)";
		panel.style.zIndex = "9999";
		panel.style.overflow = "auto";
		panel.style.resize = "horizontal";
		container.appendChild(panel);

		const pre = document.createElement("pre");
		pre.style.color = "white";
		pre.style.fontFamily = "monospace";
		pre.style.fontSize = "12px";
		pre.style.padding = "10px";
		pre.textContent = JSON.stringify(layout, null, 2);
		panel.appendChild(pre);
	},
	closeDevtools() {
		document.getElementById("ccwsui-devtools")?.remove();
	},
};

const go = new Go();
WebAssembly.instantiateStreaming(
	fetch("/static/ccwsui.wasm"),
	go.importObject,
).then((result) => {
	go.run(result.instance);
});
