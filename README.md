# Scrapenstein

This project contains an assortment of API scrapers for various services that I use personally and professionally. The results of the scrapes are stored in Postgres, where tools such as [Metabase](https://www.metabase.com/), [Google Data Studio](https://datastudio.google.com/), and others can be used to perform analysis.

Though I've structured and built this project around my own needs, please feel free to use and contribute back if you find it to be useful!

## Status

There may be large, breaking changes at any point in time.

## Prerequisites

* A Postgres DB. 10.x or higher is recommended.
* Go 1.14+ or higher.

## Install

In absence of tagged releases, the best bet for installing Scrapenstein is `go get`:
```shell script
go get -u github.com/gtaylor/scrapenstein
```

You can then run the `scrapeinstein` command.

## Usage

The included CLI is designed with exploration in mind. Run the `scrapeinstein` command and take a look at the sub-commands within.

## Support

These sources are offered as-is with no official support provided. 

## License

Scrapenstein is licensed under the [MIT License](./LICENSE).
