package sessionresolver

import (
	"bitbucket.org/HeilaSystems/session"
	"container/list"
	"fmt"
)

type SessionResolverConfig struct {
	Repo session.SessionRepo
}

type defaultSessionResolver struct {
	ll *list.List
}

func Builder() session.SessionResolverBuilder {
	return &defaultSessionResolver{ll: list.New()}
}

func (cr *defaultSessionResolver) SetRepo(repo session.SessionRepo) session.SessionResolverBuilder {
	cr.ll.PushBack(func(cfg *SessionResolverConfig){
		cfg.Repo =  repo
	})
	return cr
}

func (cr *defaultSessionResolver) Build() (session.SessionResolver, error) {
	sessionCfg := &SessionResolverConfig{}
	for e := cr.ll.Front(); e != nil; e = e.Next() {
		f := e.Value.(func(cfg *SessionResolverConfig))
		f(sessionCfg)
	}
	if sessionCfg.Repo == nil {
		return nil, fmt.Errorf("cannot initalize configurations without repo")
	}
	return &sessionWrapper{repo: sessionCfg.Repo},nil
}