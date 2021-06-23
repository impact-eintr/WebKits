package emq

type Client struct {
	*StdBroker
}

func NewClient() *Client {
	return &Client{NewBroker()}
}

func (c *Client) Publish(topic string, msg interface{}) error {
	return c.publish(topic, msg)
}

func (c *Client) Subscribe(topic string) (<-chan interface{}, error) {
	return c.subscribe(topic)
}

func (c *Client) Unsubscribe(topic string, sub <-chan interface{}) error {
	return c.unsubscribe(topic, sub)
}

func (c *Client) Close() {
	c.close()
}

func (c *Client) SetConditions(cap int) {
	c.setConditions(cap)
}

func (c *Client) GetPayLoad(sub <-chan interface{}) interface{} {
	for val := range sub {
		if val != nil {
			return val
		}
	}
	return nil
}
