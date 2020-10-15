# GRACE Decider

Program to monitor status of an IaaC provisioning request for GRACE PaaS

## Repository contents

**cmd**: Source code for Go Program

## Usage

1. Download latest release from [GitHub](https://github.com/GSA/grace-decider/releases)
2. Unzip
3. Copy binary executable file to a directory in your `$PATH` and make executable
4. Set the following environment variables:

```
export GITHUB_TOKEN=<personal access token>
export CIRCLE_TOKEN=<personal api key>
```

4. Run the command:

```
grace-decider --owner <Repo Owner> --repo <Repo Name> --pr <PR Number>
```

## Public domain

This project is in the worldwide [public domain](LICENSE.md). As stated in [CONTRIBUTING](CONTRIBUTING.md):

> This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
>
> All contributions to this project will be released under the CC0 dedication. By submitting a pull request, you are agreeing to comply with this waiver of copyright interest.
