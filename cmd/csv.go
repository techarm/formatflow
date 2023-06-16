package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/iancoleman/orderedmap"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

var input, output string
var pretty bool
var includeColumns string
var excludeColumns string
var keyFormat string

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Convert csv to json",
	Long:  "Convert csv to json",
	Run: func(cmd *cobra.Command, args []string) {
		if includeColumns != "" && excludeColumns != "" {
			fmt.Println("Error: Cannot use both --include and --exclude at the same time")
			os.Exit(1)
		}

		err := convertCSVToJSON(input, output, pretty, includeColumns, excludeColumns, keyFormat)
		if err != nil {
			fmt.Println("Error: ", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(csvCmd)

	csvCmd.Flags().StringVarP(&input, "input", "i", "", "input csv file path (required)")
	csvCmd.Flags().StringVarP(&output, "output", "o", "", "output json file path (optional)")
	csvCmd.Flags().BoolVarP(&pretty, "pretty", "p", false, "pretty print the output json (optional)")
	csvCmd.Flags().StringVarP(&includeColumns, "include", "n", "", "a comma-separated list of column names to be included in the json output (optional)")
	csvCmd.Flags().StringVarP(&excludeColumns, "exclude", "e", "", "a comma-separated list of column names to be excluded in the json output (optional)")
	csvCmd.Flags().StringVarP(&keyFormat, "keyFormat", "k", "default", "output json key format: camel, lowerCamel, snake or default (optional)")
	csvCmd.MarkFlagRequired("input")
}

func convertCSVToJSON(input string, output string, pretty bool, includeColumns, excludeColumns, keyFormat string) error {
	file, err := os.OpenFile(input, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	includeMap := make(map[string]bool)
	if includeColumns != "" {
		columns := strings.Split(includeColumns, ",")
		for _, column := range columns {
			includeMap[column] = false
		}

		for _, headerName := range records[0] {
			if _, exists := includeMap[headerName]; exists {
				includeMap[headerName] = true
			}
		}

		for column, included := range includeMap {
			if !included {
				return fmt.Errorf("column name %s does not exist in csv file", column)
			}
		}
	}

	excludeMap := make(map[string]bool)
	if excludeColumns != "" {
		columns := strings.Split(excludeColumns, ",")
		for _, column := range columns {
			excludeMap[column] = false
		}

		for _, headerName := range records[0] {
			if _, exists := excludeMap[headerName]; exists {
				excludeMap[headerName] = true
			}
		}

		for column, excluded := range excludeMap {
			if !excluded {
				return fmt.Errorf("column name %s does not exist in csv file", column)
			}
		}
	}

	// var data []map[string]string
	var data []*orderedmap.OrderedMap
	header := records[0]
	for _, row := range records[1:] {
		// item := make(map[string]string)
		item := orderedmap.New()
		for i, cell := range row {
			key := header[i]
			if _, exists := excludeMap[key]; exists {
				continue
			}
			if includeColumns == "" || includeMap[key] {
				// apply keyFormat
				switch keyFormat {
				case "camel":
					key = strcase.ToCamel(key)
				case "lowerCamel":
					key = strcase.ToLowerCamel(key)
				case "nake":
					key = strcase.ToSnake(key)
				default:
					// use default key
				}
				// item[key] = strings.Trim(cell, " ")
				item.Set(key, strings.Trim(cell, " "))
			}
		}
		// data = append(data, item)
		data = append(data, item)
	}

	var jsonData []byte
	if pretty {
		jsonData, err = json.MarshalIndent(data, "", "  ")
	} else {
		jsonData, err = json.Marshal(data)
	}

	if err != nil {
		return err
	}

	if output == "" {
		fmt.Println(string(jsonData))
		fmt.Printf("\nProcessed %d records, output to the console\n", len(data))
	} else {
		outFile, err := os.Create(output)
		if err != nil {
			return err
		}
		defer outFile.Close()

		_, err = io.WriteString(outFile, string(jsonData))
		if err != nil {
			return err
		}
		fmt.Printf("Processed %d records, output to file: %s\n", len(data), output)
	}

	return nil
}
