package db

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func WrapError(err error, message string) error {
	if e, ok := err.(*pq.Error); ok {
		var sb strings.Builder

		fmt.Fprintf(&sb, "%v\n", err.Error())

		if len(e.Severity) > 0 {
			fmt.Fprintf(&sb, "\tSeverity: %v\n", e.Severity)
		}
		if len(e.Code) > 0 {
			fmt.Fprintf(&sb, "\tCode: %v\n", e.Code)
		}
		if len(e.Message) > 0 {
			fmt.Fprintf(&sb, "\tMessage: %v\n", e.Message)
		}
		if len(e.Detail) > 0 {
			fmt.Fprintf(&sb, "\tDetail: %v\n", e.Detail)
		}
		if len(e.Hint) > 0 {
			fmt.Fprintf(&sb, "\tHint: %v\n", e.Hint)
		}
		if len(e.Position) > 0 {
			fmt.Fprintf(&sb, "\tPosition: %v\n", e.Position)
		}
		if len(e.InternalPosition) > 0 {
			fmt.Fprintf(&sb, "\tInternalPosition: %v\n", e.InternalPosition)
		}
		if len(e.Where) > 0 {
			fmt.Fprintf(&sb, "\tWhere: %v\n", e.Where)
		}
		if len(e.Schema) > 0 {
			fmt.Fprintf(&sb, "\tSchema: %v\n", e.Schema)
		}
		if len(e.Table) > 0 {
			fmt.Fprintf(&sb, "\tTable: %v\n", e.Table)
		}
		if len(e.Column) > 0 {
			fmt.Fprintf(&sb, "\tColumn: %v\n", e.Column)
		}
		if len(e.DataTypeName) > 0 {
			fmt.Fprintf(&sb, "\tDataTypeName: %v\n", e.DataTypeName)
		}
		if len(e.Constraint) > 0 {
			fmt.Fprintf(&sb, "\tConstraint: %v\n", e.Constraint)
		}
		if len(e.File) > 0 {
			fmt.Fprintf(&sb, "\tFile: %v\n", e.File)
		}
		if len(e.Line) > 0 {
			fmt.Fprintf(&sb, "\tLine: %v\n", e.Line)
		}
		if len(e.Routine) > 0 {
			fmt.Fprintf(&sb, "\tRoutine: %v\n", e.Routine)
		}

		err = errors.Wrap(errors.New(sb.String()), message)
	}

	return err
}
