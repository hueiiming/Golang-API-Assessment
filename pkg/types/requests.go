package types

type RegisterRequest struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

type SuspendRequest struct {
	Student string `json:"student"`
}

type NotificationRequest struct {
	Teacher      string `json:"teacher"`
	Notification string `json:"notification"`
}

type BulkImportRequest struct {
	Teachers []string `json:"teachers"`
	Students []string `json:"students"`
}
