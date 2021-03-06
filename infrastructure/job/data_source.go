package job

import (
	"bytes"
	"encoding/json"
	"github.com/duck8823/duci/domain/model/job"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
)

type dataSource struct {
	db LevelDB
}

// NewDataSource returns job data source
func NewDataSource(path string) (job.Repository, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &dataSource{db}, nil
}

// FindBy returns job found by ID
func (d *dataSource) FindBy(id job.ID) (*job.Job, error) {
	data, err := d.db.Get(id.ToSlice(), nil)
	if err == leveldb.ErrNotFound {
		return nil, job.ErrNotFound
	} else if err != nil {
		return nil, errors.WithStack(err)
	}

	job := &job.Job{}
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(job); err != nil {
		return nil, errors.WithStack(err)
	}
	job.ID = id
	return job, nil
}

// Save store job to data source
func (d *dataSource) Save(job job.Job) error {
	data, err := job.ToBytes()
	if err != nil {
		return errors.WithStack(err)
	}
	if err := d.db.Put(job.ID.ToSlice(), data, nil); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
