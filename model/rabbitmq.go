package model

type IndexCreate struct {
	Index    string
	BodyJson any `json:"mappings"`
}

type AliasCreate struct {
	Index string `json:"index"`
	Name  string `json:"alias"`
}

type EsAction struct {
	Actions []any `json:"actions"`
}

type RabbitMQPayload struct {
	HttpMethod string `json:"http_method"`
	Uri        string `json:"uri"`
	Body       any    `json:"body"`
}
