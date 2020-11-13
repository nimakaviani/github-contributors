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
github-contrib fetches contribution info for github

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
  -r, --repo string       project repo
  -u, --unauthenticated   unauthenticated gh call

Use "github-contrib [command] --help" for more information about a command.
```

## Output

Check the list of issues or PRs for insights

```
❯ ./github-contrib issues --repo kubernetes/kubernetes
 10 / 10 [=========================================================================================================================] 100.00% 10s
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+-------+
|      ORG      |    GITHUBID    |           EMAIL            |                           ISSUE / PR                            | ASSOCIATION | STATE |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+-------+
| google.com    | adtac          | adtac@google.com           | Issue(96527) : promote API priority                             | MEMBER      | open  |
|               |                |                            | and fairness types and APIs to beta                             |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96527 |             |       |
+               +----------------+----------------------------+-----------------------------------------------------------------+             +       +
|               | apelisse       | apelisse@google.com        | Issue(96480) : Add defaults openapi                             |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96480 |             |       |
+               +----------------+----------------------------+-----------------------------------------------------------------+             +       +
|               | caesarxuchao   | xuchao@google.com          | Issue(96549) : update golang.org/x/net and golang.org/x/sys     |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96549 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+             +       +
| live.com      | Jefftree       | jeffrey.ying86@live.com    | Issue(96317) : Integrate defaults marker to defaulter-gen       |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96317 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+             +       +
| microsoft.com | andyzhangx     | xiazhang@microsoft.com     | Issue(96293) : azure file migration go beta in 1.20             |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96293 |             |       |
+               +                +                            +-----------------------------------------------------------------+             +       +
|               |                |                            | Issue(96546) : Automated cherry pick of #96355: fix pull        |             |       |
|               |                |                            | image error from multiple ACRs using azure managed              |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96546 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+       +
| redhat.com    | smarterclayton | ccoleman@redhat.com        | Issue(96375) : api: Allow MaxSurge to                           | CONTRIBUTOR |       |
|               |                |                            | be set on DaemonSets during update                              |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96375 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+       +
| uservers.net  | puerco         | adolfo.garcia@uservers.net | Issue(96544) : Update 1.18 changelog                            | MEMBER      |       |
|               |                |                            | with entries from v1.18.11                                      |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96544 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+             +       +
| zackz.dev     | knight42       | i@zackz.dev                | Issue(96550) : chore(gce): pass auth flags to KCM and KS        |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96550 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+             +       +
| gmail.com     | nilo19         | pomelonicky@gmail.com      | Issue(96545) : Update the route table                           |             |       |
|               |                |                            | tag in the route reconcile loop                                 |             |       |
|               |                |                            | https://api.github.com/repos/kubernetes/kubernetes/issues/96545 |             |       |
+---------------+----------------+----------------------------+-----------------------------------------------------------------+-------------+-------+
```
