package namespace

type NamespaceCenter struct {
	// Map或本地共享文件 考虑并发读写
}

func (n *NamespaceCenter) RegisterService(method string, addr string) {
	return
}

func (n *NamespaceCenter) GetService(method string) (addr string, err error) {
	return
}
