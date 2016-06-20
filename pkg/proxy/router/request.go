// Copyright 2016 CodisLabs. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package router

import (
	"sync"

	"github.com/CodisLabs/codis/pkg/proxy/redis"
	"github.com/CodisLabs/codis/pkg/utils/errors"
	"github.com/CodisLabs/codis/pkg/utils/sync2/atomic2"
)

type Dispatcher interface {
	Dispatch(r *Request) error
}

var ErrDiscardedRequest = errors.New("discarded request")

type Request struct {
	OpStr string
	Start int64

	Multi []*redis.Resp
	Batch *sync.WaitGroup

	batch sync.WaitGroup

	Coalesce func() error
	Response struct {
		Resp *redis.Resp
		Err  error
	}

	Broken *atomic2.Bool
	slot   *sync.WaitGroup
}

func NewRequest() *Request {
	r := &Request{}
	r.Batch = &r.batch
	return r
}

func (r *Request) Break() {
	if r.Broken != nil {
		r.Broken.Set(true)
	}
}

func (r *Request) IsBroken() bool {
	return r.Broken != nil && r.Broken.Get()
}
