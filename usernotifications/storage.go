package usernotifications

type StorageBackend interface {
	FindByUserName(string) ([]*UserNotifications, error)
	Save(usernotifications *UserNotifications) error
	GetBy(interface{}) (*UserNotifications, error)
	Gets() ([]*UserNotifications, error)
	Update(usernotifications *UserNotifications, fields ...string) error
	DeleteByID(uint) error
}

type Store interface {
	FindByUserName(userName string) (*UserNotifications, error)
	GetById(usernotificationsId uint) (*UserNotifications, error)
	Save(usernotifications *UserNotifications) error
	Delete(id uint) error
	Gets() ([]*UserNotifications, error)
}

type Storage struct {
	back    StorageBackend
	updated map[uint]int64
}

func NewStorage(back StorageBackend) *Storage {
	return &Storage{
		back:    back,
		updated: map[uint]int64{},
	}
}

func (s *Storage) GetById(usernotificationsId uint) (*UserNotifications, error) {
	usernotifications, err := s.back.GetBy(usernotificationsId)
	if err != nil {
		return nil, err
	}

	return usernotifications, err
}

func (s *Storage) FindByUserName(userName string) (*UserNotifications, error) {
	usernotificationsSlice, err := s.back.FindByUserName(userName)
	if err != nil {
		return nil, err
	}
	if len(usernotificationsSlice) > 1 {
		return nil, nil
	}
	if len(usernotificationsSlice) == 0 {
		transparentlyCreatedUserNotifications := &UserNotifications{
			UserName:                  userName,
			AcknowledgedNotifications: make(map[uint]bool),
		}
		s.Save(transparentlyCreatedUserNotifications)
		return s.FindByUserName(userName)
	}

	usernotifications := usernotificationsSlice[0]

	return usernotifications, err
}

func (s *Storage) Save(usernotifications *UserNotifications) error {
	//TODO: validation here?
	// if err := usernotifications.Clean(""); err != nil {
	// 	return err
	// }

	return s.back.Save(usernotifications)
}

func (s *Storage) Delete(id uint) error {

	_, err := s.back.GetBy(id)
	if err != nil {
		return err
	}

	return s.back.DeleteByID(id)
}

func (s *Storage) Gets() ([]*UserNotifications, error) {
	usernotificationss, err := s.back.Gets()
	if err != nil {
		return nil, err
	}

	return usernotificationss, err
}
