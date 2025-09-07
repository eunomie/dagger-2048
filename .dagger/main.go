package main

import (
	"context"
	"fmt"
	"strings"

	"dagger/dagger-2048/internal/dagger"
	"github.com/containerd/platforms"
)

type Dagger2048 struct {
	// The source code directory
	Src      *dagger.Directory
	Platform *dagger.Platform
}

// Creates a new Dagger2048 Dagger module instance
func New(
// Source directory of the application
// +optional
// +defaultPath="/"
	src *dagger.Directory,
// +optional
	platform *dagger.Platform,
) *Dagger2048 {
	return &Dagger2048{
		Src:      src,
		Platform: platform,
	}
}

// Build environment with Go tools and dependencies
func (m *Dagger2048) BuildEnv() *dagger.Container {
	return dag.Container().
		From("golang:1.25-alpine3.22").
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
func (m *Dagger2048) Build() *dagger.Container {
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
func (m *Dagger2048) Binary() *dagger.File {
	return m.Build().File("dagger2048")
}

// Run Go tests
func (m *Dagger2048) Test(ctx context.Context) (string, error) {
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
func (m *Dagger2048) Image() *dagger.Container {
	return dag.Container().
		From("alpine:3.22").
		WithWorkdir("/app").
		WithFile("/app/dagger2048", m.Build().File("/app/dagger2048")).
		WithDefaultArgs([]string{"/app/dagger2048"})
}

// Run the binary inside a container
func (m *Dagger2048) Run() *dagger.Container {
	return m.Image().Terminal(dagger.ContainerTerminalOpts{
		Cmd: []string{"/app/dagger2048"},
	})
}
