package bolt

import (
	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"

	"github.com/filebrowser/filebrowser/v2/comments"
	"github.com/filebrowser/filebrowser/v2/errors"
)

type commentsBackend struct {
	db *storm.DB
}

func (st commentsBackend) DeleteByID(id uint) error {
	return st.db.DeleteStruct(&comments.Comment{ID: id})
}

func (st commentsBackend) GetBy(i interface{}) (comment *comments.Comment, err error) {
	comment = &comments.Comment{}

	var arg string
	switch i.(type) {
	case uint:
		arg = "ID"
	// case string: //etc.
	// 	arg = "Commentname"
	default:
		return nil, errors.ErrInvalidDataType
	}

	err = st.db.One(arg, i, comment)

	if err != nil {
		if err == storm.ErrNotFound {
			return nil, errors.ErrNotExist
		}
		return nil, err
	}

	return
}

func (st commentsBackend) FindByFilePath(filePath string) ([]*comments.Comment, error) {
	var allcomments []*comments.Comment
	err := st.db.Select(q.Eq("FilePath", filePath)).Find(&allcomments)
	//TODO: 0 comments is valid, but here is an error..
	if err == storm.ErrNotFound {
		return allcomments, nil
	}
	return allcomments, err
}

func (st commentsBackend) Save(comment *comments.Comment) error {
	err := st.db.Save(comment)
	if err == storm.ErrAlreadyExists {
		return errors.ErrExist
	}
	return err
}
