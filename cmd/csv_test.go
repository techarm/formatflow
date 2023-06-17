package cmd

import (
	"os"
	"testing"
)

func TestConvertCSVToJSON(t *testing.T) {
	tests := []struct {
		name           string
		csvData        string
		includeColumns string
		excludeColumns string
		keyFormat      string
		expectedJSON   string
	}{
		{
			name: "Basic CSV to JSON conversion",
			csvData: `Name,Age,Job
Alice,30,Developer
Bob,25,Designer`,
			includeColumns: "",
			excludeColumns: "",
			keyFormat:      "default",
			expectedJSON:   `[{"Name":"Alice","Age":"30","Job":"Developer"},{"Name":"Bob","Age":"25","Job":"Designer"}]`,
		},
		{
			name: "Exclude column",
			csvData: `Name,Age,Job
Alice,30,Developer
Bob,25,Designer`,
			includeColumns: "",
			excludeColumns: "Job",
			keyFormat:      "default",
			expectedJSON:   `[{"Name":"Alice","Age":"30"},{"Name":"Bob","Age":"25"}]`,
		},
		{
			name: "Include column",
			csvData: `Name,Age,Job
Alice,30,Developer
Bob,25,Designer`,
			includeColumns: "Name,Age",
			excludeColumns: "",
			keyFormat:      "default",
			expectedJSON:   `[{"Name":"Alice","Age":"30"},{"Name":"Bob","Age":"25"}]`,
		},
		{
			name: "Key Format Camel Case",
			csvData: `first_name,age,job_title
Alice,30,Developer
Bob,25,Designer`,
			includeColumns: "",
			excludeColumns: "",
			keyFormat:      "camel",
			expectedJSON:   `[{"FirstName":"Alice","Age":"30","JobTitle":"Developer"},{"FirstName":"Bob","Age":"25","JobTitle":"Designer"}]`,
		},
		{
			name: "Key Format Lowercase Camel Case",
			csvData: `first_name,age,job_title
Alice,30,Developer
Bob,25,Designer`,
			includeColumns: "",
			excludeColumns: "",
			keyFormat:      "lowerCamel",
			expectedJSON:   `[{"firstName":"Alice","age":"30","jobTitle":"Developer"},{"firstName":"Bob","age":"25","jobTitle":"Designer"}]`,
		},
		{
			name: "Key Format Snake Case",
			csvData: `FirstName,Age,JobTitle
Alice,30,Developer
Bob,25,Designer`,
			includeColumns: "",
			excludeColumns: "",
			keyFormat:      "snake",
			expectedJSON:   `[{"first_name":"Alice","age":"30","job_title":"Developer"},{"first_name":"Bob","age":"25","job_title":"Designer"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary CSV file
			tmpCSVFile, err := os.CreateTemp("", "test_*.csv")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpCSVFile.Name())

			// Write data into the temporary CSV file
			if _, err := tmpCSVFile.Write([]byte(tt.csvData)); err != nil {
				t.Fatal(err)
			}
			tmpCSVFile.Close()

			// Create a temporary JSON file
			tmpJSONFile, err := os.CreateTemp("", "test_*.json")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tmpJSONFile.Name())
			tmpJSONFile.Close()

			// Test the convertCSVToJSON function
			err = convertCSVToJSON(tmpCSVFile.Name(), tmpJSONFile.Name(), false, tt.includeColumns, tt.excludeColumns, tt.keyFormat)
			if err != nil {
				t.Errorf("convertCSVToJSON failed: %v", err)
			}

			// Read the generated JSON file and validate its content
			jsonData, err := os.ReadFile(tmpJSONFile.Name())
			if err != nil {
				t.Fatal(err)
			}

			if string(jsonData) != tt.expectedJSON {
				t.Errorf("Expected JSON: %s, got: %s", tt.expectedJSON, string(jsonData))
			}
		})
	}
}
