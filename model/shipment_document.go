package model

type ShipmentDocument struct {
	ID             int     `json:"id"`
	TrackingNumber string  `json:"tracking_number"`
	SenderName     string  `json:"sender_name"`
	ReceiverName   string  `json:"receiver_name"`
	Origin         string  `json:"origin"`
	Destination    string  `json:"destination"`
	Status         string  `json:"status"`
	Weight         float64 `json:"weight"`
}
