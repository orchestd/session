package cache

import (
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/cache"
	"context"
)

type cacheRepo struct {
	cacheGetter           cache.CacheStorageGetter
	cacheSetter           cache.CacheStorageSetter
	sessionCollectionName string
	version               string
}

func NewSessionCacheRepo(cacheGetter cache.CacheStorageGetter, cacheSetter cache.CacheStorageSetter, collectionName string, version string) *cacheRepo {
	return &cacheRepo{cacheGetter: cacheGetter, cacheSetter: cacheSetter, sessionCollectionName: collectionName, version: version}
}

func (r cacheRepo) GetUserSessionByTokenToStruct(c context.Context, token string, dest interface{}) error {
	if err := r.cacheGetter.GetById(c, r.sessionCollectionName, token, r.version, dest); err != nil {
		return err
	}
	return nil
}

func (r cacheRepo) InsertOrUpdate(ctx context.Context, id string, obj interface{}) error {
	return r.cacheSetter.InsertOrUpdate(ctx, r.sessionCollectionName, id, r.version, obj)
}

func (r cacheRepo) GetCacheVersions(ctx context.Context) (map[string]string, error) {
	versions, err := r.cacheGetter.GetLatestVersions(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)
	for i := range versions {
		result[versions[i].CollectionName] = versions[i].Version
	}
	return result, nil
}
