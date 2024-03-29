package cache

import (
	"context"
	"fmt"
	"github.com/orchestd/cacheStorage"
	"github.com/orchestd/dependencybundler/interfaces/cache"
	"github.com/orchestd/sharedlib/slices"
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

func (r cacheRepo) GetCollectionsFilterActions(ctx context.Context, filterAction string) ([]string, error) {
	cacheCollections, err := r.cacheGetter.GetLatestVersions(ctx)
	if err != nil {
		return nil, err
	}
	result := []string{}
	for _, cacheCollection := range cacheCollections {
		if filterAction != "" && !slices.IsStrExist(cacheCollection.LockVersionUpon, filterAction) {
			continue
		}
		result = append(result, cacheCollection.CollectionName)
	}
	return result, err
}

func (r cacheRepo) GetCacheVersions(ctx context.Context, now time.Time, filterAction string, filterType string) (map[string]string, error) {
	cacheCollections, err := r.cacheGetter.GetLatestVersions(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string)

	for _, cacheCollection := range cacheCollections {
		if filterAction != "" && !slices.IsStrExist(cacheCollection.LockVersionUpon, filterAction) {
			continue
		}
		if filterType != "" && cacheCollection.CacheType != filterType {
			continue
		}
		var latestVersion cacheStorage.Version
		for _, v := range cacheCollection.Versions {
			if (latestVersion.TimedTo.IsZero() || v.TimedTo.After(latestVersion.TimedTo)) && v.TimedTo.Before(now) {
				latestVersion = v
			}
		}
		if latestVersion.Version == "" {
			return result, fmt.Errorf("no version found for collection %v by date %v", cacheCollection.CollectionName, now)
		}
		result[cacheCollection.CollectionName] = latestVersion.Version
	}
	return result, nil
}
