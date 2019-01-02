# Changelog

`changelog` is a small utility to keep track of the maintenance done to a Unix-like system. It acts as a journal and relies on the sysadmin(s) to stay up to date and complete.

It can be used as part of scripts, both for adding and viewing. For example, as part of the `motd` to show the last changes done to a system, or as part of an orchestration system to keep track of the last actions.

`changelog` has been tested on FreeBSD, OpenBSD, NetBSD, Solaris and Debian GNU/Linux. It should run on other Unix-like systems as well.

## Usage

Run `changelog help` to see the general help, and `changelog help [show|add]` for more detailed help.

### Add a log entry

Only root can create entries. When using `sudo`/`doas`, `changelog` tries to guess the original user.

`changelog add` doesn't return any output and exits with code 0 if add was successful.

**Standard add:**

```sh
$ sudo changelog add "The new changelog entry to add."
```

**Add as another user:**

```sh
$ sudo changelog add -u foobar "The new changelog entry to add."
```

**Add from stdin:**

```sh
$ sudo echo "The new changelog entry to add." | changelog add -s
```

### View the changelog

**View the 5 (default) previous entries:**

```
$ changelog view
2018-12-31T15:08  arclab  And another one with multiple lines.
                          Like this.
                          Because why not.
2018-12-31T15:07  arclab  A long sentence that will wrap around depending on the terminal width
                          while still maintaining alignment with the column. Isn't it nice?
2018-12-31T15:06  àrçláb  An example UTF-8 message - ℓεετ sρεακ αηλοηε?
2018-12-31T15:05  arclab  This is a simple sample message.
```

**View entries formated to a width of 80 columns:**

```
$ changelog view -w 80
2018-12-31T15:08  arclab  And another one with multiple lines.
                          Like this.
                          Because why not.
2018-12-31T15:07  arclab  A long sentence that will wrap around depending on
                          the terminal width while still maintaining alignment
                          with the column. Isn't it nice?
2018-12-31T15:06  àrçláb  An example UTF-8 message - ℓεετ sρεακ αηλοηε?
2018-12-31T15:05  arclab  This is a simple sample message.
```

**View only the last 10 entries:**

```
$ changelog view -l 10
```

**View the entries done by a specific user:**

```
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
$ git checkout release/1.0.2
$ make
```

**Install it (as root):**

```sh
$ sudo make install
```

**Remove it (as root):**

```sh
$ sudo make uninstall
```

(or simply delete `rm /usr[/local]/bin/changelog`.)

Optional: delete `/var/log/changelog.db` to remove its database.

## What's New

**1.0.2**

- Bug fix: minimum terminal width could be miscalculated.
- Documentation: Added example outputs to this readme.

**1.0.1**

- Bug fix: error on `make install`.

**1.0.0**

- Initial release.

## License

MIT. See `LICENSE` file.