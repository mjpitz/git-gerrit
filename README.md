# git-gerrit

A helper tool for working with Gerrit changesets and git.  

## Status

Mostly started with some simple readonly operations to get started.

- [x] List upstream changes.
  - Ability to provide a custom query: `git gerrit changes -q "status:open"`
- [x] Show details about a change.
  - `gerrit show :changeId`
- [x] Show an activity log of events related to an associated change.
  - `gerrit log :changeId`
- [ ] Checkout / create a new changeset.
- [ ] Pull updates from a changeset into a local branch.
- [ ] Push a changeset.

## Usage

_Note:_ tool must be run from within a repository root currently.

### Installation

```
go install github.com/mjpitz/git-gerrit@latest

# for a shorter syntax...
alias gerrit=git-gerrit

# otherwise, you can access using `git gerrit`
```

### Listing available change sets

```
$ gerrit changes

UPDATED                         CHANGE ID                                       SUBJECT                                                                                                                                            
2021-11-17 09:40:44 -0600 CST   I3cad05d0c2c8c9c9b1cad8b182fb459ccf3732ea       ...
2021-11-17 05:54:31 -0600 CST   I83d34367aab1f3c0d46a044f54980b2d50174b19       ...
2021-11-16 14:16:38 -0600 CST   I0cc497873eb5732623ef2d9bc5f78ba1cc48c6b8       ...
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
