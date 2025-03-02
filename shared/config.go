package shared

import (
	"github.com/pincher95/esctl/internal/client"
)

var (
	Client                *client.Client
	Context               string
	ElasticsearchProtocol string
	ElasticsearchUsername string
	ElasticsearchPassword string
	ElasticsearchHost     string
	ElasticsearchPort     int
	Debug                 bool
)
