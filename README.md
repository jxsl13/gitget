# gitget

`gitget` is a library to fetch single files from remote or local locations, primarily git repositories.


## import

```
go get github.com/jxsl13/gitget
```

## supported uri formats

```go
remote := []string{
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

local := []string{
    "test/openapi.yaml",
    "file://test/openapi.yaml",
    "./test/openapi.yaml",
    "file://./test/openapi.yaml",
}
```

## examples

```go
import (
    "fmt"
    "github.com/jxsl13/gitget"
)
func main() {
    data, err := gitget.Get("https://github.com/deepmap/oapi-codegen.git/examples/petstore-expanded/petstore-expanded.yaml@master")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
}
```