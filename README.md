# repo-installer
clones and installs repositories that have make install command available at the root directory

# Usage

Clone the repository and cd into the repository directory

`git clone https://github.com/TomOgoma/repo-installer.git`

`cd repo-installer`

Run the installer with the help flag to get usage information

`./installer -help`

Here is a template for typical usage:

`./installer -[github|bitbucket] <account>/<repository>`

where `account` is the account name e.g. `tomogoma` in the URL
`https://github.com/TomOgoma/repo-installer.git` and `repository`
is the repository name to fetch e.g. `repo-installer` in the same URL
