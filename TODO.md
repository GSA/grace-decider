# To Do #

1. Refactor the go code into a handler module and a cmd and lambda module that
wraps the handler.
1. Refactor the build process to build the cmd and lambda binaries
1. Change `servicenow.go` to update the ticket's `automation_status` instead
of its `state`.
1. Wrtie the lambda handler to be triggered via an SQS event
1. Add terraform code to deploy the lambda function
1. Right now, it waits for the PR to be merged. If the PR is closed but not merged,
it throws an error.  After the PR is merged, it checks CircleCI.  It throws an
error if any of the jobs associated with the PR's commit fail. It assumes success if
the `apply_terraform` job finishes successfully. Would like to handle all of the
following provisioning request status/failure points:

      1. CircleCI checks of branch successful (possible failure state)
      1. PR approved
      1. PR closed but not approved (failure state)
      1. PR merged
      1. PR closed but not merged (failure state)
      1. CircleCI apply successful (possible failure state, complete state for non-EC2 provisioning)

1. Right now it updates the RITM's state (old deprecated method of indicating automation status).
It adds a comment for errors and a comment indicating provisioning is complete, but would like to do the following:

      1. get Endpoint of RDS and update ticket with this info
      1. update ticket with appropriate NLB/ALB info
      1. update ticket with AWS account numbers for new tenant provisioning request
      1. update/create CMDB CI's as appropriate
      1. create other STASKS not included in current ServiceNow provisioning workflows

1. Consider refactoring into a series of micro-services
1. Move credentials to Secrets Manager instead of environment variables
