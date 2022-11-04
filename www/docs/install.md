# Install

Generally, ECS Deployer should be run from your Continuous Integration/Continuous Deployment system (like [GitHub Actions](ci/github.md)).

However, if you would like to install the app you have the options listed below.


## Install on CI/CD Services

!!! note "Using ECS Deployer on CI/CD is the preferred method!"

* [GitHub Actions](ci/github.md)

----

## Install the pre-compiled binary

### homebrew tap

```sh
brew install ecsdeployer/tap/ecsdeployer
```

### scoop

```sh
scoop bucket add ecsdeployer https://github.com/ecsdeployer/scoop-bucket.git
scoop install ecsdeployer
```

### deb, rpm and apk packages

Download the `.deb`, `.rpm` or `.apk` packages from the [latest release on GitHub](https://github.com/ecsdeployer/ecsdeployer/releases/latest) and install them with the appropriate tools.

### go install

```sh
go install ecsdeployer.com/ecsdeployer@latest
```

### bash script
```sh
curl -sfL https://ecsdeployer.com/run.sh | bash -s -- deploy --config CONFIGFILE
```

!!! note ""
    This is not installing ECS Deployer. This will only _run_ the app.


## Compiling from source

If you just want to build from source for whatever reason, follow these steps:

**clone:**

```sh
git clone https://github.com/ecsdeployer/ecsdeployer
cd ecsdeployer
```

**get the dependencies:**

```sh
go mod tidy
```

**build:**

```sh
go build -o ecsdeployer .
```

**verify it works:**

```sh
./ecsdeployer --version
```
