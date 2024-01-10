package grafana

import (
	"context"

	"github.com/grafana-tools/sdk"
	log "github.com/sirupsen/logrus"
)

type Grafana struct {
	client *sdk.Client
	ctx    context.Context
	dir    string
}

func New(url, apiKey, dir string) Grafana {
	var err error

	gc := Grafana{
		ctx: context.Background(),
		dir: dir,
	}
	if gc.client, err = sdk.NewClient(url, apiKey, httpClient); err != nil {
		log.Fatalln(err)
	}
	return gc
}
