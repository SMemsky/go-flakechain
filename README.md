# Go-Flakechain

Note: this is project is still in development.

TODO: Write this.

TODO: Rename golang namespaces to github address

# Build

Minimal Go version supported is 1.11.

Building for 32-bit platforms might not succeed.

Base stuff
```
go get -u github.com/SMemsky/go-flakechain/cmd/daemon
```

To run (1.10 won't work, sorry :) )
```
go run github.com/SMemsky/go-flakechain/cmd/daemon
go run github.com/SMemsky/go-flakechain/cmd/wallet
```

Should also work on windows. Just remember to install latest Go and Git.

# Fork?```
25
Or
26
``
27
```

Just a quick note here. It's a golang and all import paths depend on this repository.
To make a fork (for pull request?) do following

Note: Based on [this](http://blog.campoy.cat/2014/03/github-and-go-forking-pull-requests-and.html)

1) Make an empty repo on your github account
2) `go get github.com/SMemsky/go-flakechain`
3) `cd $GOPATH/src/github.com/SMemsky/go-flakechain`
4) `git remote add yourname-fork https://github.com/yourname/reponame`
5) `git pull --rebase yourname-fork`
6) `git push yourname-fork`
