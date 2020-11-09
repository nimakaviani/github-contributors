# Github Contributors

Gain some loose understanding of who from what companies contribute to a GitHub project.


## Installation

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
+------------+----------------+------------------------+-----------------------------------------------------------------+-------------+-------+
|    ORG     |    GITHUBID    |         EMAIL          |                           ISSUE / PR                            | ASSOCIATION | STATE |
+------------+----------------+------------------------+-----------------------------------------------------------------+-------------+-------+
| gmail.com  | neolit123      | neolit123@gmail.com    | Issue(96378) : vendor: update system-validators to v1.3.0       | MEMBER      | open  |
|            |                |                        | https://api.github.com/repos/kubernetes/kubernetes/issues/96378 |             |       |
+------------+----------------+------------------------+-----------------------------------------------------------------+-------------+       +
| google.com | chelseychen    | chelseychen@google.com | Issue(96379) : Set priority of Event v1 higher than v1beta1     | CONTRIBUTOR |       |
|            |                |                        | https://api.github.com/repos/kubernetes/kubernetes/issues/96379 |             |       |
+            +----------------+------------------------+-----------------------------------------------------------------+-------------+       +
|            | karan          | karangoel@google.com   | Issue(96381) : Fix command and arg in NPD e2e                   | MEMBER      |       |
|            |                |                        | https://api.github.com/repos/kubernetes/kubernetes/issues/96381 |             |       |
+            +----------------+------------------------+-----------------------------------------------------------------+-------------+       +
|            | serathius      | siarkowicz@google.com  | Issue(96374) : [WIP] Create example component                   | CONTRIBUTOR |       |
|            |                |                        | for integrating with component-base                             |             |       |
|            |                |                        | https://api.github.com/repos/kubernetes/kubernetes/issues/96374 |             |       |
+------------+----------------+------------------------+-----------------------------------------------------------------+-------------+       +
| redhat.com | exdx           | dsover@redhat.com      | Issue(96382) : Draining nodes with stand-alone                  | NONE        |       |
|            |                |                        | pods managed by custom controllers                              |             |       |
|            |                |                        | https://api.github.com/repos/kubernetes/kubernetes/issues/96382 |             |       |
+            +----------------+------------------------+-----------------------------------------------------------------+-------------+       +
|            | gnufied        | hekumar@redhat.com     | Issue(96376) : Move fsGroupChangePolicy feature to beta         | MEMBER      |       |
|            |                |                        | https://api.github.com/repos/kubernetes/kubernetes/issues/96376 |             |       |
+            +----------------+------------------------+-----------------------------------------------------------------+-------------+       +
|            | smarterclayton | ccoleman@redhat.com    | Issue(96375) : api: Allow MaxSurge to                           | CONTRIBUTOR |       |
|            |                |                        | be set on DaemonSets during update                              |             |       |
|            |                |                        | https://api.github.com/repos/kubernetes/kubernetes/issues/96375 |             |       |
+------------+----------------+------------------------+-----------------------------------------------------------------+-------------+-------+
```
