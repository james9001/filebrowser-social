package bolt

import (
	"github.com/asdine/storm/v3"

	"github.com/filebrowser/filebrowser/v2/auth"
	"github.com/filebrowser/filebrowser/v2/comments"
	"github.com/filebrowser/filebrowser/v2/notifications"
	"github.com/filebrowser/filebrowser/v2/reactions"
	"github.com/filebrowser/filebrowser/v2/settings"
	"github.com/filebrowser/filebrowser/v2/share"
	"github.com/filebrowser/filebrowser/v2/storage"
	"github.com/filebrowser/filebrowser/v2/usernotifications"
	"github.com/filebrowser/filebrowser/v2/users"
)

// NewStorage creates a storage.Storage based on Bolt DB.
func NewStorage(db *storm.DB) (*storage.Storage, error) {
	userStore := users.NewStorage(usersBackend{db: db})
	commentStore := comments.NewStorage(commentsBackend{db: db})
	notificationsStore := notifications.NewStorage(notificationsBackend{db: db})
	userNotificationsStore := usernotifications.NewStorage(usernotificationsBackend{db: db})
	reactionsStore := reactions.NewStorage(reactionsBackend{db: db})
	shareStore := share.NewStorage(shareBackend{db: db})
	settingsStore := settings.NewStorage(settingsBackend{db: db})
	authStore := auth.NewStorage(authBackend{db: db}, userStore)

	err := save(db, "version", 2) //nolint:gomnd
	if err != nil {
		return nil, err
	}

	return &storage.Storage{
		Auth:              authStore,
		Comments:          commentStore,
		Notifications:     notificationsStore,
		UserNotifications: userNotificationsStore,
		Reactions:         reactionsStore,
		Users:             userStore,
		Share:             shareStore,
		Settings:          settingsStore,
	}, nil
}
