package oauth2

type EmptyOption struct{}

func (EmptyOption) apply(*interface{}) (err error) { return }

type funcOption struct {
	f func(interface{}) error
}

func (fo *funcOption) apply(h interface{}) error {

	return fo.f(h)
}

func newFuncOption(f func(interface{}) error) *funcOption {
	return &funcOption{
		f: f,
	}
}
