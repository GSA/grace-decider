# GRACE Decider

Program to monitor status of an IaaC provisioning request for GRACE PaaS

# WIP

This is a work-in-progress.  See [TODO](https://github.com/GSA/grace-decider/TODO.md)

## Repository contents

**decider**: Source code for Go Program
**cmd**: Go command line wrapper for handler module
**lambda**: Go lambda function wrapper for handler module

## Usage (command line)

1. Download latest release from [GitHub](https://github.com/GSA/grace-decider/releases)
2. Unzip
3. Copy binary executable file to a directory in your `$PATH` and make executable
4. Set the following environment variables:

```
export GITHUB_TOKEN=<personal access token>
export CIRCLE_TOKEN=<personal api key>
export SN_USER=<ServiceNow Username>
export SN_PASSWORD=<ServiceNow Password>
export SN_INSTANCE=<ServiceNow Instance>
```

5. grace-decider --help:

```
-owner string
    GitHub repository owner (default "GSA")
-pr int
    Pull Request Number
-repo string
    GitHub repository name
-request string
    JSON input file
```

6. Example:

```
grace-decider --request=/tmp/RITM0769780.json --repo="g-grace" --pr="233"
```

## Usage (Lambda function)

tbd

## Public domain

This project is in the worldwide [public domain](LICENSE.md). As stated in [CONTRIBUTING](CONTRIBUTING.md):

> This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
>
> All contributions to this project will be released under the CC0 dedication. By submitting a pull request, you are agreeing to comply with this waiver of copyright interest.
