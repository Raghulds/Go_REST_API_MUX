package entity

type Task struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Completed bool       `json:"completed"`
	Subtasks  []*SubTask `json:"subtasks"`
}

type SubTask struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}
