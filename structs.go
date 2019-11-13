package main

import r "gopkg.in/rethinkdb/rethinkdb-go.v5"

type Token struct {
	Id  string `gorethink:"id,omitempty"`
	Uid string `gorethink:"uid"`
}

type S3Bucket struct {
	Endpoint        string `gorethink:"endpoint"`
	Bucket          string `gorethink:"bucket"`
	AccessKeyId     string `gorethink:"access_key_id"`
	SecretAccessKey string `gorethink:"secret_access_key"`
	Region          string `gorethink:"region"`
}

type Domain struct {
	Id        string    `gorethink:"id,omitempty"`
	Public    bool      `gorethink:"public"`
	Whitelist []string  `gorethink:"whitelist"`
	Blacklist []string  `gorethink:"blacklist"`
	Owner     string    `gorethink:"owner"`
	Bucket    *S3Bucket `gorethink:"bucket"`
}

type User struct {
	Id           string   `gorethink:"id,omitempty"`
	Domain       string   `gorethink:"domain"`
	Tokens       []string `gorethink:"tokens"`
	NamingScheme string   `gorethink:"naming_scheme"`
	Encryption   bool     `gorethink:"encryption"`
}

func (user User) NewDomain(DomainName string) {
	DomainInfo := Domain{
		Id:        DomainName,
		Public:    false,
		Whitelist: []string{user.Id},
		Blacklist: []string{},
		Owner:     user.Id,
	}
	err := r.Table("domains").Insert(DomainInfo).Exec(RethinkConnection)
	if err != nil {
		panic(err)
	}
}
