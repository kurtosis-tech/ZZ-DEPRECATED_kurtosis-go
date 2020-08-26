package services

const (
	exampleServicePort = 80
)
type ExampleServiceImpl struct{
	IPAddr string
}

func (e ExampleServiceImpl) GetHelloWorldSocket() Socket {
	return Socket{
		IPAddr: e.IPAddr,
		Port: exampleServicePort,
	}
}
