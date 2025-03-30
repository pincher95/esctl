package index

type IndexSettingsResponse map[string]any

type Index interface {
	UpdateIndexSettings(endpoint, index *string, body *map[string]any, flat bool) (*IndexSettingsResponse, error)
	GetAliases(endpoint, index *string) (*IndexAliasResponse, error)
}

type index struct {
	Index
}

func NewIndex() Index {
	return &index{}
}
