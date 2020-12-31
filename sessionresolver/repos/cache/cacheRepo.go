package cache

import "context"

type CacheFunctions interface {
	GetById(c context.Context, collectionName string, id interface{}, dest interface{}) error
}
type cacheRepo struct {
	cache                 CacheFunctions
	sessionCollectionName string
}

func NewCacheRepo(cache CacheFunctions,collectionName string) *cacheRepo {
	return &cacheRepo{cache: cache,sessionCollectionName: collectionName}
}

func (r cacheRepo) GetUserSessionByTokenToStruct(c context.Context,token string,dest interface{}) error {
	if err := r.cache.GetById(c , r.sessionCollectionName , token,&dest);err != nil {
		return err
	}
	return nil
}