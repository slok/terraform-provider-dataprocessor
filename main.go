package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/slok/terraform-provider-dataprocessor/internal/provider"
)

const providerName = "registry.terraform.io/slok/dataprocessor"

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := providerserver.Serve(ctx, provider.New, providerserver.ServeOpts{
		Address: providerName,
	})

	return err
}

func main() {
	err := run(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running Terraform provider: %s", err)
		os.Exit(1)
	}
}
