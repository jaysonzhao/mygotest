package vo

// Response
type HealthCheckResponse struct {
	Message string `json:"message"`
}

type Cat struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Deploy struct {
	Namespace string `json:"ns"`
	Name      string `json:"name"`
	Image     string `json:"image"`
}
