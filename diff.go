package main

import (
	"bytes"
	"strconv"

	"g.tizu.dev/CCWSUI/components"
	"g.tizu.dev/CCWSUI/web/webmsg"
)

func assignIDs(n *components.WireNode, prefix string) {
	n.ID = prefix
	for i := range n.Children {
		childPrefix := strconv.Itoa(i)
		if prefix != "" {
			childPrefix = prefix + "." + childPrefix
		}
		assignIDs(&n.Children[i], childPrefix)
	}
}

func diffTrees(old, new *components.WireNode) []webmsg.PatchEntry {
	var patches []webmsg.PatchEntry
	diffNode(old, new, "", &patches)
	return patches
}

func diffNode(old, new *components.WireNode, path string, patches *[]webmsg.PatchEntry) {
	if old == nil && new != nil {
		*patches = append(*patches, webmsg.PatchEntry{
			Action: webmsg.PatchInsert,
			Path:   path,
			Node:   new,
		})
		return
	}
	if old != nil && new == nil {
		*patches = append(*patches, webmsg.PatchEntry{
			Action: webmsg.PatchRemove,
			Path:   path,
		})
		return
	}
	if old.Kind != new.Kind {
		*patches = append(*patches, webmsg.PatchEntry{
			Action: webmsg.PatchReplace,
			Path:   path,
			Node:   new,
		})
		return
	}
	if !bytes.Equal(old.Props, new.Props) {
		*patches = append(*patches, webmsg.PatchEntry{
			Action: webmsg.PatchSetProps,
			Path:   path,
			Props:  new.Props,
		})
	}
	diffChildLists(path, old.Children, new.Children, patches)
}

func diffChildLists(parentPath string, old, new []components.WireNode, patches *[]webmsg.PatchEntry) {
	oi, ni := 0, 0
	for oi < len(old) || ni < len(new) {
		childPath := strconv.Itoa(ni)
		if parentPath != "" {
			childPath = parentPath + "." + childPath
		}

		if oi >= len(old) {
			for j := ni; j < len(new); j++ {
				insertPath := strconv.Itoa(j)
				if parentPath != "" {
					insertPath = parentPath + "." + insertPath
				}
				*patches = append(*patches, webmsg.PatchEntry{
					Action: webmsg.PatchInsert,
					Path:   insertPath,
					Node:   &new[j],
				})
			}
			break
		}
		if ni >= len(new) {
			for j := len(old) - 1; j >= oi; j-- {
				removePath := strconv.Itoa(j)
				if parentPath != "" {
					removePath = parentPath + "." + removePath
				}
				*patches = append(*patches, webmsg.PatchEntry{
					Action: webmsg.PatchRemove,
					Path:   removePath,
				})
			}
			break
		}

		if old[oi].Kind == new[ni].Kind {
			diffNode(&old[oi], &new[ni], childPath, patches)
			oi++
			ni++
		} else if ni+1 < len(new) && old[oi].Kind == new[ni+1].Kind {
			*patches = append(*patches, webmsg.PatchEntry{
				Action: webmsg.PatchInsert,
				Path:   childPath,
				Node:   &new[ni],
			})
			ni++
		} else if oi+1 < len(old) && old[oi+1].Kind == new[ni].Kind {
			*patches = append(*patches, webmsg.PatchEntry{
				Action: webmsg.PatchRemove,
				Path:   childPath,
			})
			oi++
		} else {
			*patches = append(*patches, webmsg.PatchEntry{
				Action: webmsg.PatchReplace,
				Path:   childPath,
				Node:   &new[ni],
			})
			oi++
			ni++
		}
	}
}


