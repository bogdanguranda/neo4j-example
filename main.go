package main

import (
	"github.com/jmcvetta/neoism"
	"github.com/Sirupsen/logrus"
)

func main() {
	db, _ := neoism.Connect("http://neo4j:root@localhost:7474/db/data")
	initializeDB(db)

	res := []struct {
		A   string `json:"a.name"`
		Rel string `json:"type(r)"`
		B   string `json:"b.name"`
	}{}

	cq := neoism.CypherQuery{
		Statement: `
			MATCH (a:Person)-[r]->(b)
			WHERE a.name = {name}
			RETURN a.name, type(r), b.name
    			`,
		Parameters: neoism.Props{"name": "Napoleon Bonaparte"},
		Result:     &res,
	}

	db.Cypher(&cq)

	for _, element := range res {
		logrus.Infoln(element.A)
		logrus.Infoln(element.Rel)
		logrus.Infoln(element.B)
	}
}

func initializeDB(db *neoism.Database) {
	cqCleanDB := neoism.CypherQuery{
		Statement: `
			MATCH (n) DETACH
			DELETE n
    			`,
	}
	db.Cypher(&cqCleanDB)

	nodePerson1, _ := db.CreateNode(neoism.Props{"name": "Captain Kirk"})
	nodePerson1.AddLabel("Person")

	nodePerson2, _ := db.CreateNode(neoism.Props{"name": "Napoleon Bonaparte"})
	nodePerson2.AddLabel("Person")

	nodeCountry1, _ := db.CreateNode(neoism.Props{"name": "France"})
	nodeCountry1.AddLabel("Country")

	nodeCountry2, _ := db.CreateNode(neoism.Props{"name": "USA"})
	nodeCountry2.AddLabel("Country")

	allegiance1, _ := db.CreateNode(neoism.Props{"name": "The Frence Empire"})
	allegiance1.AddLabel("Allegience")

	nodePerson1.Relate("lives in", nodeCountry2.Id(), neoism.Props{})
	nodePerson2.Relate("lives in", nodeCountry1.Id(), neoism.Props{})
	nodePerson2.Relate("serves", allegiance1.Id(), neoism.Props{})
}
