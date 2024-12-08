package model

type QueryOptions struct {
	Params   *QueryParams
	WithLogs bool
}

type QueryParams struct {
	Number string
}
