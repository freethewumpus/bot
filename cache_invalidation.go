package main

import "fmt"

func InvalidateDomainCache(domain string) {
	RedisConnection.Del(fmt.Sprintf("d:%s", domain))
}

func InvalidateUserCache(tokens []string) {
	for _, k := range tokens {
		RedisConnection.Del(fmt.Sprintf("u:%s", k))
	}
}
