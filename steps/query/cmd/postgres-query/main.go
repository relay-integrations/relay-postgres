package main

import (
	"context"
	"errors"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	"github.com/puppetlabs/relay-sdk-go/pkg/log"
	"github.com/puppetlabs/relay-sdk-go/pkg/outputs"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"

	"github.com/relay-integrations/relay-postgres/steps/query/pkg/query"
)

const (
	DefaultOutputKey    = "results"
	DefaultQueryTimeout = 5 * time.Minute
)

type ConnectionSpec struct {
	URL string `json:"url"`
}

type Spec struct {
	Connection ConnectionSpec `json:"connection"`
	Statement  string         `json:"statement"`
}

func main() {
	spec, err := populateSpec()
	if err != nil {
		log.FatalE(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, DefaultQueryTimeout)
	defer cancel()

	err = retry.Wait(ctx, func(ctx context.Context) (bool, error) {
		dberr := run(spec)

		if dberr != nil {
			return retry.Repeat(dberr)
		}

		return retry.Done(nil)
	})

	if err != nil {
		log.FatalE(err)
	}
}

func populateSpec() (*Spec, error) {
	specURL, err := taskutil.MetadataSpecURL()
	if err != nil {
		return nil, err
	}

	opts := taskutil.DefaultPlanOptions{SpecURL: specURL}

	spec := &Spec{}
	if err := taskutil.PopulateSpecFromDefaultPlan(&spec, opts); err != nil {
		return nil, err
	}

	if spec.Connection.URL == "" {
		return nil, errors.New("spec is missing a connection URL")
	}

	if spec.Statement == "" {
		return nil, errors.New("spec is missing a statement to execute")
	}

	return spec, nil
}

func run(spec *Spec) error {
	if runner, err := query.New(spec.Connection.URL); err != nil {
		return err
	} else {
		defer runner.Close()

		res, err := runner.Query(spec.Statement)

		if err != nil {
			return err
		}

		if client, err := outputs.NewDefaultOutputsClientFromNebulaEnv(); err != nil {
			return err
		} else {
			if err := client.SetOutput(context.Background(), DefaultOutputKey, res); err != nil {
				return err
			}
		}

	}

	return nil
}
