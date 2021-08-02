package service

import (
	"time"

	"github.com/viant/mly/shared/datastore"
)

type Option interface {
	Apply(c *Service)
}

type storerOption struct {
	storer datastore.Storer
}

func (o *storerOption) Apply(c *Service) {
	c.datastore = o.storer
	c.useDatastore = true
}

func WithDataStorer(storer datastore.Storer) Option {
	return &storerOption{storer: storer}
}

type maxEvaluatorWait struct {
	time.Duration
}

func