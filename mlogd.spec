Name: mlogd
Version: 1.3.10
Release: 01
Packager: Michael P. Soulier <msoulier@digitaltorque.ca>
Summary: An svlogd replacement with more standard unix logging behaviour.
License: MIT
Group: System
Source0: %{name}-%{version}.tar.gz
BuildRequires: golang
BuildRoot: %{_tmppath}/%{name}-%{version}-root
BuildArch: x86_64
AutoReqProv: no
%define __os_install_post %{nil}
%define debug_package %{nil}

%description
This is a multilog/svlogd replacement with behaviour that is more typical of
logging on Unix, using a .log symlink to a <name>-<date>.log file, plus a post
rotation hook that allows compression to a .log.gz file.

%changelog
* Sat Jul 22 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.10-01]
- b7d8b81 Auto build: updating specfile

* Fri Jul 21 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.9-01]
- 6032de5 Auto build: updating specfile

* Fri Jul 21 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.8-01]
- e3abf3b Auto build: updating specfile

* Thu Jul 20 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.7-01]
- 5c359ca Auto build: updating specfile
- bdfc768 Fixing GOPATH in build.
- aa5cfae Fixing reference to mlib function.
- 80f8f3b Refactoring into mlib

* Wed Jul 19 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.6-01]
- bdfc768 Fixing GOPATH in build.
- aa5cfae Fixing reference to mlib function.
- 80f8f3b Refactoring into mlib

* Tue Jun 6 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.5-01]
- a1c43e2 Fixing version.
- 13bb740 Updating version
- 20d5475 Added a shutdown timer.
- 9af04b3 Adding forced rotation for mlogd through signals.
- f4343ec Fixing 32-bit build.

* Sat May 20 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.4-01]
- fc82c2e Adding deps file

* Mon Apr 10 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.3-01]
- Unknown changes

* Sat Feb 25 2017 Auto Build <do-not-reply@mitel.com>
- [1.3.2-01]
- 

* Wed Dec 7 2016 Auto Build <do-not-reply@mitel.com>
- [1.3.1-01]
- d207196 Rolling ahead for 10.0

* Wed Nov 2 2016 Auto Build <do-not-reply@mitel.com>
- [1.2.11-01]
- None

* Wed Oct  5 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.10-01]
- Fixing path to genfilelist. [MN00650422]

* Fri Sep 30 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.9-01]
- Including go-logging in our src tree.

* Wed Sep 14 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.8-01]
- Fixing a bug in the date format string.

* Wed Sep 14 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.7-01]
- Fixing build of broken symlink on relative path.
- Added assertion if space is found in filename.
- Moved rotation check to a background goroutine so it works when there is no
  input, or with blocking I/O.
- Updated format to ensure no spaces in filenames.

* Mon Sep 12 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.6-01]
- Adding parse of filename to determine creation datetime.

* Thu Sep  1 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.5-01]
- Really fixing mlogd's lack of rotation with no input. Really.

* Tue Aug 30 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.4-01]
- Fixing mlogd's lack of rotation with no input.

* Tue Aug 23 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.3-01]
- Improving detection of newly rotated file to run post on.

* Mon Aug 22 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.2-01]
- Fixing a lack of post action run on startup.

* Wed Jul 20 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2.1-01]
- Adding additional debug.

* Wed Jul 20 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.2-01]
- Fixing mlogd permissions.

* Wed Jul 20 2016 Michael Soulier <msoulier@digitaltorque.ca>
- [1.1-01]
- Initial rpm build.

%prep
%setup -q

%build
export GOPATH=$(pwd):$HOME/work/go
go build -o mlogd

%install
rm -rf $RPM_BUILD_ROOT
mkdir -p $RPM_BUILD_ROOT/usr/sbin
cp mlogd $RPM_BUILD_ROOT/usr/sbin

%clean
rm -rf $RPM_BUILD_ROOT

%files
%defattr(-,root,root)
%attr(0755,root,root) /usr/sbin/mlogd
%doc LICENSE README

%pre

%post
