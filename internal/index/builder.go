package index

// BuildFromSnapshot populates an Index from a map of environment name to
// list of secret paths, as produced by the snapshot package.
func BuildFromSnapshot(data map[string][]string) *Index {
	idx := New()
	for env, paths := range data {
		for _, p := range paths {
			idx.Add(env, p)
		}
	}
	return idx
}

// Merge combines two indexes into a new Index. Entries from both are included;
// duplicate (env, path) pairs are deduplicated.
func Merge(a, b *Index) *Index {
	out := New()
	for _, env := range a.Environments() {
		for _, p := range a.PathsForEnv(env) {
			out.Add(env, p)
		}
	}
	for _, env := range b.Environments() {
		for _, p := range b.PathsForEnv(env) {
			out.Add(env, p)
		}
	}
	return out
}

// Subset returns a new Index containing only entries for the given environments.
func Subset(idx *Index, envs []string) *Index {
	envSet := make(map[string]struct{}, len(envs))
	for _, e := range envs {
		envSet[e] = struct{}{}
	}
	out := New()
	for _, env := range idx.Environments() {
		if _, ok := envSet[env]; !ok {
			continue
		}
		for _, p := range idx.PathsForEnv(env) {
			out.Add(env, p)
		}
	}
	return out
}
