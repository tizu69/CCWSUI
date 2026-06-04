//go:build wasm && js

package main

import (
	"encoding/json"
	"log/slog"
	"time"

	"g.tizu.dev/CCWSUI/components"
)

func (j *jsapi) TotalRerender(reason string) {
	slog.Info("Rerendering!", "reason", reason)
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
		b2, _ := json.Marshal(j.context)
		j.wrap.Call("renderDevtools", reason, string(b), tookLayout, tookRender, string(b2))
	}
}
