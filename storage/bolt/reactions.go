package bolt

import (
	"reflect"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"

	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/reactions"
)

type reactionsBackend struct {
	db *storm.DB
}

func (st reactionsBackend) DeleteByID(id uint) error {
	return st.db.DeleteStruct(&reactions.Reaction{ID: id})
}

func (st reactionsBackend) GetBy(i interface{}) (reaction *reactions.Reaction, err error) {
	reaction = &reactions.Reaction{}

	var arg string
	switch i.(type) {
	case uint:
		arg = "ID"
	default:
		return nil, errors.ErrInvalidDataType
	}

	err = st.db.One(arg, i, reaction)

	if err != nil {
		if err == storm.ErrNotFound {
			return nil, errors.ErrNotExist
		}
		return nil, err
	}

	return
}

func (st reactionsBackend) FindByContextFilePath(contextFilePath string) ([]*reactions.Reaction, error) {
	var allreactions []*reactions.Reaction
	err := st.db.Select(q.Eq("ContextFilePath", contextFilePath)).Find(&allreactions)
	//TODO: 0 reactions is valid, but here is an error
	if err == storm.ErrNotFound {
		return allreactions, nil
	}
	return allreactions, err
}

func (st reactionsBackend) FindByContextFilePathAndUserName(contextFilePath string, userName string) ([]*reactions.Reaction, error) {
	var allreactions []*reactions.Reaction
	err := st.db.Select(q.And(
		q.Eq("ContextFilePath", contextFilePath),
		q.Eq("UserName", userName),
	)).Find(&allreactions)
	//TODO: 0 reactions is valid, but here is an error
	if err == storm.ErrNotFound {
		return allreactions, nil
	}
	return allreactions, err
}

func (st reactionsBackend) Gets() ([]*reactions.Reaction, error) {
	var allreactions []*reactions.Reaction
	err := st.db.All(&allreactions)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	if err != nil {
		return allreactions, err
	}

	return allreactions, err
}

func (st reactionsBackend) Save(reaction *reactions.Reaction) error {
	err := st.db.Save(reaction)
	if err == storm.ErrAlreadyExists {
		return errors.ErrExist
	}
	return err
}

func (st reactionsBackend) Update(reaction *reactions.Reaction, fields ...string) error {
	if len(fields) == 0 {
		return st.Save(reaction)
	}

	for _, field := range fields {
		val := reflect.ValueOf(reaction).Elem().FieldByName(field).Interface()
		if err := st.db.UpdateField(reaction, field, val); err != nil {
			return err
		}
	}

	return nil
}

func (st reactionsBackend) FetchPage(pagenum int) ([]*reactions.Reaction, error) {
	var reactions []*reactions.Reaction
	err := st.db.Select().Limit(10).Skip(10 * pagenum).OrderBy("ID").Reverse().Find(&reactions)
	//TODO: 0 reactions is valid, but here is an error
	if err == storm.ErrNotFound {
		return reactions, nil
	}
	return reactions, err
}

func (st reactionsBackend) GetRange(lowId int, highId int) ([]*reactions.Reaction, error) {
	var reactions []*reactions.Reaction
	err := st.db.Range("ID", lowId, highId, &reactions, storm.Reverse())
	//TODO: 0 reactions is valid, but here is an error
	if err == storm.ErrNotFound {
		return reactions, nil
	}
	return reactions, err
}
