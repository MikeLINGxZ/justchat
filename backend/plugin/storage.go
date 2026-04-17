package plugin

import (
	"path/filepath"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"go.etcd.io/bbolt"
	bbolterrors "go.etcd.io/bbolt/errors"
)

type PluginStorage struct {
	db *bbolt.DB
}

func NewPluginStorage() (*PluginStorage, error) {
	dataPath, err := utils.GetDataPath()
	if err != nil {
		logger.Error("failed to get data path: %v", err)
		return nil, err
	}

	dbPath := filepath.Join(dataPath, "plugin_storage.db")
	db, err := bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		logger.Error("failed to open plugin storage db: %v", err)
		return nil, err
	}

	return &PluginStorage{db: db}, nil
}

func (s *PluginStorage) Close() error {
	return s.db.Close()
}

func (s *PluginStorage) Get(pluginId, key string) ([]byte, error) {
	var result []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(pluginId))
		if bucket == nil {
			return nil
		}
		v := bucket.Get([]byte(key))
		if v != nil {
			result = make([]byte, len(v))
			copy(result, v)
		}
		return nil
	})
	return result, err
}

func (s *PluginStorage) Set(pluginId, key string, value []byte) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(pluginId))
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), value)
	})
}

func (s *PluginStorage) Delete(pluginId, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(pluginId))
		if bucket == nil {
			return nil
		}
		return bucket.Delete([]byte(key))
	})
}

func (s *PluginStorage) DeleteAll(pluginId string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket([]byte(pluginId))
		if err == bbolterrors.ErrBucketNotFound {
			return nil
		}
		return err
	})
}
