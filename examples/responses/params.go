package main

type FindPetsQuery struct {
	Tags  *[]string `query:"tags" description:"tags to filter by"`
	Limit *int32    `query:"limit" description:"maximum number of results to return"`
}
