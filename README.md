# MiniScrape

Simple webpages scrapper written in GO

## Installation

Clone the repository

```shell
git clone https://github.com/pestanko/miniscrape.git
```

Enter the clonned repository

```shell
cd miniscrape
```

Install the dependencies

```shell
go get .
```

## Build the scraper

```shell
make build
```

## Run the scraper

```shell
go run main.go scrape
```

Scrape the single webpage:

```shell
# For food category
go run main.go scrape -C food -N ubaumanu
```

### Run the server

```shell
make run-serve
```

## Add/Edit available webpages

The webpages list is located in ``./config/default.yml``.

## License

Miniscrape is released under the Apache 2.0 license. See LICENSE
