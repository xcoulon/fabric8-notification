package collector

import "context"

type ReceiverResolver func(context.Context, string) (users []Receiver, templateValues map[string]interface{}, err error)
type ParamValidator func(context.Context, map[string]interface{}) error

type Receiver struct {
	FullName string
	EMail    string
}

type Registry interface {
	Get(string) (ReceiverResolver, bool)
	Validator(string) (ParamValidator, bool)
}

type LocalRegistry struct {
	reg        map[string]ReceiverResolver
	validators map[string]ParamValidator
}

func (r *LocalRegistry) Get(nType string) (ReceiverResolver, bool) {
	res, b := r.reg[nType]
	return res, b
}

func (r *LocalRegistry) Validator(nType string) (ParamValidator, bool) {
	v, found := r.validators[nType]
	return v, found
}

func (r *LocalRegistry) Register(nType string, resolver ReceiverResolver, validator ParamValidator) {
	if r.reg == nil {
		r.reg = map[string]ReceiverResolver{}
	}
	r.reg[nType] = resolver
	if r.validators == nil {
		r.validators = map[string]ParamValidator{}
	}
	if validator != nil {
		r.validators[nType] = validator
	}
}
