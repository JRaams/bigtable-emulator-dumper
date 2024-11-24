package bigtable

import (
	"cloud.google.com/go/bigtable"
	"github.com/jraams/bigtable-emulator-dumper/config"
)

type Service struct {
	cfg         *config.Config
	adminClient *bigtable.AdminClient
	client      *bigtable.Client
}

type TableData struct {
	TableName string    `json:"table_name"`
	Rows      []RowData `json:"rows"`
}

type RowData struct {
	RowKey   string                `json:"row_key"`
	Families map[string]FamilyData `json:"families"`
}

type FamilyData map[string]any
