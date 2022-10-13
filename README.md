# `wasimg`

`wasimg` bundles a Wasm module into an image and pushes it to an OCI registry.

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
