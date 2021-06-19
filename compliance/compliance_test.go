package compliance

import (
	"context"
	"strings"
	"testing"

	"github.com/openconfig/gribigo/device"
	"github.com/openconfig/gribigo/negtest"
	"github.com/openconfig/gribigo/testcommon"
)

func TestCompliance(t *testing.T) {
	for _, tt := range TestSuite {
		t.Run(tt.In.ShortName, func(t *testing.T) {
			creds, err := device.TLSCredsFromFile(testcommon.TLSCreds())
			if err != nil {
				t.Fatalf("cannot load credentials, got err: %v", err)
			}
			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()
			d, err := device.New(ctx, creds)

			if err != nil {
				t.Fatalf("cannot start server, %v", err)
			}

			if tt.FatalMsg != "" {
				if got := negtest.ExpectFatal(t, func(t testing.TB) {
					tt.In.Fn(d.GRIBIAddr(), t)
				}); !strings.Contains(got, tt.FatalMsg) {
					t.Fatalf("did not get expected fatal error, got: %s, want: %s", got, tt.FatalMsg)
				}
				return
			}

			if tt.ErrorMsg != "" {
				if got := negtest.ExpectError(t, func(t testing.TB) {
					tt.In.Fn(d.GRIBIAddr(), t)
				}); !strings.Contains(got, tt.ErrorMsg) {
					t.Fatalf("did not get expected error, got: %s, want: %s", got, tt.ErrorMsg)
				}
			}

			// Any unexpected error will be caught by being called directly on t from the fluent library.
			tt.In.Fn(d.GRIBIAddr(), t)
		})
	}
}