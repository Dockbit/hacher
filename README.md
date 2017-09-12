# Hacher :hocho:

```hacher``` is a simple CLI tool to cache project artifacts. It tracks the changes to one or more dependency files, and if any of them has changed, caches related artifacts directory. For example, you could cache ```bundle``` directory on [Gemfile](http://bundler.io/man/gemfile.5.html) changes or a ```node_modules``` when [package.json](https://docs.npmjs.com/files/package.json) changes.

## Installation

Install [Golang](https://golang.org/doc/install) and:

    make install

```hacher``` needs several environment variables to operate:

* ```HACHER_PATH``` - The path where cached artifacts will be stored.
* ```HACHER_KEEP``` - The number of caches to keep, defaults to ```3```. You probably won't need more, since it's only useful if someone reverts the dependency file, hence cache could be reused.

### Building in Docker

You can use the [Dockerfile](./Dockerfile) to download and install the hacher binary inside a Docker image. This is handy in case you don't have the Go environment setup locally.

Specifying the [version](https://github.com/Dockbit/hacher/releases) to install is available via the `RELEASE` Docker build argument.

```docker build --tag hacher --build-arg RELEASE=v0.1.0 .```

## Usage

Let's say you wanted to speed up ```npm install``` during your CI builds, so you could:

    # Get the cached version of node_modules to the current directory
    hacher get -k node_modules -f package.json

    npm install

    # Cache node_modules on package.json changes
    hacher set -k node_modules -f package.json ./node_modules

To get more help:

    hacher --help

### Running in Docker

If you previously baked hacher in a Docker image, you can run it by exposing the files to be cached as Docker volumes to the container.

```
docker run --volume $PWD:/source \
           --volume $PWD/hacher_cache:/cache  \
           --env HACHER_PATH=/cache \
           --workdir /source \
           hacher set -k node_modules -f package.json ./node_modules
```

## ToDo

Here are some things to be added later. Contributions are welcome!

* Optionally store cache in Amazon S3
* Use Golang native archive utilities
* Conditional exec ```hacher exec ...``` to run operation only if there is no cache
* Some testing would be nice :innocent:


## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).

