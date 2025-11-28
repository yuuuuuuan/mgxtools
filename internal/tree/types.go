package tree

type Node struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Type        string  `json:"type"`
	Size        int64   `json:"size"`
	UpdatedTime float64 `json:"updated_time"`
	Children    []Node  `json:"children"`
}
