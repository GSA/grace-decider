package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/andrewstuart/servicenow"
	"github.com/google/go-github/v32/github"
	"github.com/jszwedko/go-circleci"
)

// req is a provisioning request object
type req struct {
	circleClient *circleci.Client
	email        string
	githubClient *github.Client
	githubURL    string
	inFile       string
	ritm         *ritm
	owner        string
	prNumber     int
	pr           *github.PullRequest
	repoName     string
	snowClient   *servicenow.Client
	reqMap       map[string]interface{}
}

// ritm type for the parsed ServiceNow RITM results JSON
type ritm struct {
	Action          string `json:"action"`            // ec2
	CatalogItemName string `json:"cat_item_name"`     // "GRACE-PaaS AWS RDS Provisioning Request",
	Comments        string `json:"comments"`          // "",
	DevCount        string `json:"development_count"` // 1,
	Identifier      string `json:"identifier"`        // "test-rds",
	Name            string `json:"name"`              // "TestDB",
	Number          string `json:"number"`            // "RITM0001001",
	OpenedBy        string `json:"opened_by"`         // "by@email.com",
	ProdCount       string `json:"production_count"`  // 1,
	RequestedFor    string `json:"requested_for"`     // "for@email.com",
	Supervisor      string `json:"supervisor"`        // "supervisor@email.com",
	SysID           string `json:"sys_id"`            // "99aa00000aa9aa00a9a99999a99aaa99",
	TestCount       string `json:"test_count"`        // 1,
}

func newReq() (*req, error) {
	var r req
	flags, output, err := r.parseFlags(os.Args[0], os.Args[1:])
	if err != nil {
		fmt.Println(output)
		return &r, err
	}

	err = r.check()
	if err != nil {
		flags.PrintDefaults()
		return &r, err
	}

	err = r.parseRITM()
	if err != nil {
		return &r, err
	}

	r.email = "grace-staff@gsa.gov"
	r.githubURL = "https://github.com/GSA/"
	r.circleClient = newCircleClient(os.Getenv("CIRCLE_TOKEN"))
	r.githubClient = newAuthenticatedClient()
	r.snowClient = newSnowClient()
	r.pr, err = r.getPR()

	return &r, err
}

func (r *req) parseFlags(progName string, args []string) (*flag.FlagSet, string, error) {
	flags := flag.NewFlagSet(progName, flag.ContinueOnError)
	var buf bytes.Buffer
	flags.SetOutput(&buf)
	flags.StringVar(&r.inFile, "request", "", "JSON input file")
	flags.StringVar(&r.owner, "owner", "GSA", "GitHub repository owner")
	flags.StringVar(&r.repoName, "repo", "", "GitHub repository name")
	flags.IntVar(&r.prNumber, "pr", 0, "Pull Request Number")
	err := flags.Parse(args)
	if err != nil {
		return flags, buf.String(), err
	}
	return flags, buf.String(), nil
}

func (r *req) checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		if r.ritm != nil {
			err := r.updateRITM(err)
			if err != nil {
				fmt.Println(err)
			}
		}
		os.Exit(1)
	}
}

// HandlePR ... Tracks status of GRACE PaaS IaC Provisioning request Pull Request (PR)
func HandlePR(opt ...*req) {
	var r *req
	var err error
	if len(opt) > 0 {
		r = opt[0]
	} else {
		r, err = newReq()
		r.checkErr(err)
	}
	r.trackStatus()
}

func (r *req) trackStatus() {
	pr, err := r.getPR()
	r.checkErr(err)

	err = waitForMerge(pr)
	r.checkErr(err)

	err = waitForApply(pr)
	r.checkErr(err)

	err = r.updateRITM(nil)
	r.checkErr(err)

	fmt.Println("Processing complete")
}

func (r *req) parseRITM() error {
	fmt.Printf("Parsing RITM from: %s\n", r.inFile)
	jsonFile, err := os.Open(r.inFile) // #nosec G304
	if err != nil {
		return err
	}

	defer jsonFile.Close() // #nosec G307
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	var ritm ritm
	err = json.Unmarshal(byteValue, &ritm)
	if err != nil {
		return err
	}
	r.ritm = &ritm

	var myMap map[string]interface{}
	err = json.Unmarshal(byteValue, &myMap)
	if err != nil {
		return err
	}
	r.reqMap = myMap

	return nil
}

func (r *req) check() error {
	if r.inFile == "" {
		return fmt.Errorf("request must be set")
	}

	if r.owner == "" {
		return fmt.Errorf("owner must be set")
	}

	if r.repoName == "" {
		return fmt.Errorf("repo must be set")
	}

	if r.prNumber == 0 {
		return fmt.Errorf("pr number must be set")
	}

	if os.Getenv("GITHUB_TOKEN") == "" {
		return fmt.Errorf("environment variable GITHUB_TOKEN must be set")
	}

	if os.Getenv("CIRCLE_TOKEN") == "" {
		return fmt.Errorf("environment variable CIRCLE_TOKEN must be set")
	}

	if os.Getenv("SN_INSTANCE") == "" {
		return fmt.Errorf("environment variable SN_INSTANCE must be set")
	}

	if os.Getenv("SN_PASSWORD") == "" {
		return fmt.Errorf("environment variable SN_PASSWORD must be set")
	}

	if os.Getenv("SN_USER") == "" {
		return fmt.Errorf("environment variable SN_USER must be set")
	}
	return nil
}
