package types

type RegisterRequest struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

type SuspendRequest struct {
	Student string `json:"student"`
}

type NotificationRequest struct {
	Teacher string `json:"teacher"`
	Message string `json:"notification"`
}

type CommonStudentsResponse struct {
	Students []string `json:"students"`
}

type NotificationResponse struct {
	Recipients []string `json:"recipients"`
}
