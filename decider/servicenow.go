package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/andrewstuart/servicenow"
)

func newSnowClient() *servicenow.Client {
	return &servicenow.Client{
		Username: os.Getenv("SN_USER"),
		Password: os.Getenv("SN_PASSWORD"),
		Instance: os.Getenv("SN_INSTANCE"),
	}
}

func (r *req) updateRITM(e error) error {
	fmt.Printf("Updating %s (%s)\n", r.ritm.Number, r.ritm.SysID)
	table := "sc_req_item"
	var out map[string]interface{}
	var state = 2 // Work in Progress
	var comment = strings.ToUpper(r.ritm.Action) + " provisioned via GRACE-PaaS CI/CD Pipeline"
	if e != nil {
		state = 8 // Reopened
		comment = fmt.Sprintf("Error provisioning %s: %v", strings.ToUpper(r.ritm.Action), e)
	}

	body := map[string]interface{}{
		"state":    state,
		"comments": comment,
	}

	fmt.Printf("Body: %v:\n", body)
	err := r.snowClient.PerformFor(table, "update", r.ritm.SysID, nil, body, &out)
	fmt.Printf("Out: %v\n", out)
	return err
}
