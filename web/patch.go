//go:build wasm && js

package main

import (
	"strconv"
	"strings"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
)

func applyPatches(root *components.WireNode, patches []webmsg.PatchEntry) {
	for _, p := range patches {
		switch p.Action {
		case webmsg.PatchSetProps:
			if p.Path == "" {
				root.Props = p.Props
			} else if node := nodeAtPath(root, p.Path); node != nil {
				node.Props = p.Props
			}
		case webmsg.PatchReplace:
			if p.Path == "" && p.Node != nil {
				*root = *p.Node
			} else if parent, idx := parentAndChild(root, p.Path); parent != nil && idx >= 0 && idx < len(parent.Children) && p.Node != nil {
				parent.Children[idx] = *p.Node
			}
		case webmsg.PatchRemove:
			if parent, idx := parentAndChild(root, p.Path); parent != nil && idx >= 0 && idx < len(parent.Children) {
				parent.Children = append(parent.Children[:idx], parent.Children[idx+1:]...)
			}
		case webmsg.PatchInsert:
			if parent, idx := parentAndChild(root, p.Path); parent != nil && p.Node != nil {
				if idx >= 0 && idx <= len(parent.Children) {
					parent.Children = append(parent.Children[:idx], append([]components.WireNode{*p.Node}, parent.Children[idx:]...)...)
				} else {
					parent.Children = append(parent.Children, *p.Node)
				}
			}
		}
	}
}

func parsePath(path string) []int {
	if path == "" {
		return nil
	}
	parts := strings.Split(path, ".")
	indices := make([]int, len(parts))
	for i, p := range parts {
		indices[i], _ = strconv.Atoi(p)
	}
	return indices
}

func nodeAtPath(root *components.WireNode, path string) *components.WireNode {
	indices := parsePath(path)
	node := root
	for _, idx := range indices {
		if idx < 0 || idx >= len(node.Children) {
			return nil
		}
		node = &node.Children[idx]
	}
	return node
}

func parentAndChild(root *components.WireNode, path string) (*components.WireNode, int) {
	indices := parsePath(path)
	if len(indices) == 0 {
		return nil, -1
	}
	node := root
	for _, idx := range indices[:len(indices)-1] {
		if idx < 0 || idx >= len(node.Children) {
			return nil, -1
		}
		node = &node.Children[idx]
	}
	return node, indices[len(indices)-1]
}
