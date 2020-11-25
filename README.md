# Github Contributors

Gain some loose understanding of who from what companies contribute to a GitHub project.


## Installation

Simply download the binary for your operating system from the [Release
Page](https://github.com/nimakaviani/github-contributors/releases).

Or if you are curious enough to build the code, try below:

```
go get github.com/nimakaviani/github-contributors
```

you can also build it by running

```
> ./hack/build.sh
```

## Usage

Extract the list of contributors by looking at issues or PRs

```
❯ ./github-contrib -h
github-contrib fetches contribution info for a github repo

Usage:
  github-contrib [flags]
  github-contrib [command]

Available Commands:
  help        Help about any command
  issues      Analyze issues
  prs         Analyze PRs
  version     Print the version number of Hugo

Flags:
  -c, --count int         count of items to analyze (default 30)
  -d, --debug             debug mode
  -e, --expand            expand user info (default true)
  -h, --help              help for github-contrib
  -o, --output string     output format (json)
  -r, --repo string       project repo
  -u, --unauthenticated   unauthenticated gh call

Use "github-contrib [command] --help" for more information about a command.
```

## Output

Get higher level insights from the list of contributors

```
❯ ./github-contrib --repo kubernetes/kubernetes --expand=false
Analyzig the top 30 contributors on kubernetes/kubernetes
 30 / 30 [==================] 100.00% 18s
+----------------+-------+------------+
|      ORG       | COUNT | PERCENTAGE |
+----------------+-------+------------+
| microsoft.com  |     1 | 3.3%       |
+----------------+-------+------------+
| redhat.com     |     4 | 13.3%      |
+----------------+-------+------------+
| google.com     |    15 | 50.0%      |
+----------------+-------+------------+
| liggitt.net    |     1 | 3.3%       |
+----------------+-------+------------+
| bedafamily.com |     1 | 3.3%       |
+----------------+-------+------------+
| gmail.com      |     6 | 20.0%      |
+----------------+-------+------------+
```

Or get an extended view when checking contributors, issues, or PRs.

```
❯ ./github-contrib issues --repo kubernetes/kubernetes
 10 / 10 [=========================================================================================================================] 100.00% 10s
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+-------+
|      ORG      |    GITHUBID    |           EMAIL            |                           ISSUE / PR                            | ASSOCIATION | STATE |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+-------+
| google.com    | adtac          | axxxc@google.com           | Issue(96527) : promote API priority                             | MEMBER      | open  |
|               |                |                            | and fairness types and APIs to beta                             |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96527 |             |       |
+               +----------------+----------------------------+-----------------------------------------------------------------+             +       +
|               | apelisse       | axxxxxxx@google.com        | Issue(96480) : Add defaults openapi                             |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96480 |             |       |
+               +----------------+----------------------------+-----------------------------------------------------------------+             +       +
|               | caesarxuchao   | xxxxxx@google.com          | Issue(96549) : update golang.org/x/net and golang.org/x/sys     |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96549 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+             +       +
| live.com      | Jefftree       | jxxxxxxxxxxxxx@live.com    | Issue(96317) : Integrate defaults marker to defaulter-gen       |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96317 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+             +       +
| microsoft.com | andyzhangx     | xxxxxxxx@microsoft.com     | Issue(96293) : azure file migration go beta in 1.20             |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96293 |             |       |
+               +                +                            +-----------------------------------------------------------------+             +       +
|               |                |                            | Issue(96546) : Automated cherry pick of #96355: fix pull        |             |       |
|               |                |                            | image error from multiple ACRs using azure managed              |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96546 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+       +
| redhat.com    | smarterclayton | cxxxxxxx@redhat.com        | Issue(96375) : api: Allow MaxSurge to                           | CONTRIBUTOR |       |
|               |                |                            | be set on DaemonSets during update                              |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96375 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+       +
| uservers.net  | puerco         | axxxxxxxxxxxx@uservers.net | Issue(96544) : Update 1.18 changelog                            | MEMBER      |       |
|               |                |                            | with entries from v1.18.11                                      |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96544 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+             +       +
| zackz.dev     | knight42       | i@zackz.dev                | Issue(96550) : chore(gce): pass auth flags to KCM and KS        |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96550 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+             +       +
| gmail.com     | nilo19         | pxxxxxxxxxx@gmail.com      | Issue(96545) : Update the route table                           |             |       |
|               |                |                            | tag in the route reconcile loop                                 |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96545 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+-------+
```

## Contribute

Of course, contributions are very welcome. File issues or PRs if you find this
useful.
