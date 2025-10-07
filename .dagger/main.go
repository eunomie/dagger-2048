package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/dagger/internal/dagger"

	"github.com/containerd/platforms"
)

type Dagger struct {
	// The source code directory
	// +private
	Src *dagger.Directory
	// +private
	Platform *dagger.Platform
	// +private
	GoVersion string
	// +private
	AlpineVersion string
}

// Creates a new Dagger2048 Dagger module instance
func New(
	ctx context.Context,
// Source directory of the application.
// +optional
// +defaultPath="/"
	src *dagger.Directory,
// +optional
	platform *dagger.Platform,
// Go version to build the game.
// You can define it in a `.env` file
	goVersion string,
// Alpine version.
// You can define it in a `.env` file
	alpineVersion string,
) (*Dagger, error) {
	if platform == nil {
		if p, err := dag.DefaultPlatform(ctx); err != nil {
			return nil, err
		} else {
			platform = &p
		}
	}
	return &Dagger{
		Src:           src,
		Platform:      platform,
		GoVersion:     goVersion,
		AlpineVersion: alpineVersion,
	}, nil
}

// Build environment with Go tools and dependencies
func (m *Dagger) BuildEnv() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("golang:%s-alpine%s", m.GoVersion, m.AlpineVersion)).
		WithWorkdir("/app").
		With(func(c *dagger.Container) *dagger.Container {
			if m.Platform == nil {
				return c
			}
			spec := platforms.Normalize(platforms.MustParse(string(*m.Platform)))
			c = c.
				WithEnvVariable("GOOS", spec.OS).
				WithEnvVariable("GOARCH", spec.Architecture)
			switch spec.Architecture {
			case "arm", "arm64":
				switch spec.Variant {
				case "", "v8":
				default:
					c = c.WithEnvVariable("GOARM", strings.TrimPrefix(spec.Variant, "v"))
				}
			}
			return c
		}).
		WithMountedCache("/go/pkg/mod", dag.CacheVolume("go-mod")).
		WithMountedCache("/root/.cache/go-build", dag.CacheVolume("go-build")).
		WithDirectory(
			".",
			m.Src,
			dagger.ContainerWithDirectoryOpts{
				Include: []string{
					"go.mod",
					"go.sum",
				},
			}).
		WithExec([]string{"go", "mod", "download"}).
		WithDirectory(".", m.Src.WithoutDirectory("web")).
		WithEnvVariable("CGO_ENABLED", "0")
}

// Build Go application
func (m *Dagger) Build() *dagger.Container {
	return m.BuildEnv().
		WithExec([]string{"go", "build", "-ldflags", "-s -w", "-o", "dagger2048", "./main.go"})
}

// Return the built binary
//
// To return a binary for the current platform, run the following command:
//
//	dagger -c '. --platform current | binary | export dagger2048'
//
// Then you can run
//
//	./dagger2048
func (m *Dagger) Binary() *dagger.File {
	return m.Build().File("dagger2048")
}

// Run Go tests
func (m *Dagger) Test(ctx context.Context) (string, error) {
	ctr := m.BuildEnv().
		WithExec([]string{"go", "test", "./..."}, dagger.ContainerWithExecOpts{Expect: dagger.ReturnTypeAny})
	out, err := ctr.CombinedOutput(ctx)
	if err != nil {
		return "", err
	}
	if e, err := ctr.ExitCode(ctx); err != nil {
		return "", err
	} else if e != 0 {
		return "", fmt.Errorf("backend tests failed:\n%s", out)
	}
	return out, nil
}

// Create runnable docker image
func (m *Dagger) Image() *dagger.Container {
	return dag.Container().
		From(fmt.Sprintf("alpine:%s", m.AlpineVersion)).
		WithWorkdir("/app").
		WithFile("/app/dagger2048", m.Build().File("/app/dagger2048")).
		WithDefaultArgs([]string{"/app/dagger2048"})
}

// Run the binary inside a container
func (m *Dagger) Run() *dagger.Container {
	return m.Image().Terminal(dagger.ContainerTerminalOpts{
		Cmd: []string{"/app/dagger2048"},
	})
}
