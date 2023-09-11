package data

// logSorter is a custom type to make sorting Logs easier.
type logSorter []*Log

// Len implements sort.Interface.
func (s logSorter) Len() int {
	return len(s)
}

// Less implements sort.Interface.
func (s logSorter) Less(i, j int) bool {
	if s[i].Author < s[j].Author {
		return true
	} else if s[i].Author > s[j].Author {
		return false
	}

	// If Author Emails are equal, compare by Date
	if s[i].Date.AsTime().Before(s[j].Date.AsTime()) {
		return true
	} else if s[i].Date.AsTime().After(s[j].Date.AsTime()) {
		return false
	}

	// If Dates are equal, compare by Path
	if s[i].Path < s[j].Path {
		return true
	} else if s[i].Path > s[j].Path {
		return false
	}

	// If Dates and Paths are equal, compare by Diff
	return s[i].Diff < s[j].Diff
}

// Swap implements sort.Interface.
func (s logSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
