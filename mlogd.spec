Name: mlogd
Version: 1.10.0
Release: 01el8
Packager: Michael P. Soulier <msoulier@digitaltorque.ca>
Summary: An svlogd replacement with more standard unix logging behaviour.
License: MIT
Group: System
Source0: %{name}-%{version}.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-root
AutoReqProv: no
%define __os_install_post %{nil}
%define debug_package %{nil}

%ifarch x86_64
ExclusiveArch: x86_64
%define arch x86_64
# temporary hack to keep mag.py happy
BuildArch: x86_64
%else

%ifarch aarch64
ExclusiveArch: aarch64
%define arch aarch64
%else
%{error:"Unsupported build architecture %{arch}"}
%endif

%endif


%description
This is a multilog/svlogd replacement with behaviour that is more typical of
logging on Unix, using a .log symlink to a <name>-<date>.log file, plus a post
rotation hook that allows compression to a .log.gz file.

%changelog
* Mon Jan 8 2024 Michael Soulier <michael.soulier@mitel.com>
- [1.10.0-01el8]
- Rolling ahead for 12.0.

* Wed Dec 14 2022 Auto build <do-not-reply@mitel.com>
- [1.8.1-01el8]
- a7d55b1 Rolling minor version to make room for a possible SP stream

* Sat Dec 10 2022 Auto build <do-not-reply@mitel.com>
- [1.7.1-01el8]
- 

* Thu Dec 8 2022 Michael Soulier <michael.soulier@mitel.com>
- [1.7.0-01el8]
- Rolling ahead for Rocky 8.

* Mon Aug 22 2022 Auto build <do-not-reply@mitel.com>
- [1.6.1-01el7]
- d526a4d Fixing build instructions in specfile.
- 2041ff4 Turning mlogd into a go module
- aa2f988 rolling ahead after stream split

* Fri Feb 26 2021 Auto build <do-not-reply@mitel.com>
- [1.5.16-01el7]
- af72fc8 Copying mlib into mlogd for now.
- 85ac15e Fixing GOPATH to pick up deps in mitel-msl-tug.
- dfcbebf Adding arch hack to specfile to keep mag.py happy for now.
- 5daa302 Changing arch detection in specfiles.
- ce14868 Changing arch detection in specfiles.
- 91bc47a Fixing if condition parsing when targetplatform is not supplied.
- b2c40e5 Ooops
- 18683e9 Adding a conditional BuildArch to specfile. [MBG-5149]

* Thu Feb 18 2021 Auto build <do-not-reply@mitel.com>
- [1.5.15-01el7]
- af72fc8 Copying mlib into mlogd for now.
- 85ac15e Fixing GOPATH to pick up deps in mitel-msl-tug.

* Thu Feb 18 2021 Auto build <do-not-reply@mitel.com>
- [1.5.14-01el7]
- dfcbebf Adding arch hack to specfile to keep mag.py happy for now.

* Thu Feb 11 2021 Auto build <do-not-reply@mitel.com>
- [1.5.13-01el7]
- 5daa302 Changing arch detection in specfiles.
- ce14868 Changing arch detection in specfiles.

* Thu Jan 28 2021 Auto build <do-not-reply@mitel.com>
- [1.5.12-01el7]
- 91bc47a Fixing if condition parsing when targetplatform is not supplied.
- b2c40e5 Ooops
- 18683e9 Adding a conditional BuildArch to specfile. [MBG-5149]

* Thu Aug 20 2020 Auto build <do-not-reply@mitel.com>
- [1.5.11-01el7]
- 3a95b59 Migrate .whitesource configuration file to inheritance mode

* Thu Aug 20 2020 Auto build <do-not-reply@mitel.com>
- [1.5.10-01el7]
- 3a95b59 Migrate .whitesource configuration file to inheritance mode

