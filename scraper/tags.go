package scraper

func MakeTagsResolver(provided []string) TagsResolver {
	return TagsResolver{
		provided: makeStringSet(provided),
	}
}

type TagsResolver struct {
	provided map[string]bool
}

func (t *TagsResolver) isMatch(expected []string) bool {
	expSet := makeStringSet(expected)

	for key := range t.provided {
		if !expSet[key] {
			return false
		}
	}

	return true
}

func makeStringSet(provided []string) map[string]bool {
	set := make(map[string]bool)

	for _, item := range provided {
		set[item] = true
	}

	return set
}
