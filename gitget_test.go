package gitget

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	ctx := context.Background()

	urls := []string{
		// raw url that does not point to a git repo but to a concrete file
		"https://raw.githubusercontent.com/deepmap/oapi-codegen/master/examples/petstore-expanded/petstore-expanded.yaml",

		// git repos (github flavor)
		// only with a known ssh key
		"git@github.com:22:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"git@github.com:22:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml@master",

		"git@github.com:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"git@github.com:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml@master",

		"ssh://git@github.com:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"ssh://git@github.com:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml@master",

		"ssh://git@github.com:22:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"ssh://git@github.com:22:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml@master",

		// works without any ssh keys
		"https://github.com/deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"https://github.com/deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml@master",

		// git repos (bitbucket flavor)
		"https://bitbucket.org/atlassian/openapi-diff.git/test/e2e/fixtures/openapi3/source-spec.json",

		// only with a known ssh key (bitbucket must know your public key)
		"ssh://git@bitbucket.org/atlassian/openapi-diff.git/test/e2e/fixtures/openapi3/source-spec.json",
	}

	for _, u := range urls {
		_, err := parseUrl(u)
		require.NoError(t, err)
	}

	// local files
	// TODO: might also point to local directories containing git repos
	local := []string{
		"test/openapi.yaml",
		"file://test/openapi.yaml",
		"./test/openapi.yaml",
		"file://./test/openapi.yaml",
	}

	for idx, fp := range local {
		_, err := parseUrl(fp)
		require.Error(t, err)

		_, err = readLocalFile(fp)
		require.NoErrorf(t, err, "local error at #%d", idx+1)
	}

	for idx, u := range urls {
		data, err := Get(ctx, u, GetOptions{InsecureSkipTLS: true})
		require.NoErrorf(t, err, "urls error at #%d", idx+1)
		require.NotEmpty(t, string(data))
		require.Contains(t, string(data), "openapi")

	}

}
