package client

type TLS struct {
	Enable             bool
	InsecureSkipVerify bool
	CertFile           string
}

type Options struct {
	Tls TLS
}

type Option func(*Options)

func DefaultOptions() Options {
	return Options{
		Tls: TLS{
			Enable:             false,
			InsecureSkipVerify: false,
			CertFile:           "",
		},
	}
}

func WithTls(tls TLS) func(*Options) {
	return func(o *Options) {
		o.Tls = tls
	}
}
