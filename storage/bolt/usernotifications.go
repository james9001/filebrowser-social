package bolt

import (
	"reflect"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/usernotifications"
)

type usernotificationsBackend struct {
	db *storm.DB
}

func (st usernotificationsBackend) DeleteByID(id uint) error {
	return st.db.DeleteStruct(&usernotifications.UserNotifications{ID: id})
}

func (st usernotificationsBackend) GetBy(i interface{}) (notification *usernotifications.UserNotifications, err error) {
	notification = &usernotifications.UserNotifications{}

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

func (st usernotificationsBackend) FindByUserName(userName string) ([]*usernotifications.UserNotifications, error) {
	var allusernotifications []*usernotifications.UserNotifications
	err := st.db.Select(q.Eq("UserName", userName)).Find(&allusernotifications)
	//TODO: 0 usernotifications is valid, but here is an error
	if err == storm.ErrNotFound {
		return allusernotifications, nil
	}
	return allusernotifications, err
}

func (st usernotificationsBackend) Gets() ([]*usernotifications.UserNotifications, error) {
	var allusernotifications []*usernotifications.UserNotifications
	err := st.db.All(&allusernotifications)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	if err != nil {
		return allusernotifications, err
	}

	return allusernotifications, err
}

func (st usernotificationsBackend) Save(notification *usernotifications.UserNotifications) error {
	err := st.db.Save(notification)
	if err == storm.ErrAlreadyExists {
		return errors.ErrExist
	}
	return err
}

func (st usernotificationsBackend) Update(notification *usernotifications.UserNotifications, fields ...string) error {
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
