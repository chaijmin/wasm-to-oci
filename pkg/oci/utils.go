package oci

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	auth "oras.land/oras-go/pkg/auth/docker"
	orascnt "oras.land/oras-go/pkg/content"
	orasctx "oras.land/oras-go/pkg/context"
)

func newORASContext(insecure, useHTTP bool) (context.Context, remotes.Resolver, *orascnt.Memorystore) {
	ctx := orasctx.Background()
	memoryStore := orascnt.NewMemoryStore()
	cli, err := auth.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading auth file: %v\n", err)
	}

	client := http.DefaultClient
	if insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	resolver, err := cli.Resolver(context.Background(), client, useHTTP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading resolver: %v\n", err)
		resolver = docker.NewResolver(docker.ResolverOptions{})
	}

	return ctx, resolver, memoryStore
}
