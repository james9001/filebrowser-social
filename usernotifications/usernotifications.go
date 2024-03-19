package usernotifications

type UserNotifications struct {
	ID                        uint          `storm:"id,increment" json:"id"`
	UserName                  string        `json:"userName"`
	AcknowledgedNotifications map[uint]bool `json:"acknowledgedNotifications"`
}
