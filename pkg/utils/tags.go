package utils

// MakeTagsResolver make an instance of the tags resolver for provided slice of
// strings
func MakeTagsResolver(provided []string) TagsResolver {
	return TagsResolver{
		provided: makeStringSet(provided),
	}
}

// TagsResolver for resolving tags - for provided set of tags
type TagsResolver struct {
	provided map[string]bool
}

// IsMatch to check whether the expected tags list is subset of provided
func (t *TagsResolver) IsMatch(expected []string) bool {
	expSet := makeStringSet(expected)

	if len(t.provided) == 0 {
		return false
	}

	for key := range t.provided {
		if !expSet[key] {
			return false
		}
	}

	return true
}

// makeStringSet - make a string set from provided string slice
func makeStringSet(provided []string) map[string]bool {
	set := make(map[string]bool)

	for _, item := range provided {
		set[item] = true
	}

	return set
}
