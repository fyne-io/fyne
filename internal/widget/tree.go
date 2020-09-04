package widget

// AddTreePath adds the given path to the given parent->children map
func AddTreePath(data map[string][]string, path ...string) {
	parent := ""
	for _, p := range path {
		children := data[parent]
		add := true
		for _, c := range children {
			if c == p {
				add = false
				break
			}
		}
		if add {
			data[parent] = append(children, p)
		}
		parent = p
	}
}
