package gitget

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetRefs(t *testing.T) {
	ctx := context.Background()

	urls := []string{
		// raw url that doe snot point to a git repo
		"https://raw.githubusercontent.com/deepmap/oapi-codegen/master/examples/petstore-expanded/petstore-expanded.yaml",

		// git repos
		"git@github.com:22:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"git@github.com:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"ssh://git@github.com:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"ssh://git@github.com:22:deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
		"https://github.com/deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml",
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

	for _, fp := range local {
		_, err := parseUrl(fp)
		require.Error(t, err)

		_, err = readLocalFile(fp)
		require.NoError(t, err)
	}

	for _, u := range urls {
		data, err := GetFile(ctx, u)
		require.NoError(t, err)
		require.NotEmpty(t, string(data))
		require.Contains(t, string(data), "openapi")

	}

}
