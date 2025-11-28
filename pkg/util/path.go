package util

func NameFromPath(p string) string {
	if len(p) == 0 {
		return ""
	}
	if p[len(p)-1] == '/' {
		p = p[:len(p)-1]
	}
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[i+1:]
		}
	}
	return p
}
