function estimatePx() {
	const dummy = document.createElement("div");
	dummy.style.width = "var(--px)";
	document.body.appendChild(dummy);
	const px = dummy.getBoundingClientRect().width;
	dummy.remove();
	return px;
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
});
