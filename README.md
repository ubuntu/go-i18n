# Welcome to go-i18n

[actions-image]: https://github.com/ubuntu/go-i18n/actions/workflows/qa.yaml/badge.svg?branch=main
[actions-url]: https://github.com/ubuntu/go-i18n/actions?query=branch%3Amain+event%3Apush

[license-image]: https://img.shields.io/badge/License-MIT-blue.svg

[codecov-image]: https://codecov.io/gh/ubuntu/go-i18n/branch/master/graph/badge.svg
[codecov-url]: https://codecov.io/gh/ubuntu/go-i18n

[reference-documentation-image]: https://pkg.go.dev/badge/github.com/ubuntu/go-i18n.svg
[reference-documentation-url]: https://pkg.go.dev/github.com/ubuntu/go-i18n

[goreport-image]: https://goreportcard.com/badge/github.com/ubuntu/go-i18n
[goreport-url]: https://goreportcard.com/report/github.com/ubuntu/go-i18n

[![Code quality][actions-image]][actions-url]
[![License][license-image]](LICENSE)
[![Code coverage][codecov-image]][codecov-url]
[![Reference documentation][reference-documentation-image]][reference-documentation-url]
[![Go Report Card][goreport-image]][goreport-url]

This is the code repository for **go-i18n**, a go-text wrapper joining gettext support for Linux and Windows.

This package allows to transparently embeds local translation or lookup in system path for installed translations on both platforms. It also includes a composite Github action to update translation on any folder. It will initialize and loads translation, ready to be used by [gotext](https://github.com/leonelquinteros/gotext).

For usage in your own project, please refer to the [reference documentation]([reference-documentation-url]).

## Reusable github action

A reusable action will extract any translatable strings using [gotext](https://github.com/leonelquinteros/gotext) functions in your code. Those will generate an up to date `<domain>.pot` file in the destination directory. Any `<locale>.po` file inside this directory will then be refreshed with the new available translations.

Usage example:

```yaml

```

To bootstrap a new locale, you can `cp <domain>.pot <locale>.po` and commit it.

## Troubleshooting

The project is using the slog package from Go 1.21. You can increase the verbosity of your embedding code to have more logs printed.

## Get involved

This is an [open source](LICENSE) project and we warmly welcome community contributions, suggestions, and constructive feedback. If you're interested in contributing, please take a look at our [Contribution guidelines](CONTRIBUTING.md) first.

- to report an issue, please file a bug report against our repository, using a bug template.
- for suggestions and constructive feedback, report a feature request bug report, using the proposed template.

## Get in touch

We're friendly! We have a community forum at [https://discourse.ubuntu.com](https://discourse.ubuntu.com) where we discuss feature plans, development news, issues, updates and troubleshooting.

For news and updates, follow the [Ubuntu twitter account](https://twitter.com/ubuntu) and on [Facebook](https://www.facebook.com/ubuntu).
