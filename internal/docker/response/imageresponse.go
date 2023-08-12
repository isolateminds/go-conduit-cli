package response

type Image struct {
	Stream   string `json:"stream"`
	Status   string `json:"status"`
	Progress string `json:"progress"`
	Aux      struct {
		ID string `json:"id"`
	} `json:"aux"`
	ErrorDetail struct {
		Message string `json:"message"`
	} `json:"errorDetail"`
	Error string `json:"error"`
}
