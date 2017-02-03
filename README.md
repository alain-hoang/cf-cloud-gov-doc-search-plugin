# cf-cloud-gov-doc-search-plugin

This is a [plugin](https://github.com/cloudfoundry/cli/tree/master/plugin_examples) for [Cloud Foundry CLI](https://github.com/cloudfoundry/cli). This plugin searches the cloud.gov docs for the entered search term and optionally attempt to navigate to the URL provided.


## Pre-Requisites

### Go

Install [Go](https://golang.org/).

On Mac OS X with brew

	$ brew install go


## Building the plugin
To build the plugin make a checkout and use glide to install the dependencies

        $ export GOPATH=$(pwd)
        $ cd src/github.com/alain-hoang/cgds
	$ glide install
	$ go build
	$ cf install-plugin cgds

## Usage
To use the plugin:

	$ cf cloud-gov-doc-search <search term>


## Removal
To remove the plugin:

	$ cf uninstall-plugin "CloudGovDocSearchPlugin"
	


### Public domain

This project is in the worldwide [public domain](LICENSE.md). As stated in [CONTRIBUTING](CONTRIBUTING.md):

> This project is in the public domain within the United States, and copyright and related rights in the work worldwide are waived through the [CC0 1.0 Universal public domain dedication](https://creativecommons.org/publicdomain/zero/1.0/).
>
> All contributions to this project will be released under the CC0 dedication. By submitting a pull request, you are agreeing to comply with this waiver of copyright interest.