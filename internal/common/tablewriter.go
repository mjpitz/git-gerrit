package common

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/mjpitz/myago"
	"github.com/olekukonko/tablewriter"
)

const tableWriterContextKey = myago.ContextKey("table_writer")

// TableWriter obtains a tablewriter.Table from the provided context.
func TableWriter(ctx context.Context) *tablewriter.Table {
	v := ctx.Value(tableWriterContextKey)
	if v == nil {
		return nil
	}
	return v.(*tablewriter.Table)
}

// SetupTableWriter creates a new table writer and appends it to the context.
func SetupTableWriter(ctx context.Context, out io.Writer) context.Context {
	table := tablewriter.NewWriter(out)

	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	// do this in callers
	// table.SetHeader([]string{"Name", "Status", "Role", "Version"})
	// table.Append(row)
	// table.AppendBulk(rows)
	return context.WithValue(ctx, tableWriterContextKey, table)
}

func extract(v interface{}) (headers, row []string) {
	value := reflect.Indirect(reflect.ValueOf(v))

	if value.Kind() != reflect.Struct {
		return
	}

	headers = make([]string, 0, value.NumField())
	row = make([]string, 0, value.NumField())

	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i)
		fieldType := value.Type().Field(i)

		var cell string
		switch c := fieldValue.Interface().(type) {
		case fmt.Stringer:
			cell = c.String()
		default:
			switch fieldValue.Kind() {
			case reflect.String:
				cell = fieldValue.String()
			case reflect.Bool:
				cell = strconv.FormatBool(fieldValue.Bool())
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				cell = strconv.FormatInt(fieldValue.Int(), 10)
			case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				cell = strconv.FormatUint(fieldValue.Uint(), 10)
			}
		}

		header := fieldType.Name
		parts := strings.Split(fieldType.Tag.Get("json"), ",")
		if parts[0] != "" {
			header = parts[0]
		}

		headers = append(headers, header)
		row = append(row, cell)
	}

	return headers, row
}

// WriteRow using golang reflection to get column names from the provided structure.
func WriteRow(table *tablewriter.Table, v interface{}) {
	headers, row := extract(v)
	if headers == nil || row == nil {
		return
	}

	if table.NumLines() == 0 {
		table.SetHeader(headers)
	}

	table.Append(row)
}
