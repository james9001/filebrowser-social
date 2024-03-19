package notifications

type StorageBackend interface {
	FindByContextFilePath(string) ([]*Notification, error)
	Save(notification *Notification) error
	GetBy(interface{}) (*Notification, error)
	Gets() ([]*Notification, error)
	Update(notification *Notification, fields ...string) error
	DeleteByID(uint) error
	FetchPage(pagenum int) ([]*Notification, error)
	GetRange(lowId int, highId int) ([]*Notification, error)
}

type Store interface {
	GetById(notificationId uint) (*Notification, error)
	FindByContextFilePath(contextFilePath string) ([]*Notification, error)
	Save(notification *Notification) error
	Delete(id uint) error
	Gets() ([]*Notification, error)
	FetchPage(pagenum int) ([]*Notification, error)
	GetRange(lowId int, highId int) ([]*Notification, error)
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

func (s *Storage) GetById(notificationId uint) (*Notification, error) {
	notification, err := s.back.GetBy(notificationId)
	if err != nil {
		return nil, err
	}

	return notification, err
}

func (s *Storage) FindByContextFilePath(contextFilePath string) ([]*Notification, error) {
	notifications, err := s.back.FindByContextFilePath(contextFilePath)
	if err != nil {
		return nil, err
	}

	return notifications, err
}

func (s *Storage) Save(notification *Notification) error {
	//TODO: validation here?
	// if err := notification.Clean(""); err != nil {
	// 	return err
	// }

	return s.back.Save(notification)
}

func (s *Storage) Delete(id uint) error {

	_, err := s.back.GetBy(id)
	if err != nil {
		return err
	}

	return s.back.DeleteByID(id)
}

func (s *Storage) Gets() ([]*Notification, error) {
	notifications, err := s.back.Gets()
	if err != nil {
		return nil, err
	}

	return notifications, err
}

func (s *Storage) FetchPage(pagenum int) ([]*Notification, error) {
	notifications, err := s.back.FetchPage(pagenum)
	if err != nil {
		return nil, err
	}

	return notifications, err
}

func (s *Storage) GetRange(lowId int, highId int) ([]*Notification, error) {
	notifications, err := s.back.GetRange(lowId, highId)
	if err != nil {
		return nil, err
	}

	return notifications, err
}
