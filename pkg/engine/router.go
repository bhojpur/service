package engine

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"sync"

	"github.com/bhojpur/service/pkg/engine/config"
	engine "github.com/bhojpur/service/pkg/engine/core"
	"github.com/bhojpur/service/pkg/engine/logger"
)

// router
type router struct {
	config *config.WorkflowConfig
}

func newRouter(config *config.WorkflowConfig) engine.Router {
	return &router{config: config}
}

// router interface
func (r *router) Route(appID string) engine.Route {
	logger.Debugf("%sapp[%s] workflowconfig is %#v", processorLogPrefix, appID, r.config)
	return newRoute(r.config)
}

func (r *router) Clean() {
	r.config = nil
}

// route interface
type route struct {
	data sync.Map
}

func newRoute(config *config.WorkflowConfig) *route {
	if config == nil {
		logger.Errorf("%sworkflowconfig is nil", processorLogPrefix)
		return nil
	}
	r := route{
		data: sync.Map{},
	}
	logger.Debugf("%sworkflowconfig %+v", processorLogPrefix, *config)
	for i, app := range config.Functions {
		r.Add(i, app.Name)
	}

	return &r
}

func (r *route) Add(index int, name string) {
	logger.Debugf("%sroute add: %s", processorLogPrefix, name)
	r.data.Store(index, name)
}

func (r *route) Exists(name string) bool {
	var ok bool
	logger.Debugf("%srouter[%v] exists name: %s", processorLogPrefix, r, name)
	r.data.Range(func(key interface{}, val interface{}) bool {
		if val.(string) == name {
			ok = true
			return false
		}
		return true
	})

	return ok
}

func (r *route) GetForwardRoutes(current string) []string {
	idx := -1
	r.data.Range(func(key interface{}, val interface{}) bool {
		if val.(string) == current {
			idx = key.(int)
			return false
		}
		return true
	})

	routes := make([]string, 0)
	r.data.Range(func(key interface{}, val interface{}) bool {
		if key.(int) > idx {
			routes = append(routes, val.(string))
		}
		return true
	})

	return routes
}
