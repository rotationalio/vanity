package vanity_test

import (
	"net/http"
	"net/url"
	"testing"

	. "github.com/rotationalio/vanity"
	"github.com/rotationalio/vanity/config"
	"github.com/stretchr/testify/require"
)

type expected struct {
	redirect   string
	importMeta string
	sourceMeta string
}

func TestGoPackage(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		testCases := []struct {
			in       *GoPackage
			conf     *config.Config
			req      *http.Request
			expected expected
		}{
			{
				in:   &GoPackage{Repository: "https://github.com/rotationalio/confire"},
				conf: &config.Config{Domain: "go.rotational.io", DefaultBranch: "main"},
				req:  &http.Request{URL: &url.URL{Path: "/confire/validate"}},
				expected: expected{
					redirect:   "https://godoc.org/go.rotational.io/confire/validate",
					importMeta: "go.rotational.io/confire git https://github.com/rotationalio/confire",
					sourceMeta: "go.rotational.io/confire https://github.com/rotationalio/confire https://github.com/rotationalio/confire/tree/main{/dir} https://github.com/rotationalio/confire/blob/main{/dir}/{file}#L{line}",
				},
			},
		}

		for i, tc := range testCases {
			// Resolve the go package data
			require.NoError(t, tc.in.Resolve(tc.conf), "test case %d failed: could not resolve", i)

			// Finalize package for the request
			pkg := tc.in.WithRequest(tc.req)

			// Perform assertions
			require.Equal(t, tc.expected.redirect, pkg.Redirect(), "test case %d failed: bad redirect", i)
			require.Equal(t, tc.expected.importMeta, pkg.GoImportMeta(), "test case %d failed: bad import meta", i)
			require.Equal(t, tc.expected.sourceMeta, pkg.GoSourceMeta(), "test case %d failed: bad source meta", i)
		}
	})
}
