package types

type CommonStudentsResponse struct {
	Students []string `json:"students"`
}

type NotificationResponse struct {
	Recipients []string `json:"recipients"`
}
