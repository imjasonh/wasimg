/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

var variant = flag.String("variant", "", "variant to set (spin, slight)")

func main() {
	flag.Parse()

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
		OS:           "wasi",
		Architecture: "wasm",
		Variant:      *variant,
		Config: v1.Config{
			Entrypoint: []string{fn},
		},
	})
	if err != nil {
		log.Fatalf("mutate.ConfigFile: %v", err)
	}

	l, err := tarball.LayerFromFile(fn)
	if err != nil {
		log.Fatalf("tarball.LayerFromFile: %v", err)
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
