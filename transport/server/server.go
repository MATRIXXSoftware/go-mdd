package server

type TLS struct {
	Enable         bool
	SelfSignedCert bool
	CertFile       string
	KeyFile        string
}

type Options struct {
	Tls TLS
}

type Option func(*Options)

func DefaultOptions() Options {
	return Options{
		Tls: TLS{
			Enable:         false,
			SelfSignedCert: false,
			CertFile:       "",
			KeyFile:        "",
		},
	}
}

func WithTls(tls TLS) func(*Options) {
	return func(o *Options) {
		o.Tls = tls
	}
}
