package comments

import (
	"time"
)

type Comment struct {
	ID          uint      `storm:"id,increment" json:"id"`
	FilePath    string    `json:"filePath"`
	CommentText string    `json:"commentText"`
	UserName    string    `json:"userName"`
	CreatedTime time.Time `json:"createdTime"`
}
