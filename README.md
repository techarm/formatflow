# FormatFlow

FormatFlow is an open-source project for converting different data formats, especially from CSV to JSON. This project is written in Go and uses the cobra library to build CLI applications.

## Features

1. Conversion from CSV to JSON
2. Selectively include or exclude certain columns of data
3. Option to output JSON in a pretty-printed manner
4. Option to output JSON keys in camel, lowerCamel, snake or default naming
5. The order of keys in the generated JSON matches the order of columns in the CSV

## Installation

You can use Go's package management tool `go get` to install:

```
go install github.com/techarm/formatflow
```

## Usage

The basic usage is as follows:

```
formatflow csv -i input.csv -o output.json
```

Where `-i` or `--input` argument specifies the input CSV file name, `-o` or `--output` argument specifies the output JSON file name. If `-o` is not specified, the output will be printed to the console.

Optional arguments include:

- `-p` or `--pretty`: pretty-print JSON
- `--include`: only include specified columns, separated by commas
- `--exclude`: exclude specified columns, separated by commas
- `-k` or `--keyFormat`: choose the naming style of output JSON keys, optional values are `camel`, `lowerCamel`, `snake` or `default`

Example:

```
formatflow csv -i input.csv -o output.json --include column1,column2 --pretty --keyFormat lowerCamel
```

This command will convert the data in `column1` and `column2` of the `input.csv` file into JSON format with camelCase naming, and then write it to `output.json` file in a pretty-printed manner.

## Note

The `--include` and `--exclude` arguments cannot be used at the same time, if used together, an error will be returned.

## Contributing

Contributions of any kind are welcome, including but not limited to submitting issues, suggesting improvements, improving code, etc.

## License

This project is licensed under the MIT license, see the LICENSE file for details.
