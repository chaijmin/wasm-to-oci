package oci

import (
	"io/ioutil"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/oras-project/oras-go/pkg/oras"
	log "github.com/sirupsen/logrus"
)

// Push pushes a WASM module an OCI registry
func Push(ref, module string, insecure, useHTTP bool) error {
	ctx, resolver, store := newORASContext(insecure, useHTTP)

	contents, err := ioutil.ReadFile(module)
	if err != nil {
		return err
	}

	desc := store.Add(module, ContentLayerMediaType, contents)
	layers := []ocispec.Descriptor{desc}

	pushOpts := []oras.PushOpt{
		oras.WithConfigMediaType(ConfigMediaType),
		oras.WithNameValidation(nil),
	}

	manifest, err := oras.Push(ctx, resolver, ref, store, layers, pushOpts...)
	if err != nil {
		return err
	}

	log.Infof("Pushed: %v", ref)
	log.Infof("Size: %v", desc.Size)
	log.Infof("Digest: %v", manifest.Digest)

	return nil
}
