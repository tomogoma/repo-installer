# repo-installer
clones and installs repositories that have make install command available at the root directory

# Usage

Clone the repository and cd into the repository directory

`git clone https://github.com/TomOgoma/repo-installer.git`

`cd repo-installer`

The binary can either be installed:

`make install`

`repo-installer -help`

or run directly from the cloned directory:

`./repo-installer -help`

Here is a template for typical usage:

`repo-installer -[github|bitbucket] <account>/<repository>`

where `account` is the account name e.g. `TomOgoma` in the URL
`https://github.com/TomOgoma/repo-installer.git` and `repository`
is the repository name to fetch e.g. `repo-installer` in the same URL

# Uninstall

To uninstall, run:

`make uninstall`

# Build

To do a manual build run

`make build`
