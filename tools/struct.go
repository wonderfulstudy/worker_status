package tools

var URL string = "https://dev-squid.cybros.network/graphql"
var RESPONSE WorkersResponse

type WorkersResponse struct {
	Data struct {
		Workers []struct {
			Address      string `json:"address"`
			OwnerAddress string `json:"ownerAddress"`
			Status       string `json:"status"`
			Uptime       int    `json:"uptime"`
		} `json:"workers"`
	} `json:"data"`
}

var AddressQuery string = `
	query MyQuery {
		workers(where: {ownerAddress_eq: "%s"}) {
			address
			ownerAddress
			status
			uptime
		}
	}`

var AddressesQuery string = `
	query MyQuery {
		workers(where: {ownerAddress_in: ["%v"]}) {
			address
			ownerAddress
			status
			uptime
		}
	}`
