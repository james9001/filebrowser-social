package reactions

import (
	"time"
)

type Reaction struct {
	ID              uint      `storm:"id,increment" json:"id"`
	ReactionValue   string    `json:"reactionValue"`
	ContextFilePath string    `json:"contextFilePath"`
	UserName        string    `json:"userName"`
	CommentID       uint      `json:"commentId"`
	CreatedTime     time.Time `json:"createdTime"`
}
