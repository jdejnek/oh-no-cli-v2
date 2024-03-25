package main

import "github.com/jdejnek/oh-no-cui/http_client"

type simPage struct {
	sim http_client.Sim
}

func simPageModel() simPage {
	return simPage{}
}
