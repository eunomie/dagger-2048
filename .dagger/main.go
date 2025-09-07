package main

import (
	"context"
	"fmt"

	"dagger/dagger-2048/internal/dagger"
)

type Dagger2048 struct {
	// The source code directory
	Src *dagger.Directory
}

// Creates a new Dagger2048 Dagger module instance
func New(
// Source directory of the application
// +optional
// +defaultPath="/"
	src *dagger.Directory,
) *Dagger2048 {
	return &Dagger2048{
		Src: src,
	}
}

// Build environment with Go tools and dependencies
func (m *Dagger2048) BuildEnv() *dagger.Container {
	return dag.Container().
		From("golang:1.25-alpine3.22").
		WithWorkdir("/app").
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
