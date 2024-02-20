package pagerduty

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	pd "github.com/PagerDuty/terraform-provider-pagerduty/pagerduty"
)

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("PAGERDUTY_PARALLEL"); v != "" {
		t.Parallel()
	}

	if v := os.Getenv("PAGERDUTY_TOKEN"); v == "" {
		t.Fatal("PAGERDUTY_TOKEN must be set for acceptance tests")
	}

	if v := os.Getenv("PAGERDUTY_USER_TOKEN"); v == "" {
		t.Fatal("PAGERDUTY_USER_TOKEN must be set for acceptance tests")
	}
}

func testAccCheckAttributes(n string, fn func(map[string]string) error) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes
		return fn(a)
	}
}

func testAccProtoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"pagerduty": func() (tfprotov5.ProviderServer, error) {
			ctx := context.Background()
			providers := []func() tfprotov5.ProviderServer{
				pd.Provider().GRPCProvider,
				providerserver.NewProtocol5(New()),
			}

			muxServer, err := tf5muxserver.NewMuxServer(ctx, providers...)
			if err != nil {
				return nil, err
			}

			return muxServer.ProviderServer(), nil
		},
	}
}
