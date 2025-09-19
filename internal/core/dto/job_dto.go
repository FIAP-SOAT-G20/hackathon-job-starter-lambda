package dto

type StartJobInput struct {
	Name string `json:"job_name"`

	Variables map[string]string `json:"variables"`
}
