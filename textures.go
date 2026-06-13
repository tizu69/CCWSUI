package main

import (
	"bytes"
	"image/png"
)

var coreTextures = map[string][]byte{}

func init() {
	tex, _ := staticFS.ReadDir("static/tex")
	for _, tex := range tex {
		data, err := staticFS.ReadFile("static/tex/" + tex.Name())
		if err != nil {
			panic(err)
		}
		data, err = stripColorProfile(data)
		if err != nil {
			panic(err)
		}
		coreTextures[tex.Name()[:len(tex.Name())-4]] = data
	}

	icons, _ := staticFS.ReadDir("static/icon")
	for _, icon := range icons {
		data, err := staticFS.ReadFile("static/icon/" + icon.Name())
		if err != nil {
			panic(err)
		}
		data, err = stripColorProfile(data)
		if err != nil {
			panic(err)
		}
		coreTextures["@icon/"+icon.Name()[:len(icon.Name())-4]] = data
	}
}

func stripColorProfile(data []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
