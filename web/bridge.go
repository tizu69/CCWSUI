//go:build wasm && js

package main

import (
	"fmt"
	"image/color"
	"syscall/js"
)

func (j *jsapi) GetDimensions() (w, h int) {
	ret := j.wrap.Call("getDimensions")
	return ret.Get("w").Int(), ret.Get("h").Int()
}

func (j *jsapi) GetMousePos() (x, y int) {
	ret := j.wrap.Call("getMousePos")
	return ret.Get("x").Int(), ret.Get("y").Int()
}

func (j *jsapi) GetMouseScroll() (dx, dy int) {
	ret := j.wrap.Get("mouseScroll")
	return ret.Get("dx").Int(), ret.Get("dy").Int()
}

func (j *jsapi) GetMouseDown() bool {
	return j.wrap.Get("mouseDown").Bool()
}

func (j *jsapi) GetShiftDown() bool {
	return j.wrap.Get("shiftDown").Bool()
}

func (j *jsapi) GetCtrlDown() bool {
	return j.wrap.Get("ctrlDown").Bool()
}

func (j *jsapi) GetAltDown() bool {
	return j.wrap.Get("altDown").Bool()
}

func (j *jsapi) Clear() {
	j.wrap.Call("clear")
}

func (j *jsapi) Scissor(x, y, w, h int) {
	j.wrap.Call("scissor", x, y, w, h)
}

func (j *jsapi) PopScissor() {
	j.wrap.Call("popScissor")
}

func (j *jsapi) RenderText(x, y int, text string, color color.Color) {
	r, g, b, a := color.RGBA()
	hex := fmt.Sprintf("#%02x%02x%02x%02x", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
	j.wrap.Call("renderText", x, y, text, hex)
}

func (j *jsapi) RequireTexture(path string) {
	j.requiredTextures = append(j.requiredTextures, path)
}

func (j *jsapi) PrepareTextures(wait bool, path ...string) {
	waitch := make(chan any)
	jspath := make([]any, len(path))
	for i, p := range path {
		jspath[i] = js.ValueOf(p)
	}
	j.wrap.Call("prepareTextures", jspath...).
		Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
			close(waitch)
			return nil
		}))
	if wait {
		<-waitch
	}
}

func (j *jsapi) RenderTex(x, y, w, h int, path string, sx, sy int) {
	j.wrap.Call("renderTex", x, y, w, h, path, sx, sy)
}

func (j *jsapi) RenderTexMixel(x, y, w, h int, path string, sx, sy, sw, sh int, nn bool) {
	j.wrap.Call("renderTex", x, y, w, h, path, sx, sy, sw, sh, nn)
}

func (j *jsapi) RenderTexPattern(x, y, w, h int, path string, sx, sy, sw, sh int) {
	j.wrap.Call("renderTexPattern", x, y, w, h, path, sx, sy, sw, sh)
}

func (j *jsapi) GetTexSize(path string) (w, h int) {
	ret := j.wrap.Call("getTexSize", path)
	return ret.Get("w").Int(), ret.Get("h").Int()
}

func (j *jsapi) GetTexBorders(path string) (t, l, b, r int) {
	if borders, ok := j.texBordersCache[path]; ok {
		return borders[0], borders[1], borders[2], borders[3]
	}
	ret := j.wrap.Get("textureBorders").Get(path)
	if ret.IsUndefined() {
		return 0, 0, 0, 0
	}
	j.texBordersCache[path] = [4]int{
		ret.Get("t").Int(), ret.Get("l").Int(),
		ret.Get("b").Int(), ret.Get("r").Int(),
	}
	return ret.Get("t").Int(), ret.Get("l").Int(),
		ret.Get("b").Int(), ret.Get("r").Int()
}

func (j *jsapi) GuessTextWidth(text string) int {
	if width, ok := j.textWidthCache[text]; ok {
		return width
	}
	width := j.wrap.Call("guessTextWidth", text).Int()
	j.textWidthCache[text] = width
	return width
}

func (j *jsapi) GetLineHeight() int {
	return j.wrap.Get("lineheight").Int()
}

func (j *jsapi) RenderOverlay(fn func()) {
	j.overlays = append(j.overlays, fn)
}
