package task

type SummaryOutput struct {
	Total     int          `json:"total"`
	Completed int          `json:"completed"`
	Pending   int          `json:"pending"`
	Tasks     []TaskOutput `json:"tasks"`
}
