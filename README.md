# `wasimg`

`wasimg` bundles a Wasm module into an OCI manifest, and pushes it to an OCI registry.

For example:

```
$ wasimg my-module.wasm gcr.io/my-project/my-module
gcr.io/my-project/my-module@sha256:a7bb950a6cf95fd1dfc55907ec997e37840c71f6840d0c32481e7b8392490022
```

It doesn't require Docker or Dockerfiles, and it reuses your pre-configured registry credentials by default.

### Installation

```
go install github.com/imjasonh/wasimg@latest
```

### Details

The OCI manifest looks like this:

```
{
  "schemaVersion": 2,
  "mediaType": "application/vnd.oci.image.manifest.v1+json",
  "config": {
    "mediaType": "application/vnd.oci.image.config.v1+json",
    "size": 261,
    "digest": "sha256:e3c7550643da8b5b9954e479c2ea6aa5fc679896d9ee35511f9e4056f3343af0"
  },
  "layers": [
    {
      "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
      "size": 23,
      "digest": "sha256:81da0491c5af5635831f6a3febb5d9bfd66987ba3ecc42e58dc3d80938c25705"
    }
  ]
}
```

...and the config looks like this:

```
{
  "architecture": "wasm",
  "created": "0001-01-01T00:00:00Z",
  "history": [
    {
      "created": "0001-01-01T00:00:00Z"
    }
  ],
  "os": "wasi",
  "rootfs": {
    "type": "",
    "diff_ids": [
      "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    ]
  },
  "config": {
    "Entrypoint": [
      "my-module.wasm"
    ]
  }
}
```
