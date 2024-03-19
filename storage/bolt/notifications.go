package bolt

import (
	"reflect"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/notifications"
)

type notificationsBackend struct {
	db *storm.DB
}

func (st notificationsBackend) DeleteByID(id uint) error {
	return st.db.DeleteStruct(&notifications.Notification{ID: id})
}

func (st notificationsBackend) GetBy(i interface{}) (notification *notifications.Notification, err error) {
	notification = &notifications.Notification{}

	var arg string
	switch i.(type) {
	case uint:
		arg = "ID"
	default:
		return nil, errors.ErrInvalidDataType
	}

	err = st.db.One(arg, i, notification)

	if err != nil {
		if err == storm.ErrNotFound {
			return nil, errors.ErrNotExist
		}
		return nil, err
	}

	return
}

func (st notificationsBackend) FindByContextFilePath(contextFilePath string) ([]*notifications.Notification, error) {
	var allnotifications []*notifications.Notification
	err := st.db.Select(q.Eq("ContextFilePath", contextFilePath)).Find(&allnotifications)
	//TODO: 0 notifications is valid, but here is an error
	if err == storm.ErrNotFound {
		return allnotifications, nil
	}
	return allnotifications, err
}

func (st notificationsBackend) Gets() ([]*notifications.Notification, error) {
	var allnotifications []*notifications.Notification
	err := st.db.All(&allnotifications)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	if err != nil {
		return allnotifications, err
	}

	return allnotifications, err
}

func (st notificationsBackend) Save(notification *notifications.Notification) error {
	err := st.db.Save(notification)
	if err == storm.ErrAlreadyExists {
		return errors.ErrExist
	}
	return err
}

func (st notificationsBackend) Update(notification *notifications.Notification, fields ...string) error {
	if len(fields) == 0 {
		return st.Save(notification)
	}

	for _, field := range fields {
		val := reflect.ValueOf(notification).Elem().FieldByName(field).Interface()
		if err := st.db.UpdateField(notification, field, val); err != nil {
			return err
		}
	}

	return nil
}

func (st notificationsBackend) FetchPage(pagenum int) ([]*notifications.Notification, error) {
	var notifications []*notifications.Notification
	err := st.db.Select().Limit(10).Skip(10 * pagenum).OrderBy("ID").Reverse().Find(&notifications)
	//TODO: 0 notifications is valid, but here is an error
	if err == storm.ErrNotFound {
		return notifications, nil
	}
	return notifications, err
}

func (st notificationsBackend) GetRange(lowId int, highId int) ([]*notifications.Notification, error) {
	var notifications []*notifications.Notification
	err := st.db.Range("ID", lowId, highId, &notifications, storm.Reverse())
	//TODO: 0 notifications is valid, but here is an error
	if err == storm.ErrNotFound {
		return notifications, nil
	}
	return notifications, err
}
