package namespace

type Finder interface {
	RegisterService(serviceName string, addr string)
	GetService(method string) (addr string, err error)
}