* Wed Jul 22 2020 Auto build <do-not-reply@mitel.com>
- [1.5.9-01el7]
- 870c686 Adding udp streaming to mlogd. [MBG-4276]

* Tue Jul 21 2020 Auto build <do-not-reply@mitel.com>
- [1.5.8-01el7]
- 870c686 Adding udp streaming to mlogd. [MBG-4276]

* Tue Jul 7 2020 Auto build <do-not-reply@mitel.com>
- [1.5.7-01el7]
- d36ad6e Add .whitesource configuration file

* Sat Jun 27 2020 Auto build <do-not-reply@mitel.com>
- [1.5.6-01el7]
- d36ad6e Add .whitesource configuration file

* Sat Feb 29 2020 Auto build <do-not-reply@mitel.com>
- [1.5.5-01el7]
- fe9ed80 Making maxage == 0 mean disabled. [MBG-4124:solved]

* Fri Dec 20 2019 Auto build <do-not-reply@mitel.com>
- [1.5.4-01el7]
- a56ae0f Fixing a typo in the -stdout option

* Fri Oct 25 2019 Auto build <do-not-reply@mitel.com>
- [1.5.3-01el7]
- dc027cf Adding the timestamp to the stdout option.
- c2c81d5 Adding a --stdout option.

* Thu Dec 13 2018 Auto build <do-not-reply@mitel.com>
- [1.5.2-01el7]
- 1e070d6 Adding a common timestamp format. [MBG-2913:solved]

* Wed Jul 4 2018 Auto build <do-not-reply@mitel.com>
- [1.5.1-01el7]
- c9542c5 Rolling version for 11.0 stream.
- 244ff39 Rolling version for 10.2 stream.
- 0cb02af Built RPi version.

* Tue Jun 26 2018 Michael P. Soulier <michael.soulier@mitel.com>
- [1.5.0-01el7]
- Bumping version for 11.0 stream.

* Tue Jun 26 2018 Michael P. Soulier <michael.soulier@mitel.com>
- [1.4.0-01el7]
- Bumping version for 10.2 stream.

* Thu Nov 2 2017 Auto build <do-not-reply@mitel.com>
- [1.3.15-01]
- 68bdcd6 Adding a check for deleting the just-rotated file.
- 1de10ab Updating signal handlers to loop indefinitely, except for handle_shutdown.
- 0cbcb2d fixing logger info statement
- 71d35cf Flushing test output
- a8aad80 Adding some debug
- 5ca45e4 Simplifying the fake_input script
- 82d1161 First mlogd .deb package

* Thu Nov 2 2017 Auto build <do-not-reply@mitel.com>
- [1.3.14-01]
- 68bdcd6 Adding a check for deleting the just-rotated file.
- 1de10ab Updating signal handlers to loop indefinitely, except for handle_shutdown.
- 0cbcb2d fixing logger info statement
- 71d35cf Flushing test output
- a8aad80 Adding some debug
- 5ca45e4 Simplifying the fake_input script
- 82d1161 First mlogd .deb package

* Thu Nov 2 2017 Auto build <do-not-reply@mitel.com>
- [1.3.13-01]
- 68bdcd6 Adding a check for deleting the just-rotated file.
- 1de10ab Updating signal handlers to loop indefinitely, except for handle_shutdown.
- 0cbcb2d fixing logger info statement
- 71d35cf Flushing test output
- a8aad80 Adding some debug
- 5ca45e4 Simplifying the fake_input script
- 82d1161 First mlogd .deb package

* Fri Aug 4 2017 Auto build <do-not-reply@mitel.com>
- [1.3.12-01]
- 5e5d6b3 Moving to use os.Stat for file size.
- 510cebc New dep

* Fri Jul 28 2017 Auto build <do-not-reply@mitel.com>
- [1.3.11-01]
- 5e5d6b3 Moving to use os.Stat for file size.
- 510cebc New dep

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
export PATH=$PATH:/usr/local/go/bin
go build

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
