
# Development

### Requirements

The requirements for building the operator are fairly minimal.

 * Go 1.13+
 * Operator SDK 1.0+
 * Bash or equivalent
 * Make
 * Podman

By default, the project uses [Podman][podman_link] for building container images, this can be changed to `docker` or `buildah` by setting the `IMAGE_BUILDER` variable to the tool of choice.

### Building from Source

There are several targets provided in the Makefile to build and release the operator binaries from source.

#### Build

Run the provided target to build the operator. A container image wil be created locally.

``` bash
make image-build
```

### Release

Push a locally created container image to a container registry for deployment.

``` bash
make image-push
```

[podman_link]:https://podman.io
