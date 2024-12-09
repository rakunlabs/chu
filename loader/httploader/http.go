package httploader

import "github.com/worldline-go/klient"

type Loader struct {
	client *klient.Client
}

func New(opts ...klient.OptionClientFn) *Loader {
	opts = append([]klient.OptionClientFn{
		klient.WithDisableBaseURLCheck(true),
	}, opts...)

	client, err := klient.New(opts...)
	if err != nil {
		panic(err)
	}

	return &Loader{
		client: client,
	}
}
