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
		console.log("Fonts ready?", document.fonts.status);
		console.log(
			"Images loaded?",
			[...document.images].every((img) => img.complete),
		);
		const offsetX = el.offsetLeft % gridPx;
		const offsetY = el.offsetTop % gridPx;
		console.log(el, offsetX, offsetY);
		el.style.transform = `translate(${-offsetX}px, ${-offsetY}px)`;
	});
}

const px = estimatePx();
document.fonts.ready.then(() => snapToGrid(document.body, px));
document.body.addEventListener("htmx:afterSwap", (e) => {
	document.fonts.ready.then(() => snapToGrid(e.detail.target, px));
});
