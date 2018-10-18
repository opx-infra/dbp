# Forking from an Upstream

Depending on your upstream, there are three main ways of interacting with Debian package repositories:

1. No upstream: Debian native package
1. Upstream is Debian (or a Debian downstream like Ubuntu)
1. Upstream is not Debian

Upstream type determines how our package versions are bumped. Version examples:

| Upstream type | Upstream version | Upstream tag | Release tag |
| --- | --- | --- | --- |
| not Debian |  `4.5.6` | `v4.5.6` | `debian/4.5.6-0opx1` |
| Debian | `4.5.6-2` | `debian/4.5.6-2` | `debian/4.5.6-2opx1` |
| Ubuntu | `4.5.6-2ubuntu3` | `debian/4.5.6-2ubuntu3` | `debian/4.5.6-2ubuntu3opx1` |

## Debian native package
1. Develop on master
2. Release on `debian/X.Y.Z` tags

## Upstream is Debian
Development is done on the `opx/master` branch. Development is done through Debian patches with [*gbp pq*](https://honk.sigxcpu.org/projects/git-buildpackage/manual-html/gbp.patches.html).

### Workflow for new upstream versions
1. Merge upstream (Debian) tag into `opx/master` branch.
1. Resolve conflicts
1. Create/adapt packaging
	- Add changelog entry with `gbp dch`
	- Reset `opxX` to `opx1`
1. Raise pull request
1. Merge, tag, and release

### Workflow for new OPX changes
1. Follow normal OPX workflow, forking from `opx/master`
1. Import patches with `gbp pq import`
1. Make changes (to source only, no packaging)
1. Export patches with `gbp pq export`
1. Create/adapt packaging
	- Add changelog entry with `gbp dch` (keep same upstream version, only increment final `opxX`
1. Raise pull request
1. Merge, tag, and release

## Upstream is not Debian
Development is done on the `opx/master` branch. No patches or patch queues are used.

### Workflow for new upstream versions
1. Merge upstream tag into `opx/master` branch.
1. Resolve conflicts
1. Create/adapt packaging
	- Add changelog entry with `gbp dch`
	- Reset `opxX` to `opx1`
1. Raise pull request
1. Merge
1. Tag and release

### Workflow for new OPX changes
1. Follow normal OPX workflow, forking from `opx/master`
1. Make changes
1. Create/adapt packaging
	- Add changelog entry with `gbp dch` (keep same upstream version, only increment final `opxX`
1. Raise pull request
1. Merge
1. Tag and release

## Demo of creating a Debian fork and an initial release
```bash
PKG=net-snmp
UPSTREAM=5.7.3+dfsg-1.7
OPX_VER=${UPSTREAM}opx1
git clone -o upstream https://salsa.debian.org/debian/net-snmp.git

# Create our development branch
git -C $PKG checkout -b opx/master debian/$UPSTREAM

# Create new changelog for initial OPX release
# Add new changelog entry, changing 5.7.3+dfsg-2 to 5.7.3+dfsg-1.7opx1
dbp shell -c "cd $PKG; gbp dch --release"
git -C $PKG add debian/changelog

# Commit our changes
git -C $PKG commit -sm 'Initial commit for OPX'

# Build it! And tag it on successful build
dbp --dist stretch --no-extra-sources build $PKG --gbp='--git-tag'

# Push our changes and remove the repo
git -C $PKG remote add origin https://github.com/theucke/net-snmp
git -C $PKG push origin opx/master
git -C $PKG push origin debian/$OPX_VER
rm -rf $PKG
```

## Demo of modifying the Debian fork
```bash
# Clone from our fork
git clone https://github.com/theucke/net-snmp
git -C $PKG checkout opx/master

# Create dev branch
git -C $PKG checkout --track opx/master -b bugfix/python-indentation

# Import patch queue
dbp shell -c "cd $PKG; gbp pq import"

# Fix tab character in python file for something to change
sed -i $'s/\t/        /' $PKG/python/setup.py
git -C $PKG add python/setup.py

# Commit changes and export patch queue
git -C $PKG commit -asm "python: Fix indentation"
dbp shell -c "cd $PKG; gbp pq export"

# Add new changelog entry, changing 5.7.3+dfsg-1.7opx1 to 5.7.3+dfsg-1.7opx2
# Version should be changed automatically for you
dbp shell -c "cd $PKG; gbp dch --release"
git -C $PKG commit --amend -a

# Raise pull request
git pull-request
```

## Demo of merging in new upstream release
```bash
# Clone from our fork
git clone https://github.com/theucke/net-snmp
git -C $PKG checkout opx/master

# Import new upstream version and resolve conflicts
# Move new changelog entry to Debian changelog (must be a better way, Git didn't detect rename)
git -C $PKG merge debian/5.7.3+dfsg-1.8

# Add new changelog entry, changing 5.7.3+dfsg-1.7opx2 to 5.7.3+dfsg-1.8opx1
dbp shell -c "cd $PKG; gbp dch --release"
git -C $PKG commit --amend -a

# Build for verification (normally CI will do this)
dbp --dist stretch --no-extra-sources build $PKG --gbp='--git-tag'

# Push our new tag
git -C $PKG push origin opx/master
git -C $PKG push origin debian/5.7.3+dfsg-1.8opx1
```
