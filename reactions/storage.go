package reactions

import (
	"github.com/filebrowser/filebrowser/v2/errors"
)

type StorageBackend interface {
	FindByContextFilePath(string) ([]*Reaction, error)
	Save(Reaction *Reaction) error
	GetBy(interface{}) (*Reaction, error)
	Gets() ([]*Reaction, error)
	Update(Reaction *Reaction, fields ...string) error
	DeleteByID(uint) error
	FetchPage(pagenum int) ([]*Reaction, error)
	GetRange(lowId int, highId int) ([]*Reaction, error)
	FindByContextFilePathAndUserName(contextFilePath string, userName string) ([]*Reaction, error)
}

type Store interface {
	GetById(ReactionId uint) (*Reaction, error)
	FindByContextFilePath(contextFilePath string) ([]*Reaction, error)
	Save(Reaction *Reaction) error
	Delete(id uint) error
	Gets() ([]*Reaction, error)
	FetchPage(pagenum int) ([]*Reaction, error)
	GetRange(lowId int, highId int) ([]*Reaction, error)
	FindByContextFilePathAndUserName(contextFilePath string, userName string) ([]*Reaction, error)
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

func (s *Storage) GetById(ReactionId uint) (*Reaction, error) {
	Reaction, err := s.back.GetBy(ReactionId)
	if err != nil {
		return nil, err
	}

	return Reaction, err
}

func (s *Storage) FindByContextFilePath(contextFilePath string) ([]*Reaction, error) {
	Reactions, err := s.back.FindByContextFilePath(contextFilePath)
	if err != nil {
		return nil, err
	}

	return Reactions, err
}

func (s *Storage) FindByContextFilePathAndUserName(contextFilePath string, userName string) ([]*Reaction, error) {
	Reactions, err := s.back.FindByContextFilePathAndUserName(contextFilePath, userName)
	if err != nil {
		return nil, err
	}

	return Reactions, err
}

func (s *Storage) Save(Reaction *Reaction) error {
	Reactions, someErr := s.back.FindByContextFilePathAndUserName(Reaction.ContextFilePath, Reaction.UserName)
	if someErr != nil {
		return someErr
	}

	for _, a := range Reactions {
		if Reaction.CommentID == a.CommentID && Reaction.ReactionValue == a.ReactionValue {
			return errors.ErrExist
		}
	}

	return s.back.Save(Reaction)
}

func (s *Storage) Delete(id uint) error {

	_, err := s.back.GetBy(id)
	if err != nil {
		return err
	}

	return s.back.DeleteByID(id)
}

func (s *Storage) Gets() ([]*Reaction, error) {
	Reactions, err := s.back.Gets()
	if err != nil {
		return nil, err
	}

	return Reactions, err
}

func (s *Storage) FetchPage(pagenum int) ([]*Reaction, error) {
	Reactions, err := s.back.FetchPage(pagenum)
	if err != nil {
		return nil, err
	}

	return Reactions, err
}

func (s *Storage) GetRange(lowId int, highId int) ([]*Reaction, error) {
	Reactions, err := s.back.GetRange(lowId, highId)
	if err != nil {
		return nil, err
	}

	return Reactions, err
}
