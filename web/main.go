//go:build wasm && js

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"syscall/js"
	"time"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
	"github.com/Marlliton/slogpretty"
	"github.com/coder/websocket"
)

type jsapi struct {
	wrap             js.Value
	Root             components.Native
	DevTools         bool
	requiredTextures []string
	overlays         []func()
	textWidthCache   map[string]int
	textWidthScale   int
}

func NewJSAPI() *jsapi {
	j := &jsapi{wrap: js.Global().Get("ccwsui")}

	j.wrap.Set("totalRerender", js.FuncOf(j.totalRerender))
	j.wrap.Set("openDevtools", js.FuncOf(j.openDevtools))

	wait := make(chan any)
	j.wrap.Call("prepare").Call("then", js.FuncOf(func(this js.Value, args []js.Value) any {
		close(wait)
		return nil
	}))
	<-wait

	return j
}

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
	hexcolor := fmt.Sprintf("#%02x%02x%02x%02x", uint8(r>>8),
		uint8(g>>8), uint8(b>>8), uint8(a>>8))
	j.wrap.Call("renderText", x, y, text, hexcolor)
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

func (j *jsapi) totalRerender(this js.Value, args []js.Value) any {
	j.TotalRerender(args[0].String())
	return nil
}

func (j *jsapi) RenderOverlay(fn func()) {
	j.overlays = append(j.overlays, fn)
}

func (j *jsapi) TotalRerender(reason string) {
	slog.Info("Rerendering!")
	j.Clear()
	if j.Root == nil {
		return
	}

	t := time.Now()
	w, h := j.GetDimensions()
	scale := j.wrap.Get("scale").Int()
	if scale != j.textWidthScale {
		j.textWidthScale = scale
		j.textWidthCache = map[string]int{}
	}
	l := j.Root.Layout(j, components.Rect{X: 0, Y: 0, W: w, H: h})
	tookLayout := time.Since(t).Nanoseconds()

	t = time.Now()
	j.Root.Render(j, l)
	for _, fn := range j.overlays {
		fn()
	}
	j.overlays = nil
	tookRender := time.Since(t).Nanoseconds()

	if j.requiredTextures != nil {
		j.PrepareTextures(false, j.requiredTextures...)
		j.requiredTextures = nil
	}

	if j.DevTools {
		b, _ := json.Marshal(l)
		j.wrap.Call("renderDevtools", reason, string(b), tookLayout, tookRender)
	}
}

func pathKey(path []int) string {
	if len(path) == 0 {
		return ""
	}
	parts := make([]string, len(path))
	for i, p := range path {
		parts[i] = strconv.Itoa(p)
	}
	return strings.Join(parts, ".")
}

func (j *jsapi) openDevtools(this js.Value, args []js.Value) any {
	j.DevTools = !j.DevTools
	j.TotalRerender("DevTools toggled")
	if !j.DevTools {
		j.wrap.Call("closeDevtools")
	}
	return nil
}

func (j *jsapi) SocketURL() string {
	return j.wrap.Get("socketURL").String()
}

func main() {
	slog.SetDefault(slog.New(slogpretty.New(os.Stdout, &slogpretty.Options{
		Level: slog.LevelDebug, AddSource: true, Colorful: true,
		Multiline: true, TimeFormat: "01.02.06 3:04PM",
	})))

	j := NewJSAPI()

	slog.Info("Connecting to gateway!", "url", j.SocketURL())
	ws, _, err := websocket.Dial(context.Background(), j.SocketURL(), &websocket.DialOptions{})
	if err != nil {
		panic(err)
	}
	defer ws.CloseNow()
	ws.SetReadLimit(1024 * 1024) // 1MB

	go func() {
		for {
			_, msg, err := ws.Read(context.Background())
			if err != nil {
				panic(err)
			}
			slog.Info("Received message")

			var e webmsg.Envelope
			if err := json.Unmarshal(msg, &e); err != nil {
				slog.Error("Failed to unmarshal message", "err", err)
				continue
			}

			switch e.Type {
			case webmsg.TypeUpdate:
				var u webmsg.Update
				if err := json.Unmarshal(e.Data, &u); err != nil {
					slog.Error("Failed to unmarshal update", "err", err)
					continue
				}
				slog.Info("Update", "data", string(msg))
				newroot, err := components.FromWire(u.Root)
				if err != nil {
					slog.Error("Failed to unmarshal root", "err", err)
					continue
				}
				j.Root = newroot
				j.TotalRerender("Server-sent Tree Update")
			}
		}
	}()

	select {}
}
