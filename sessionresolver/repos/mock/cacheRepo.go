package mock

import (
	"context"
	"encoding/json"
	"fmt"
)

type cacheRepoMock struct {
	versions map[string]string
	sessions map[string]interface{}
}

func (c cacheRepoMock) GetUserSessionByTokenToStruct(context context.Context, token string, dest interface{}) error {
	val, ok := c.sessions[token]
	if !ok {
		return fmt.Errorf("notFound")
	}
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return err
	}
	return nil
}

func (c cacheRepoMock) InsertOrUpdate(ctx context.Context, id string, obj interface{}) error {
	if obj == nil {
		return fmt.Errorf("cannotInsertNilObject")
	}
	c.sessions[id] = obj
	return nil
}

func (c cacheRepoMock) GetCacheVersions(ctx context.Context) (map[string]string, error) {
	return c.versions , nil
}

func NewCacheRepoMock(versionsFakeData map[string]string , sessionsFakeData map[string]interface{}) *cacheRepoMock {
	return &cacheRepoMock{
		versions: versionsFakeData,
		sessions: sessionsFakeData,
	}
}
