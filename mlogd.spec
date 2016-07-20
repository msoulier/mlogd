Name: mlogd
Version: 1.14.4
Release: 01
Packager: Michael P. Soulier <msoulier@digitaltorque.ca>
Summary: An svlogd replacement with more standard unix logging behaviour.
License: MIT
Group: System
Source0: %{name}-%{version}.tar.gz
BuildRequires: golang
BuildRoot: %{_tmppath}/%{name}-%{version}-root
BuildArch: x86_64
#AutoReqProv: no
#%define __os_install_post %{nil}

%description
This is a multilog/svlogd replacement with behaviour that is more typical of
logging on Unix, using a .log symlink to a <name>-<date>.log file, plus a post
rotation hook that allows compression to a .log.gz file.

%changelog
*
- []
- Initial rpm build.

%prep
%setup -q

%build

%install
rm -rf $RPM_BUILD_ROOT
(cd root   ; find . -depth -print | cpio -dump $RPM_BUILD_ROOT)

rm -f %{name}-%{version}-%{release}-filelist
/sbin/e-smith/genfilelist \
    --file '/usr/sbin/tug-dump-dbs' 'attr(0700,root,root)' \
    $RPM_BUILD_ROOT > %{name}-%{version}-%{release}-filelist

%clean
rm -rf $RPM_BUILD_ROOT

%files -f %{name}-%{version}-%{release}-filelist
%defattr(-,root,root)

%pre

%post
