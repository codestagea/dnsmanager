package cache

import (
	"github.com/codestagea/bindmgr/internal/model"
	"github.com/codestagea/bindmgr/internal/tools"
	"sync"
	"time"
)

var (
	ErrUserNotFound = tools.ErrorWithCode{3000, "user token not found"}
	ErrTokenExpire  = tools.ErrorWithCode{3000, "token expire"}
)

type TokenCache interface {
	GetLoginUser(token string) (*model.LoginUser, error)
	DelToken(token string) error
	SetLoginUser(user *model.LoginUser) error
	RefreshToken(user *model.LoginUser) error
}

type localTokenCache struct {
	stop chan struct{}

	wg    sync.WaitGroup
	mu    sync.RWMutex
	users map[string]*model.LoginUser
}

func NewLocalTokenCache(cleanupInterval time.Duration) TokenCache {
	lc := &localTokenCache{
		users: make(map[string]*model.LoginUser),
		stop:  make(chan struct{}),
	}

	lc.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer lc.wg.Done()
		lc.cleanupLoop(cleanupInterval)
	}(cleanupInterval)

	return lc
}

func (c *localTokenCache) GetLoginUser(token string) (*model.LoginUser, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	cu, ok := c.users[token]
	if !ok {
		return nil, ErrUserNotFound
	}

	if cu.ExpireTime.Before(time.Now()) {
		return nil, ErrTokenExpire
	}

	return cu, nil
}

func (c *localTokenCache) DelToken(token string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.users, token)
	return nil
}

func (c *localTokenCache) SetLoginUser(user *model.LoginUser) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users[user.Uuid] = user
	return nil
}

func (c *localTokenCache) RefreshToken(user *model.LoginUser) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.users[user.Uuid] = user
	return nil
}

func (c *localTokenCache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-c.stop:
			return
		case <-t.C:
			c.mu.Lock()
			now := time.Now()
			for uid, cu := range c.users {
				if cu.ExpireTime.Before(now) {
					delete(c.users, uid)
				}
			}
			c.mu.Unlock()
		}
	}
}

func (c *localTokenCache) stopCleanup() {
	close(c.stop)
	c.wg.Wait()
}
