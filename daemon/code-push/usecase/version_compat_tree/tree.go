package version_compat_tree

type versionRangeTree struct {
	compatVersionBranchMap
	*latestVersionBranch
	*strictCompatQueryBranch
}

func (v *versionRangeTree) Add(entries ...IEntry) {
	for _, entry := range entries {
		evictEntryList := v.compatVersionBranchMap.Enqueue(entry)
		if evictEntryList != nil {
			for _, entry := range evictEntryList {
				v.strictCompatQueryBranch.Evict(entry)
			}
		}

		v.latestVersionBranch.Enqueue(entry)
		v.strictCompatQueryBranch.Insert(entry)
	}
}

func (v *versionRangeTree) StrictCompat(anchor ICompatQueryAnchor) ICompatQueryResult {
	queryResult := &compatQueryResult{}

	canUpdateEntry := v.strictCompatQueryBranch.Query(anchor)
	strictLatest := v.latestVersionBranch.StrictLatest(anchor)

	queryResult.canUpdateVersion = canUpdateEntry
	queryResult.latestVersion = strictLatest

	return queryResult
}

func NewVersionCompatTree(entries []IEntry) ITree {
	tree := &versionRangeTree{
		compatVersionBranchMap: newCompatVersionBranchMap(),
		latestVersionBranch: newLatestVersionBranch(),
		strictCompatQueryBranch: &strictCompatQueryBranch{},
	}

	if entries != nil {
		tree.Add(entries...)
	}

	return tree
}

type compatQueryResult struct {
	latestVersion    IEntry
	canUpdateVersion IEntry
}

func (r *compatQueryResult) LatestVersion() IEntry {
	return r.latestVersion
}

func (r *compatQueryResult) CanUpdateVersion() IEntry {
	return r.canUpdateVersion
}
