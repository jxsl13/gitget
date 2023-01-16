package gitget

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/client"
	"github.com/go-git/go-git/v5/storage/memory"
)

var (
	ErrInvalidPath = errors.New("invalid path")
	ErrInvalidUri  = errors.New("invalid uri")
	errNotSshUrl   = errors.New("not an ssh url")
	errLocalPath   = errors.New("local path")
	errNotGitUrl   = errors.New("not a git url, doe snot contain '.git'")
)

type GetOptions struct {
	InsecureSkipTLS bool
}

func Get(ctx context.Context, gitUrl string, opts ...GetOptions) ([]byte, error) {
	var opt GetOptions
	if len(opts) > 0 {
		opt = opts[0]
	}

	repoUrl, filePath, branch, err := splitRepoUrl(gitUrl)
	if err != nil {

		if errors.Is(err, errNotGitUrl) {
			// but it is a url
			resp, err := http.Get(gitUrl)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			return io.ReadAll(resp.Body)
		}

		data, ferr := readLocalFile(gitUrl)
		if ferr != nil {
			return nil, fmt.Errorf("%w: %v", err, ferr)
		}

		// unexpected error, might be a file
		return data, nil
	}

	fs := memfs.New()
	storer := memory.NewStorage()

	o := &git.CloneOptions{
		URL:             repoUrl,
		SingleBranch:    true,
		ReferenceName:   plumbing.NewBranchReferenceName(branch),
		Depth:           1,
		InsecureSkipTLS: opt.InsecureSkipTLS,
	}

	_, err = git.Clone(storer, fs, o)
	if err != nil {
		fmt.Println(err)
	}

	f, err := fs.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error in cloned git repository: %w: %v", err, filePath)
	}
	defer f.Close()
	return io.ReadAll(f)
}

func readLocalFile(localPath string) ([]byte, error) {
	localPath = strings.TrimPrefix(localPath, "file://")

	// not an url, try to get fro mlocal path
	absolutePath, err := filepath.Abs(localPath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidPath, err)
	}
	data, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidUri, err)
	}
	return data, nil
}

func splitRepoUrl(fullUrl string) (repoUrl, filePath, branch string, err error) {
	defer func() {
		if err != nil {
			return
		}

		if branch == "" {
			e, err := transport.NewEndpoint(repoUrl)
			if err != nil {
				return
			}
			cli, err := client.NewClient(e)
			if err != nil {
				return
			}
			s, err := cli.NewUploadPackSession(e, nil)
			if err != nil {
				return
			}
			info, err := s.AdvertisedReferences()
			if err != nil {
				return
			}
			refs, err := info.AllReferences()
			if err != nil {
				return
			}
			headReference := refs["HEAD"].Target()
			headBranch := headReference.String()
			branch = strings.TrimPrefix(headBranch, "refs/heads/")
		}
	}()

	u, err := parseUrl(fullUrl)
	if err != nil {
		return "", "", "", err
	}

	if strings.Contains(u.Path, "@") {
		tokens := strings.Split(u.Path, "@")
		branch, err = last(tokens)
		if err != nil {
			return "", "", "", err
		}
		u.Path = strings.Join(tokens[:len(tokens)-1], "@")
	}

	pathParts := strings.Split(u.Path, "/")
	for idx, p := range pathParts {
		if strings.Contains(p, ".git") {
			if idx+1 > len(pathParts) {
				return "", "", "", errors.New("invalid git url: expected {https|ssh}://domain.tld/{name}.git/{path}[@{branch}]")
			}

			u.Path = path.Join(pathParts[:idx+1]...)
			filePath := path.Join(pathParts[idx+1:]...)
			if filePath == "" {
				return "", "", "", fmt.Errorf("invalid file path: %s", filePath)
			}
			return u.String(), filePath, branch, nil
		}
	}

	return "", "", "", errNotGitUrl
}

func last[T any](a []T) (T, error) {
	if len(a) == 0 {
		var t T
		return t, errors.New("slice is empty")
	}
	return a[len(a)-1], nil
}

func parseUrl(urlStr string) (*url.URL, error) {
	u, err := detectSSH(urlStr)
	if err == nil {
		return u, nil
	}

	u, err = url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "file" {
		return nil, errLocalPath
	}

	return u, nil
}

var sshPattern = regexp.MustCompile(`^(?:([a-zA-Z+-.]+)*://)?(?:([^@\s]+)@)([^:]+(?::\d+)?):(/?.+)$`)

func detectSSH(src string) (*url.URL, error) {
	matched := sshPattern.FindStringSubmatch(src)
	if len(matched) == 0 {
		return nil, errNotSshUrl
	}
	if matched[1] != "" && !strings.Contains(matched[1], "ssh") {
		return nil, fmt.Errorf("%w: %s", errNotSshUrl, src)
	}

	user := matched[2]
	host := matched[3]
	path := matched[4]
	qidx := strings.Index(path, "?")
	if qidx == -1 {
		qidx = len(path)
	}

	var u url.URL
	u.Scheme = "ssh"
	u.User = url.User(user)
	u.Host = host
	u.Path = path[0:qidx]
	if qidx < len(path) {
		q, err := url.ParseQuery(path[qidx+1:])
		if err != nil {
			return nil, fmt.Errorf("%w: %s", errNotSshUrl, err)
		}
		u.RawQuery = q.Encode()
	}

	return &u, nil
}
