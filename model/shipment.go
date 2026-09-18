package model

type CreateShipmentRequest struct {
	TrackingNumber string  `json:"tracking_number"`
	SenderName     string  `json:"sender_name"`
	ReceiverName   string  `json:"receiver_name"`
	Origin         string  `json:"origin"`
	Destination    string  `json:"destination"`
	Weight         float64 `json:"weight"`
}

type UpdateShipmentRequest struct {
	SenderName   string  `json:"sender_name"`
	ReceiverName string  `json:"receiver_name"`
	Origin       string  `json:"origin"`
	Destination  string  `json:"destination"`
	Status       string  `json:"status"`
	Weight       float64 `json:"weight"`
}
