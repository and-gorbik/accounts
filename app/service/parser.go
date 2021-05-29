package service

import (
	"accounts/util"
)

type parser struct {
	err       error
	strValues []string
	results   []interface{}
}

func NewParser(values []string) *parser {
	return &parser{
		strValues: values,
		results:   make([]interface{}, 0, len(values)),
	}
}

func (p *parser) OnlyValue() *parser {
	if p.err != nil {
		return p
	}

	if len(p.strValues) != 1 {
		p.err = errValuesLen
	}

	return p
}

func (p *parser) Int() *parser {
	if p.err != nil {
		return p
	}

	for _, value := range p.strValues {
		res, err := util.ParseInt(value)
		if err != nil {
			p.err = err
			break
		}

		p.results = append(p.results, res)
	}

	return p
}

func (p *parser) String() *parser {
	if p.err != nil {
		return p
	}

	for _, value := range p.strValues {
		p.results = append(p.results, value)
	}

	return p
}

func (p *parser) Bool() *parser {
	if p.err != nil {
		return p
	}

	for _, value := range p.strValues {
		intVal, err := util.ParseInt(value)
		if err != nil {
			p.err = err
			break
		}

		if intVal != 0 && intVal != 1 {
			p.err = errInvalidValue
			break
		}

		p.results = append(p.results, intVal == 1)
	}

	return p
}

func (p *parser) Timestamp() *parser {
	if p.err != nil {
		return p
	}

	for _, value := range p.strValues {
		res, err := util.ParseTimestamp(value)
		if err != nil {
			p.err = err
			break
		}

		p.results = append(p.results, res)
	}

	return p
}

func (p *parser) Parse() ([]interface{}, error) {
	if p.err != nil {
		return nil, p.err
	}

	return p.results, nil
}
