# git-gerrit

A helper tool for working with Gerrit change sets and git.  

## Status

Mostly started with some simple readonly operations to get started.

- [x] List upstream changes.
  - Ability to provide a custom query: `git gerrit changes -q "status:open"`
- [x] Show details about a change.
  - `gerrit show :changeId`
- [x] Show an activity log of events related to an associated change.
  - `gerrit log :changeId`
- [x] Checkout an existing change set.
  - `gerrit checkout :changeID`
- [ ] Create a new change set.
- [ ] Split an active change set.
- [ ] Push an updated change set.

## Usage

_Note:_ tool must be run from within a repository root currently.

### Installation

```
go install github.com/mjpitz/git-gerrit/cmd/git-gerrit@latest

# for a shorter syntax...
alias gerrit=git-gerrit

# otherwise, you can access using `git gerrit`
```

### Listing available change sets

```
$ gerrit changes

CHANGE ID       SUBJECT                 UPDATED
6359            ...                     Thu, 18 Nov 2021 10:39:54 CST
6218            ...                     Thu, 18 Nov 2021 10:38:06 CST
6298            ...                     Thu, 18 Nov 2021 09:21:59 CST
6354            ...                     Thu, 18 Nov 2021 08:48:22 CST
```

### Checking out a change set

When checking out a change set, a new branch is actually created (gerrit branches )

```
$ gerrit checkout 1234  # or gerrit checkout Change-Id

$ git branch | grep "*" 
* gerrit/6218
```

### Showing a change set

```
$ gerrit show I3cad05d0c2c8c9c9b1cad8b182fb459ccf3732ea

REVIEW:  https://localhost:8080/c/mjpitz/git-gerrit/+/1
SUBJECT: initial version of the tool
OWNER:   mya <noreply@example.com>
CREATED: Fri, 29 Oct 2021 14:44:06 -0500
UPDATED: Wed, 17 Nov 2021 06:52:40 -0600

STATUS:           MERGED    +413 -0
WORK IN PROGRESS: false
MERGEABLE:        false

PROJECT:  mjpitz/git-gerrit
BRANCH:   main
COMMIT:   c1df86cf99132ce83952a2c35e7d76d4e136ca02

REVIEWERS:
- Robot <>
- Person 1 <noreply@example.com>
- Person 2 <noreply@example.com>
- ...

```

### Showing the change log

```
$ gerrit log I3cad05d0c2c8c9c9b1cad8b182fb459ccf3732ea
```
