package controller

type EventTarget struct {
	MediaType  string `json:"mediaType"`
	Size       int    `json:"size"`
	Digest     string `json:"digest"`
	Length     int    `json:"length"`
	Repository string `json:"repository"`
	URL        string `json:"url"`
	Tag        string `json:"tag"`
}

type EventRequest struct {
	ID        string `json:"id"`
	Addr      string `json:"addr"`
	Host      string `json:"host"`
	Method    string `json:"method"`
	UserAgent string `json:"useragent"`
}

type EventSource struct {
	Addr       string `json:"addr"`
	InstanceID string `json:"instanceID"`
}

type Event struct {
	ID        string                 `json:"id"`
	TimeStamp string                 `json:"timestamp"`
	Action    string                 `json:"action"`
	Target    EventTarget            `json:"target"`
	Request   EventRequest           `json:"request"`
	Actor     map[string]interface{} `json:"actor"`
	Source    EventSource            `json:"source"`
}

type Events struct {
	Events []Event `json:"events"`
}
