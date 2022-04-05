package cache

import (
	"bitbucket.org/HeilaSystems/cacheStorage"
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/cache"
	"context"
	"fmt"
	"time"
)

type cacheRepo struct {
	cacheGetter           cache.CacheStorageGetterWrapper
	cacheSetter           cache.CacheStorageSetterWrapper
	sessionCollectionName string
	version               string
}

func NewSessionCacheRepo(cacheGetter cache.CacheStorageGetterWrapper, cacheSetter cache.CacheStorageSetterWrapper, collectionName string, version string) *cacheRepo {
	return &cacheRepo{cacheGetter: cacheGetter, cacheSetter: cacheSetter, sessionCollectionName: collectionName, version: version}
}

func (r cacheRepo) GetUserSessionByTokenToStruct(c context.Context, token string, dest interface{}) (bool, error) {
	err := r.cacheGetter.GetById(c, r.sessionCollectionName, token, r.version, dest)
	if err != nil {
		if err.IsNotFound() {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (r cacheRepo) InsertOrUpdate(ctx context.Context, id string, obj interface{}) error {
	return r.cacheSetter.InsertOrUpdate(ctx, r.sessionCollectionName, id, r.version, obj)
}

func (r cacheRepo) GetCacheVersions(ctx context.Context, now time.Time) (map[string]string, error) {
	versions, err := r.cacheGetter.GetLatestVersions(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)

	for i := range versions {
		var latestVersion cacheStorage.Version
		for _, v := range versions[i].Versions {
			if (latestVersion.TimedTo.IsZero() || v.TimedTo.After(latestVersion.TimedTo)) && v.TimedTo.Before(now) {
				latestVersion = v
			}
		}
		if latestVersion.Version == "" {
			return result, fmt.Errorf("no version found for collection %v by date %v", versions[i].CollectionName, now)
		}
		result[versions[i].CollectionName] = latestVersion.Version
	}
	return result, nil
}
