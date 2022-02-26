package sentinel

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
	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/pkg/errors"
)

type propertyDataSource struct {
	datasource.Base
	rules string
}

func loadRules(rules string, newDatasource func(rules string) (datasource.DataSource, error)) error {
	if rules != "" {
		ds, err := newDatasource(rules)
		if err != nil {
			return err
		}

		err = ds.Initialize()
		if err != nil {
			return err
		}
	}

	return nil
}

func newFlowRuleDataSource(rules string) (datasource.DataSource, error) {
	return newDataSource(rules, datasource.NewFlowRulesHandler(datasource.FlowRuleJsonArrayParser))
}

func newCircuitBreakerRuleDataSource(rules string) (datasource.DataSource, error) {
	return newDataSource(rules, datasource.NewCircuitBreakerRulesHandler(datasource.CircuitBreakerRuleJsonArrayParser))
}

func newHotSpotParamRuleDataSource(rules string) (datasource.DataSource, error) {
	return newDataSource(rules, datasource.NewHotSpotParamRulesHandler(datasource.HotSpotParamRuleJsonArrayParser))
}

func newIsolationRuleDataSource(rules string) (datasource.DataSource, error) {
	return newDataSource(rules, datasource.NewIsolationRulesHandler(datasource.IsolationRuleJsonArrayParser))
}

func newSystemRuleDataSource(rules string) (datasource.DataSource, error) {
	return newDataSource(rules, datasource.NewSystemRulesHandler(datasource.SystemRuleJsonArrayParser))
}

func newDataSource(rules string, handlers ...datasource.PropertyHandler) (datasource.DataSource, error) {
	ds := &propertyDataSource{
		rules: rules,
	}
	for _, h := range handlers {
		ds.AddPropertyHandler(h)
	}

	return ds, nil
}

func (p propertyDataSource) ReadSource() ([]byte, error) {
	return []byte(p.rules), nil
}

func (p propertyDataSource) Initialize() error {
	src, err := p.ReadSource()
	if err != nil {
		err = errors.Errorf("Fail to read source, err: %+v", err)

		return err
	}

	return p.Handle(src)
}

func (p propertyDataSource) Close() error {
	// no op
	return nil
}
