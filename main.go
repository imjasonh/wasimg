/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

func main() {
	flag.Parse()

	hello() // this is the function that is conditionally compiled

	if len(flag.Args()) != 2 {
		log.Fatalf("Usage: %s <module>.wasm <ref>", os.Args[0])
	}
	fn := flag.Arg(0)
	refstr := flag.Arg(1)

	ref, err := name.ParseReference(refstr)
	if err != nil {
		log.Fatalf("name.ParseReference: %v", err)
	}

	img := mutate.MediaType(empty.Image, types.OCIManifestSchema1).(v1.Image)
	img = mutate.ConfigMediaType(img, types.OCIConfigJSON).(v1.Image)
	img, err = mutate.ConfigFile(img, &v1.ConfigFile{
		OS:           "wasip1",
		Architecture: "wasm",
		Config: v1.Config{
			Entrypoint: []string{fn},
		},
	})
	if err != nil {
		log.Fatalf("mutate.ConfigFile: %v", err)
	}

	f, err := os.Open(fn)
	if err != nil {
		log.Fatalf("os.Open: %v", err)
	}
	fi, err := f.Stat()
	if err != nil {
		log.Fatalf("f.Stat: %v", err)
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	if err := tw.WriteHeader(&tar.Header{
		Name: filepath.Base(fn),
		Mode: 0755,
		Size: fi.Size(),
	}); err != nil {
		log.Fatalf("tw.WriteHeader: %v", err)
	}
	if _, err := io.Copy(tw, f); err != nil {
		log.Fatalf("io.Copy: %v", err)
	}
	if err := tw.Close(); err != nil {
		log.Fatalf("tw.Close: %v", err)
	}

	l, err := tarball.LayerFromReader(&buf)
	if err != nil {
		log.Fatalf("tarball.LayerFromReader: %v", err)
	}
	img, err = mutate.Append(img, mutate.Addendum{
		Layer:     l,
		MediaType: types.OCILayer,
	})
	if err != nil {
		log.Fatalf("mutate.Append: %v", err)
	}

	d, err := img.Digest()
	if err != nil {
		log.Fatalf("img.Digest: %v", err)
	}

	if err := remote.Write(ref, img, remote.WithAuthFromKeychain(authn.DefaultKeychain)); err != nil {
		log.Fatalf("Pushing image: %v", err)
	}
	fmt.Println(ref.Context().Digest(d.String()))
}
