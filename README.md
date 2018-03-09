Two simple programs that help when drafting SSSD release notes.

# Build
The programs are written in Go because Jakub wanted an excuse to write something in Go.

To build the programs, run: `go build`:
```
go build detailed-changelog.go
go build pagure-list.go
```

# Usage
To format the closed issues for a release, run:
```
./pagure-list --project=sssd/sssd --milestone="SSSD 1.16.1"
```

To format the detailed changelog for a release, run:
```
./detailed-changelog --directory=/home/jhrozek/devel/sssd --from-tag=sssd_1_16_0
```
