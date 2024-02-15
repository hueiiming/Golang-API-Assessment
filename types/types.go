package types

type RegisterRequest struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

type CommonStudents struct {
	Students []string `json:"students"`
}

type Notification struct {
	Recipients []string `json:"recipients"`
}
