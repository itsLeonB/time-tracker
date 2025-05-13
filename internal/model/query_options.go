package model

type QueryOptions struct {
	Filters          map[string]any
	PreloadRelations []string
}
