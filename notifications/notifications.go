package notifications

import (
	"time"
)

type Notification struct {
	ID               uint      `storm:"id,increment" json:"id"`
	ContextFilePath  string    `json:"contextFilePath"`
	CausingUserName  string    `json:"causingUserName"`
	NotificationType string    `json:"notificationType"`
	CreatedTime      time.Time `json:"createdTime"`
}
