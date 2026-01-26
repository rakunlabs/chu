package loader

import (
	"slices"
)

func OrderLoaders(loaders map[LoaderName]Loader) []LoaderName {
	// Build dependency graph
	afterList := make(map[LoaderName][]LoaderName)
	lastList := make([]LoaderName, 0)
	inDegree := make(map[LoaderName]int)

	// Create afterList and inDegree maps
	for name, l := range loaders {
		empty := true

		// Convert before to after
		if len(l.LoadOrder().Before) > 0 {
			for _, before := range l.LoadOrder().Before {
				if _, ok := loaders[before]; !ok {
					continue
				}

				if !slices.Contains(afterList[name], before) {
					afterList[name] = append(afterList[name], before)
					inDegree[before]++
				}

				empty = false
			}
		}
		if len(l.LoadOrder().After) > 0 {
			for _, after := range l.LoadOrder().After {
				if _, ok := loaders[after]; !ok {
					continue
				}

				if !slices.Contains(afterList[after], name) {
					afterList[after] = append(afterList[after], name)
					inDegree[name]++
				}

				empty = false
			}
		}

		if empty {
			lastList = append(lastList, name)

			continue
		}

		if inDegree[name] == 0 {
			// No dependencies, add to inDegree map
			inDegree[name] = 0
		}
	}

	// Kahn's algorithm for topological sort with name
	var ordered []LoaderName

	queue := make([]LoaderName, 0, len(inDegree))
	inDegreeKey := make([]LoaderName, 0, len(inDegree))
	for name := range inDegree {
		inDegreeKey = append(inDegreeKey, name)
	}
	slices.Sort(inDegreeKey)
	for _, name := range inDegreeKey {
		if v, ok := inDegree[name]; ok && v == 0 {
			queue = append(queue, name)
		}
	}

	// Order afterList
	for name := range afterList {
		slices.Sort(afterList[name])
	}
	// Order lastList
	slices.Sort(lastList)

	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]
		ordered = append(ordered, name)

		for _, next := range afterList[name] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}

	ordered = append(ordered, lastList...)

	return ordered
}
