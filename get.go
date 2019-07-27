package main

import r "gopkg.in/rethinkdb/rethinkdb-go.v5"

func GetUser(UserId string) User {
	cursor, err := r.Table("users").Get(UserId).Run(RethinkConnection)
	if err != nil {
		panic(err)
	}
	if cursor.IsNil() {
		user := User{
			Id: UserId,
			Domain: "freethewump.us",
			NamingScheme: "ccccc",
			Tokens: []string{},
		}
		err = r.Table("users").Insert(&user).Exec(RethinkConnection)
		if err != nil {
			panic(err)
		}
		return user
	} else {
		var user User
		err = cursor.One(&user)
		if err != nil {
			panic(err)
		}
		return user
	}
}

func GetDomain(DomainName string) *Domain {
	cursor, err := r.Table("domains").Get(DomainName).Run(RethinkConnection)
	if err != nil {
		panic(err)
	}
	if cursor.IsNil() {
		return nil
	}
	var domain Domain
	err = cursor.One(&domain)
	if err != nil {
		panic(err)
	}
	return &domain
}

func (u User) GetOwnedDomains() []Domain {
	var OwnedDomains []Domain
	cursor, err := r.Table("domains").GetAllByIndex("owner", u.Id).Run(RethinkConnection)
	if err != nil {
		panic(err)
	}
	err = cursor.All(&OwnedDomains)
	if err != nil {
		panic(err)
	}
	return OwnedDomains
}

func (u User) GetWhitelistedDomains() []Domain {
	var OwnedDomains []Domain
	cursor, err := r.Table("domains").GetAllByIndex("whitelist", u.Id).Run(RethinkConnection)
	if err != nil {
		panic(err)
	}
	err = cursor.All(&OwnedDomains)
	if err != nil {
		panic(err)
	}
	var ParsedDomain []Domain
	for _, v := range OwnedDomains {
		if v.Owner != u.Id && !v.Public {
			ParsedDomain = append(ParsedDomain, v)
		}
	}
	return ParsedDomain
}
