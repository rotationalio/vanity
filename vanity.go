package vanity

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"go.rtnl.ai/vanity/config"
	"gopkg.in/yaml.v3"
)

const (
	protocolGit    = "git"
	protocolGitHub = "github"
	protocolGOGS   = "gogs"
)

var godoc *url.URL = &url.URL{
	Scheme: "https",
	Host:   "godoc.org",
}

var validProtocols = map[string]struct{}{
	protocolGit:    {},
	protocolGitHub: {},
	protocolGOGS:   {},
}

type GoPackage struct {
	Domain     string   `yaml:"-"`          // the vanity URL domain to use
	Module     string   `yaml:"-"`          // the module name where go.mod is located; parsed from the repository
	Package    string   `yaml:"-"`          // the full package path being requested for correct redirects
	Protocol   string   `yaml:"protocol"`   // can be "git", "github", or "gogs" -- defaults to "git"
	Repository string   `yaml:"repository"` // a path to the public repository starting with https://
	Branch     string   `yaml:"branch"`     // the name of the default branch -- defaults to "main"
	repo       *url.URL `yaml:"-"`          // the parsed repository URL
	user       string   `yaml:"-"`          // the user or organization from the repository
}

func Load(conf *config.Config) (pkgs []*GoPackage, err error) {
	if err = conf.Validate(); err != nil {
		return nil, err
	}

	var f *os.File
	if f, err = os.Open(conf.ConfigMap); err != nil {
		return nil, fmt.Errorf("could not open %q: %w", conf.ConfigMap, err)
	}
	defer f.Close()

	pkgs = make([]*GoPackage, 0)
	if err = yaml.NewDecoder(f).Decode(&pkgs); err != nil {
		return nil, fmt.Errorf("could not decode vanity urls config: %w", err)
	}

	for i, pkg := range pkgs {
		if err = pkg.Resolve(conf); err != nil {
			return nil, fmt.Errorf("could not resolve config %d: %w", i+1, err)
		}
	}

	return pkgs, nil
}

func (p *GoPackage) Resolve(conf *config.Config) (err error) {
	// Verify there is a repository
	if p.Repository == "" {
		return ErrNoRepository
	}

	// Parse the repository
	if p.repo, err = url.Parse(p.Repository); err != nil {
		return ErrInvalidRepository
	}

	parts := strings.Split(p.repo.Path, "/")
	if len(parts) != 3 {
		return ErrInvalidRepository
	}

	p.user = parts[1]
	p.Module = parts[2]

	// Check protocol
	if p.Protocol == "" {
		p.Protocol = protocolGit
	}

	if _, ok := validProtocols[p.Protocol]; !ok {
		return ErrInvalidProtocol
	}

	// Manage the configuration
	if conf != nil {
		p.Domain = conf.Domain

		if p.Branch == "" {
			p.Branch = conf.DefaultBranch
		}
	}

	// Check the ref
	if p.Branch == "" {
		p.Branch = "main"
	}

	return nil
}

func (p *GoPackage) WithRequest(r *http.Request) GoPackage {
	pkg := p.Module
	if r != nil {
		pkg = r.URL.Path
	}

	clone := GoPackage{
		Domain:     p.Domain,
		Module:     p.Module,
		Package:    pkg,
		Protocol:   p.Protocol,
		Repository: p.Repository,
		Branch:     p.Branch,
		repo:       p.repo,
		user:       p.user,
	}

	if clone.Domain == "" && r != nil {
		clone.Domain = r.Host
	}

	return clone
}

func (p GoPackage) Redirect() string {
	var uri *url.URL
	if p.Package != "" {
		uri = godoc.ResolveReference(&url.URL{Path: filepath.Join("/", p.Domain, p.Package)})
	} else {
		uri = godoc.ResolveReference(&url.URL{Path: filepath.Join("/", p.Domain, p.Module)})
	}

	return uri.String()
}

func (p GoPackage) GoImportMeta() string {
	parts := []string{
		p.Import(),
		p.Protocol,
		p.repo.String(),
	}

	return strings.Join(parts, " ")
}

func (p GoPackage) GoSourceMeta() string {
	parts := []string{
		p.Import(),
		p.repo.String(),
		"",
		"",
	}
	parts[2], parts[3] = p.Source()
	return strings.Join(parts, " ")
}

func (p GoPackage) Import() string {
	return filepath.Join(p.Domain, p.Module)
}

func (p GoPackage) Source() (string, string) {
	switch p.Protocol {
	case protocolGit, protocolGitHub:
		return p.githubSource()
	case protocolGOGS:
		return p.gogsSource()
	default:
		return "", ""
	}
}

func (p GoPackage) githubSource() (string, string) {
	base := filepath.Join(p.user, p.Module)
	directoryPath := filepath.Join(base, "tree", p.Branch+"{/dir}")
	filePath := filepath.Join(base, "blob", p.Branch+"{/dir}", "{file}#L{line}")

	uri := p.repo.ResolveReference(&url.URL{Path: "/"}).String()
	return uri + directoryPath, uri + filePath
}

func (p GoPackage) gogsSource() (string, string) {
	base := filepath.Join(p.user, p.Module)
	directoryPath := filepath.Join(base, "src", p.Branch+"{/dir}")
	filePath := filepath.Join(base, "src", p.Branch+"{/dir}", "{file}#L{line}")

	uri := p.repo.ResolveReference(&url.URL{Path: "/"}).String()
	return uri + directoryPath, uri + filePath
}
