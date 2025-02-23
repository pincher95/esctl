package output

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

type ColumnDefaults struct {
	Header string
	Type   ColumnType
}

type sortColumn struct {
	header     string
	descending bool
}

type ColumnType int

const (
	Text ColumnType = iota
	Number
	Percent
	DataSize
	Date
)

func compareValues(left, right string, columnType ColumnType) bool {
	switch columnType {
	case Text:
		return sortText(left, right)
	case Number:
		return sortNumber(left, right)
	case DataSize:
		return sortDataSize(left, right)
	case Percent:
		return sortPercent(left, right)
	case Date:
		return sortDate(left, right)
	}
	return false
}

func PrintTable(columnDefs []ColumnDefaults, data [][]string, sortCols []sortColumn) {
	// Detect empty columns (unchanged from your snippet):
	emptyColumns := make([]bool, len(columnDefs))
	for i := range columnDefs {
		empty := true
		for _, row := range data {
			if row[i] != "" {
				empty = false
				break
			}
		}
		emptyColumns[i] = empty
	}

	// Build a map to locate column index by header
	headerIndexMap := make(map[string]int)
	for i, columnDef := range columnDefs {
		// For flexibility, match headers case-insensitively
		headerIndexMap[strings.ToLower(columnDef.Header)] = i
	}

	// If we have sort columns, perform the sort
	if len(sortCols) > 0 {
		// Validate each column
		for _, sc := range sortCols {
			if _, exists := headerIndexMap[strings.ToLower(sc.header)]; !exists {
				fmt.Fprintf(os.Stderr, "header '%s' is not a valid column\n", sc.header)
				os.Exit(1)
			}
		}

		sort.SliceStable(data, func(i, j int) bool {
			for _, sc := range sortCols {
				colIndex := headerIndexMap[strings.ToLower(sc.header)]
				left, right := data[i][colIndex], data[j][colIndex]

				if left == right {
					// If they are equal, check next sort column
					continue
				}

				// Compare ascending
				cmpAsc := compareValues(left, right, columnDefs[colIndex].Type)

				// If descending, invert
				if sc.descending {
					return !cmpAsc
				}
				return cmpAsc
			}
			// All sort columns matched => keep existing order
			return false
		})
	}

	// Prepare writer
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Write headers
	for i, columnDef := range columnDefs {
		if !emptyColumns[i] {
			fmt.Fprintf(w, "%s\t", columnDef.Header)
		}
	}
	fmt.Fprintln(w)

	// Write rows
	for _, row := range data {
		for i, cell := range row {
			if !emptyColumns[i] {
				fmt.Fprintf(w, "%s\t", cell)
			}
		}
		fmt.Fprintln(w)
	}
}

func ParseSortColumns(sortByStr string) []sortColumn {
	if sortByStr == "" {
		return nil
	}
	parts := strings.Split(sortByStr, ",")
	var columns []sortColumn
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		// Look for :desc or :asc
		var col sortColumn
		if idx := strings.Index(p, ":"); idx != -1 {
			col.header = strings.TrimSpace(p[:idx])
			order := strings.ToLower(strings.TrimSpace(p[idx+1:]))
			if order == "desc" {
				col.descending = true
			}
		} else {
			// no : found -> default ascending
			col.header = p
			col.descending = false
		}
		columns = append(columns, col)
	}
	return columns
}
