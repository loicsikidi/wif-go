[![](https://github.com/loicsikidi/wif-go/actions/workflows/tests.yml/badge.svg)](https://github.com/loicsikidi/wif-go/actions/workflows/tests.yml)
[![](https://codecov.io/gh/loicsikidi/wif-go/branch/main/graph/badge.svg?token=G6DVJ1GUYU)](https://codecov.io/gh/loicsikidi/wif-go)
[![](https://img.shields.io/github/release/loicsikidi/wif-go.svg?label=version)](https://github.com/loicsikidi/wif-go/releases/latest)
[![](https://img.shields.io/github/go-mod/go-version/loicsikidi/wif-go.svg?label=go)](https://github.com/loicsikidi/wif-go)
[![](https://goreportcard.com/badge/github.com/loicsikidi/wif-go)](https://goreportcard.com/report/github.com/loicsikidi/wif-go)
[![](https://img.shields.io/badge/License-Apache%202.0-blue.svg?label=license)](https://github.com/loicsikidi/wif-go/blob/main/LICENSE)

# wif-go (Workload Identity Federation)

Tool (implemented in Golang) emulating the behavior of [Workload Identity Federation](https://cloud.google.com/iam/docs/workload-identity-federation).

Features 🚀:

* Playground in order to test interactively if a _subject token_ match or not a WIF setup. A public instance is available [here](https://play.wif.lsikidi.org)!
* `wif-go`: Package (used by the playground) emulating WIF behavior when a _subject token_ is given

## Why

Today, GCP _(Google Cloud Platforms)_ doesn't provide a way to test `Workload Identity Federation` setup beforehand (eg. unit test, web playground) in order to check if the _attribute mapping_ and/or the _attibute condition_ is suitable for your use case.

## Roadmap

Provider support:

  * [x] `oidc`
  * [] `aws`
  * [] `saml`

Optimization:

  * [] `wif-go.wasm`: Improve the size (currently ~ 16MB) in order to load the playground faster

## Acknowledgement 🫶

* The WIF Playground borrows a lot of ideas and styles from [Rego Playground](https://play.openpolicyagent.org/).
* Logo used in the playground has been generated at [Gopherize.me](https://gopherize.me/).

## Disclaimer

This is a personal project, while I do my best to ensure that everything works, I take no responsibility for issues caused by this code.