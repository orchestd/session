package cache

import (
	"bitbucket.org/HeilaSystems/dependencybundler/interfaces/cache"
	"context"
)

type CacheFunctions interface {
	GetById(c context.Context, collectionName string, id interface{}, dest interface{}) cacheRepo
}
type cacheRepo struct {
	cache                 cache.CacheStorageGetter
	sessionCollectionName string
	version string
}

func NewSessionCacheRepo(cache cache.CacheStorageGetter,collectionName string , version string) *cacheRepo {
	return &cacheRepo{cache: cache,sessionCollectionName: collectionName ,version: version}
}

func (r cacheRepo) GetUserSessionByTokenToStruct(c context.Context,token string,dest interface{}) error {
	if err := r.cache.GetById(c , r.sessionCollectionName , token,r.version,dest);err != nil {
		return err
	}
	return nil
}