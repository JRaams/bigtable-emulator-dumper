package bigtable

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"cloud.google.com/go/bigtable"
	"github.com/jraams/bigtable-emulator-dumper/config"
)

func New(ctx context.Context, cfg *config.Config) (*Service, error) {
	adminClient, err := bigtable.NewAdminClient(ctx, cfg.BigtableProject, cfg.BigtableInstance)
	if err != nil {
		return nil, errors.New("failed to create Bigtable admin client")
	}

	_, err = adminClient.Tables(ctx)
	if err != nil {
		return nil, errors.New("failed to connect with Bigtable admin client")
	}

	client, err := bigtable.NewClient(ctx, cfg.BigtableProject, cfg.BigtableInstance)
	if err != nil {
		return nil, errors.New("failed to create Bigtable client")
	}

	return &Service{
		cfg:         cfg,
		adminClient: adminClient,
		client:      client,
	}, nil
}

func (s *Service) Close() {
	if s.adminClient != nil {
		s.adminClient.Close()
	}

	if s.client != nil {
		s.client.Close()
	}
}

func (s *Service) FetchAllTables(ctx context.Context) (*[]TableData, error) {
	tableNames, err := s.adminClient.Tables(ctx)
	if err != nil {
		return nil, errors.New("failed to list tables")
	}

	data := []TableData{}
	for _, tableName := range tableNames {
		rows, err := s.readTableData(ctx, tableName)
		if err != nil {
			return nil, err
		}

		data = append(data, TableData{
			TableName: tableName,
			Rows:      *rows,
		})
	}

	return &data, nil
}

func (s *Service) FetchSingleTable(ctx context.Context, tableName string) (*[]RowData, error) {
	return s.readTableData(ctx, tableName)

}

func (s *Service) readTableData(ctx context.Context, tableName string) (*[]RowData, error) {
	table := s.client.Open(tableName)
	rows := []RowData{}

	err := table.ReadRows(ctx, bigtable.PrefixRange(""), func(row bigtable.Row) bool {
		rowData := RowData{
			RowKey:   row.Key(),
			Families: make(map[string]FamilyData),
		}
		for family, items := range row {
			if _, exists := rowData.Families[family]; !exists {
				rowData.Families[family] = make(FamilyData)
			}
			for _, item := range items {
				column := item.Column[strings.Index(item.Column, ":")+1:]

				var j any
				err := json.Unmarshal(item.Value, &j)
				if err != nil {
					rowData.Families[family][column] = string(item.Value)
				} else {
					rowData.Families[family][column] = j
				}
			}
		}
		rows = append(rows, rowData)
		return true
	})
	if err != nil {
		return nil, errors.New("failed to read rows from table")
	}

	return &rows, nil
}
