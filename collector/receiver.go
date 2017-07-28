package collector

import "context"

type ReceiverResolver func(context.Context, string) (users []Receiver, templateValues map[string]interface{}, err error)

type Receiver struct {
	FullName string
	EMail    string
}

type Registry interface {
	Get(string) (ReceiverResolver, bool)
}

type LocalRegistry struct {
	reg map[string]ReceiverResolver
}

func (r *LocalRegistry) Get(nType string) (ReceiverResolver, bool) {
	res, b := r.reg[nType]
	return res, b
}

func (r *LocalRegistry) Register(nType string, resolver ReceiverResolver) {
	if r.reg == nil {
		r.reg = map[string]ReceiverResolver{}
	}
	r.reg[nType] = resolver
}
