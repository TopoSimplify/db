package db

func BulkLoadNodes(src *DataSrc, nodes []*Node, columns string ) error {
	var vals = make([][]string, 0)
	for _, h := range nodes {
		vals = append(vals, h.ColumnValues(src.SRID))
	}
	_, err := src.Exec(SQLInsertIntoNodeTable(src.Table, columns, vals))
	return err
}
