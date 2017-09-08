## Releasing ZBC-GO

ZBC-GO uses [Goreleaser](https://github.com/goreleaser/goreleaser) for doing a release.
The configuration can be found in the [.goreleaser.yml](https://github.com/zeebe-io/zbc-go/blob/master/.goreleaser.yml) file.

### Travis-CI
Doing a release with Travis-CI only requires creating a new git tag:

```
git tag -a v0.1.0 -m "Awesome zbctl release" && git push origin v0.1.0
```

### Locally
To test the release process locally, download [goreleaser](https://github.com/goreleaser/goreleaser/releases) and run in the root of the project

```
goreleaser --skip-validate --skip-publish --snapshot
```

