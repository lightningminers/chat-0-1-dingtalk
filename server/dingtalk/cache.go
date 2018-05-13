package dingtalk

import (
	"encoding/json"
	"io/ioutil"
	"time"
	"errors"
)

type FileCache struct {
	Path string
}

type Expirable interface {
	CreatedAt() int64
	ExpiresIn() int
}

type Cache interface {
	Set(data Expirable) error
	Get(data Expirable) error
}


func NewFileCache(path string) *FileCache{
	return &FileCache{
		Path:path,
	}
}

func (f *FileCache) Get(data Expirable) error{
	bytes, err := ioutil.ReadFile(f.Path)
	if err == nil{
		err = json.Unmarshal(bytes, data)
		if err == nil{
			created := data.CreatedAt()
			expires := data.ExpiresIn()
			if err == nil && time.Now().Unix() > created + int64(expires - 60){
				err = errors.New("已过期")
			}
		}
	}
	return  err
}

func (f *FileCache) Set(data Expirable) error{
	bytes, err := json.Marshal(data)
	if err == nil{
		ioutil.WriteFile(f.Path, bytes, 0644)
	}
	return err
}