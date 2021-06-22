package emq

type Client struct {
	*StdBroker
}

func (c *Client) Publish(topic string, msg interface{}) error {

}

func (c *Client) Subscribe(topic string) (<-chan interface{}, error) {

}

func (c *Client) Unsubscribe(topic string, sub <-chan interface{}) error {

}

func (c *Client) Close() {

}

func (c *Client) Broadcast(msg interface{}, subscribers []chan interface{}) {

}

func (c *Client) SetConditions(cap int) {
	c.setConditions(cap)
}