package comments

type StorageBackend interface {
	FindByFilePath(string) ([]*Comment, error)
	Save(comment *Comment) error
	DeleteByID(uint) error
	GetBy(interface{}) (*Comment, error)
}

type Store interface {
	GetById(commentId uint) (*Comment, error)
	FindByFilePath(filePath string) ([]*Comment, error)
	Save(comment *Comment) error
	Delete(id uint) error
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

func (s *Storage) GetById(commentId uint) (*Comment, error) {
	comment, err := s.back.GetBy(commentId)
	if err != nil {
		return nil, err
	}

	return comment, err
}

func (s *Storage) FindByFilePath(filePath string) ([]*Comment, error) {
	comments, err := s.back.FindByFilePath(filePath)
	if err != nil {
		return nil, err
	}

	return comments, err
}

func (s *Storage) Save(comment *Comment) error {
	//TODO: validation here?
	// if err := comment.Clean(""); err != nil {
	// 	return err
	// }

	return s.back.Save(comment)
}

func (s *Storage) Delete(id uint) error {

	_, err := s.back.GetBy(id)
	if err != nil {
		return err
	}

	return s.back.DeleteByID(id)
}
