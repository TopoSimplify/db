package db

func BulkLoadNodes(src *DataSrc, nodes []*Node) error {
	var vals = make([][]string, 0)
	for _, h := range nodes {
		vals = append(vals, h.ColumnValues(src.SRID))
	}
	_, err := src.Exec(SQLInsertIntoNodeTable(src.NodeTable, NodeTblColumns, vals))
	return err
}
