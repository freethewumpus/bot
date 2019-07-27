package main

import r "gopkg.in/rethinkdb/rethinkdb-go.v5"

func GetPublicDomains() []Domain {
	cursor, err := r.Table("domains").GetAllByIndex("public", true).Run(RethinkConnection)
	if err != nil {
		panic(err)
	}
	var domains []Domain
	err = cursor.All(&domains)
	if err != nil {
		panic(err)
	}
	return domains
}
