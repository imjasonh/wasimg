# `wasimg`

`wasimg` bundles a Wasm module into an OCI manifest, and pushes it to an OCI registry.

For example:

```
$ wasimg my-module.wasm example.com/my-project/my-module
example.com/my-project/my-module@sha256:a7bb950a6cf95fd1dfc55907ec997e37840c71f6840d0c32481e7b8392490022
```

It doesn't require Docker or Dockerfiles, and it reuses your pre-configured registry credentials by default.

### Installation

```
go install github.com/imjasonh/wasimg@latest
```

### Usage

Build a wasm module, for example, build `wasimg` itself:

```
GOOS=wasip1 GOARCH=wasm go build -o wasimg.wasm .
```

If you already have a wasm module, you can skip this step.

Then use `wasimg` to bundle that wasm module and push it to a registry:

```
$ wasimg wasimg.wasm ttl.sh/wasimg
ttl.sh/wasimg@sha256:caea81fc44d4d92280a4bc7ceaccf15b3466792c312c9fa38446f73ce358ee3c
```

This prints the image reference of the pushed image, which you can run:

```
$ docker run \
  --runtime=io.containerd.wasmedge.v1 \
  --platform=wasip1/wasm \
  ttl.sh/wasimg@sha256:caea81fc44d4d92280a4bc7ceaccf15b3466792c312c9fa38446f73ce358ee3c
```

Try it out to see what happens!
