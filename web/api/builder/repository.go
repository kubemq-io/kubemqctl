package builder

import (
	"context"
	"fmt"
	"sync"
	"time"
)

const keepTime = 10 * time.Minute

type Repository struct {
	cache sync.Map
}

func NewRepository(ctx context.Context) *Repository {
	r := &Repository{}
	go r.run(ctx)
	return r
}

func (r *Repository) Set(m *Manifest) {
	r.cache.Store(m.Id, m)
}

func (r *Repository) Get(id int64) (*Manifest, error) {
	val, ok := r.cache.Load(id)
	if !ok {
		return nil, fmt.Errorf("invalid build id")
	}
	return val.(*Manifest), nil
}

func (r *Repository) run(ctx context.Context) {
	for {
		select {
		case <-time.After(10 * time.Second):
			r.cache.Range(func(key, value interface{}) bool {
				id, man := key.(int64), value.(*Manifest)
				if man.CreatedAt.Add(keepTime).Unix() < time.Now().Unix() {
					r.cache.Delete(id)
				}
				return true
			})
		case <-ctx.Done():
			return
		}
	}
}
