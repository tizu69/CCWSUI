function estimatePx() {
	const dummy = document.createElement("div");
	dummy.style.width = "var(--px)";
	document.body.appendChild(dummy);
	const px = dummy.getBoundingClientRect().width;
	dummy.remove();
	return px;
}
function estimateTextWidth(text) {
	const dummy = document.createElement("span");
	// we need to use ZWSP to prevent the browser from collapsing whitespace
	dummy.textContent = "\u200B" + text.split("").join("\u200B") + "\u200B";
	document.body.appendChild(dummy);
	const width = dummy.getBoundingClientRect().width;
	dummy.remove();
	return width;
}

/**
 * @param {HTMLElement} container
 * @param {number} gridPx */
function snapToGrid(container, gridPx) {
	container.querySelectorAll("[data-ccwsui-snap] > *").forEach((el) => {
		const rect = el.getBoundingClientRect();
		let snapX = rect.left % gridPx;
		let snapY = rect.top % gridPx;
		if (snapX < 0) snapX += gridPx;
		if (snapY < 0) snapY += gridPx;
		el.style.transform = `translate(${-snapX}px, ${-snapY}px)`;
	});
}

const px = estimatePx();
document.body.addEventListener("htmx:after:process", (e) => {
	document.fonts.ready.then(() => snapToGrid(e.target, px));

	e.target.querySelectorAll('[data-ccwsui="inputregion"]').forEach((el) => {
		const caret = document.createElement("span");
		let interval;
		caret.style.position = "absolute";
		function update() {
			clearInterval(interval);
			caret.style.top = `calc(${el.offsetTop}px + var(--px))`;
			const text = el.value.substring(0, el.selectionStart);
			const isEnd = el.selectionStart === el.value.length;
			caret.style.left = el.offsetLeft + estimateTextWidth(text) + "px";
			el.parentElement.appendChild(caret);
			caret.textContent = isEnd ? "_" : "|";
			caret.style.visibility = "visible";
			if (isEnd) caret.dataset.ccwsuiShadow = true;
			else delete caret.dataset.ccwsuiShadow;
			interval = setInterval(() => {
				caret.style.visibility =
					caret.style.visibility === "hidden" ? "visible" : "hidden";
			}, 500);
		}
		el.addEventListener("focus", update);
		el.addEventListener("input", update);
		el.addEventListener("selectionchange", update);
		el.addEventListener("blur", () => {
			caret.remove();
			clearInterval(interval);
		});
	});
});
