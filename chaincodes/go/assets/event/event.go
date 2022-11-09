package event

type Event struct {
	Type    string `json:"Type,omitempty"`
	Message string `json:"Message,omitempty"`
}
