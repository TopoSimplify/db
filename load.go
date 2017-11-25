package db

import "fmt"

func BulkLoadNodes(src *DataSrc, nodes []*Node) error {
	var vals = make([][]string, 0)
	for _, h := range nodes {
		vals = append(vals, []string{
			fmt.Sprintf(`%v`, h.FID),
			fmt.Sprintf(`'%v'`, Serialize(h)),
			fmt.Sprintf(`ST_GeomFromText('%v', %v)`, h.WTK, src.SRID),
		})
	}
	_, err := src.Exec(SQLInsertIntoTable(src.NodeTable, nodeTblColumns, vals))
	return err
}

