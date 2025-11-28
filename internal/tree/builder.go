package tree

import (
	"mgxtools/internal/api"
	"mgxtools/pkg/util"
)

func BuildTree(client *api.Client, path string) (Node, error) {
	items, err := client.ListDirectory(path)
	if err != nil {
		return Node{}, err
	}

	node := Node{
		Name: util.NameFromPath(path),
		Path: path,
		Type: "dir",
	}

	children := []Node{}

	for _, item := range items {
		full := path + item.Name

		if item.IsDir {
			sub := full + "/"
			child, err := BuildTree(client, sub)
			if err != nil {
				return Node{}, err
			}
			child.Name = item.Name
			child.Path = sub
			child.Type = "dir"
			child.Size = item.Size
			child.UpdatedTime = item.UpdatedTime
			children = append(children, child)
		} else {
			children = append(children, Node{
				Name:        item.Name,
				Path:        full,
				Type:        "file",
				Size:        item.Size,
				UpdatedTime: item.UpdatedTime,
			})
		}
	}

	node.Children = children
	return node, nil
}
