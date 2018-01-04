package db

func BulkLoadNodes(src *DataSrc, vals [][]string, columns string ) error {
	_, err := src.Exec(SQLInsertIntoTable(src.Table, columns, vals))
	return err
}
