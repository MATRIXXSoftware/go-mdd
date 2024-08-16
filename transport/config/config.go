package config

type ClientOptions struct {
	Tls                bool
	InsecureSkipVerify bool
	CertFile           string
}

type ClientOption func(*ClientOptions)

func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Tls:                false,
		InsecureSkipVerify: false,
		CertFile:           "",
	}
}

func WithTls() func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.Tls = true
	}
}

func WithInsecureSkipVerify() func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.InsecureSkipVerify = true
	}
}

func WithCertFile(certFile string) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.CertFile = certFile
	}
}
