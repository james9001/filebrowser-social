package storage

import (
	"github.com/filebrowser/filebrowser/v2/auth"
	"github.com/filebrowser/filebrowser/v2/comments"
	"github.com/filebrowser/filebrowser/v2/notifications"
	"github.com/filebrowser/filebrowser/v2/reactions"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/share"
	"github.com/filebrowser/filebrowser/v2/usernotifications"
	"github.com/filebrowser/filebrowser/v2/users"
)

// Storage is a storage powered by a Backend which makes the necessary
// verifications when fetching and saving data to ensure consistency.
type Storage struct {
	Users             users.Store
	Comments          comments.Store
	Notifications     notifications.Store
	UserNotifications usernotifications.Store
	Reactions         reactions.Store
	Share             *share.Storage
	Auth              *auth.Storage
	Settings          *settings.Storage
}
