package neo4jx

import (
	"encoding/json"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

const TagName = "json"

func UnmarshalRecord(dst any, record *neo4j.Record, outName string) error {
	node, ok := record.AsMap()[outName].(dbtype.Node)
	if !ok {
		return fmt.Errorf("outName is not a dbtype.Node")
	}

	b, e := json.Marshal(node.Props)
	if e != nil {
		return e
	}
	return json.Unmarshal(b, dst)
}
