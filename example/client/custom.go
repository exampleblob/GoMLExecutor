package client

type customMakerRegistry struct {
	registry map[string]func() interface{}
}

func (c *customMakerRegistry) Register(k string, gen func() interface{})