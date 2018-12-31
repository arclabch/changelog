# Changelog

`changelog` is a small utility to keep track of the maintenance done to a Unix-like system. It acts as a journal and relies on the sysadmin(s) to stay up to date and complete.

It can be used as part of scripts, both for adding and viewing. For example, as part of the `motd` to show the last changes done to a system, or as part of an orchestration system to keep track of the last actions.

It was tested on FreeBSD, OpenBSD and Linux.

## Usage

Run `changelog help` to see the general help, and `changelog help [show|add]` for more detailed help.

### Add a log entry

Only root can create entries. When using `sudo`/`doas`, `changelog` tries to guess the original user.

**Standard add:**

```sh
# changelog add "The new changelog entry to add."
```

**Add as another user:**

```sh
# changelog add -u foobar "The new changelog entry to add."
```

**Add from stdin:**

```sh
# echo "The new changelog entry to add." | changelog add -s
```

### View the changelog

**View all previous entries:**

```sh
$ changelog view
```

**View only the last 10 entries:**

```sh
$ changelog view -l 10
```

**View the entries done by a specific user:**

```sh
$ changelog view -u foobar
```

Options can be combined.

## Installation

Building `changelog` requires the [Go](https://golang.org/) and [GCC](https://gcc.gnu.org/) compilers and [Git](https://git-scm.com/).

On non-GNU systems such as FreeBSD and OpenBSD, you also require [GNU Make](https://www.gnu.org/software/make/). For those, replace `make` by `gmake` in the commands below.

**Get it:**

```sh
$ go get -u github.com/arclabch/changelog
```

As `go get` will build it by itself without following the makefile, you should also delete the compiled version (`$GOPATH/bin/changelog`) to avoid runtime issues.

**Compile it:**

You must go to the source folder and checkout the version you want to compile. If `$GOPATH` is not defined, try `~/go`.

```sh
$ cd $GOPATH/src/github.com/arclabch/changelog
$ git checkout release/1.0.1
$ make
```

**Install it (as root):**

```sh
# make install
```

**Remove it (as root):**

```sh
# make uninstall
```

(or simply delete `rm /usr[/local]/bin/changelog`.)

Optional: delete `/var/log/changelog.db` to remove its database.

## License

MIT. See `LICENSE` file.