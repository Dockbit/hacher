# Hacher :hocho:

```hacher``` is a simple CLI tool to cache project artifacts. It tracks the changes to a dependency file, and if changed caches related artifacts directory. For example, you could cache ```bundle``` directory on [Gemfile](http://bundler.io/man/gemfile.5.html) changes or a ```node_modules``` when [package.json](https://docs.npmjs.com/files/package.json) changes.

## Installation

Install [Golang](https://golang.org/doc/install) and:

    make install

```hacher``` needs several environment variables to operate:

* ```HACHER_PATH``` - The path where cached artifacts will be stored.
* ```HACHER_KEEP``` - The number of caches to keep, defaults to ```3```. You probably won't need more, since it's only useful if someone reverts the dependency file, hence cache could be reused.

## Usage

Let's say you wanted to speed up ```npm install``` during your CI builds, so you could:

    # Get the cached version of node_modules to the current directory
    hacher get -k node_modules -f package.json

    npm install

    # Cache node_modules on package.json changes
    hacher set -k node_modules -f package.json ./node_modules

To get more help:

    hacher --help

## ToDo

Here are some things to be added later. Contributions are welcome!

* Optionally store cache in Amazon S3
* Use Golang native archive utilities
* Conditional exec ```hacher exec ...``` to run operation only if there is no cache
* Some testing would be nice :innocent:


## License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).

