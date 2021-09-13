package main

import (
	"context"
	"flag"

	"github.com/puppetlabs/relay-sdk-go/pkg/log"
	"github.com/puppetlabs/relay-sdk-go/pkg/outputs"
	"github.com/puppetlabs/relay-sdk-go/pkg/taskutil"

	"github.com/relay-integrations/relay-postgres/steps/query/pkg/query"
)

// DefaultOutputKey is the key of the output that will be set when the step
// executes successfully.
const DefaultOutputKey = "results"

// ConnectionSpec contains the relevant connection information. In the Relay
// product a connection object will be created.
type ConnectionSpec struct {
	URL string `json:"url"`
}

// Spec is the schema for the actual spec.
type Spec struct {
	Connection ConnectionSpec `json:"connection"`
	Statement  string         `json:"statement"`
}

func main() {
	var (
		specURL = flag.String("spec-url", mustGetDefaultMetadataSpecURL(), "url to fetch the spec from")
	)

	flag.Parse()

	spec := mustPopulateSpec(*specURL)

	if spec.Connection.URL == "" {
		log.Fatal("spec is missing a connection URL")
	}

	if spec.Statement == "" {
		log.Fatal("spec is missing a statement to execute")
	}

	if runner, err := query.New(spec.Connection.URL); err != nil {
		log.FatalE(err)
	} else {
		defer runner.Close()

		res, err := runner.Query(spec.Statement)

		if err != nil {
			log.FatalE(err)
		}

		if client, err := outputs.NewDefaultOutputsClientFromNebulaEnv(); err != nil {
			log.FatalE(err)
		} else {
			if err := client.SetOutput(context.Background(), DefaultOutputKey, res); err != nil {
				log.FatalE(err)
			}
		}
	}
}

func mustGetDefaultMetadataSpecURL() string {
	if metadataSpecURL, err := taskutil.MetadataSpecURL(); err != nil {
		log.FatalE(err)

		panic(err)
	} else {
		return metadataSpecURL
	}
}

func mustPopulateSpec(specURL string) (spec Spec) {
	opts := taskutil.DefaultPlanOptions{SpecURL: specURL}

	if err := taskutil.PopulateSpecFromDefaultPlan(&spec, opts); err != nil {
		log.FatalE(err)
	}

	return
}
