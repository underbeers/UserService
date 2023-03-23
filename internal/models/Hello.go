package models

type Hello struct {
	// Structure for sending data to the ApiGateway
	Name      string     `json:"name"`
	Port      string     `json:"port"`
	IP        string     `json:"ip"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	URL       string   `json:"url"`
	Protected bool     `json:"protected"`
	Methods   []string `json:"methods"`
}
