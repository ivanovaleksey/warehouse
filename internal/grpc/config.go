package grpc

import "fmt"

type Config struct {
	Port int
}

func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}
