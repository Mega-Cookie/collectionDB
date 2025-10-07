Name:           collectionDB
Version:        0.0.2alpha2
Release:        1%{?dist}
Summary:        A simple webapp for managing your collections of physical media written in Go with SQLite.

License:        GPLv3
URL:            https://github.com/Mega-Cookie/collectionDB
Source0:        collectionDB-0.0.2alpha2.tar.gz

BuildRequires:  golang
BuildRequires:	systemd-rpm-macros

Provides:	%{name}-%{version}

%description
A simple webapp for managing your collections of physical media written in Go with SQLite.

%global debug_package %{nil}

%prep
%autosetup

%build
go build -v -o %{name}

%install
install -Dpm 0755 VERSION %{buildroot}%{_sysconfdir}/%{name}/VERSION
install -Dpm 0755 rpm/config.yml %{buildroot}%{_sysconfdir}/%{name}/config.yml
install -Dpm 0755 %{name} %{buildroot}%{_bindir}/%{name}
install -Dpm 644 %{name}.service %{buildroot}%{_unitdir}/%{name}.service
mkdir -p %{buildroot}%{_sharedstatedir}/%{name}/
mkdir -p %{buildroot}%{_sysconfdir}/%{name}/templates/
cp -a templates/* %{buildroot}%{_sysconfdir}/%{name}/templates/

%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%files
%dir %{_sharedstatedir}/%{name}
%dir %{_sysconfdir}/%{name}
%{_sysconfdir}/%{name}/VERSION
%{_sysconfdir}/%{name}/config.yml
%{_bindir}/%{name}
%{_unitdir}/%{name}.service
%{_sysconfdir}/%{name}/templates/*
%changelog
* Wed Oct 01 2025 Mega-Cookie - pre.alpha.1
- Implementation of basic functionalities by @Mega-Cookie
+ Fri Oct 03 2025 Mega-Cookie - 0.0.2alpha2
- Display server time in view and edit pages by @Mega-Cookie
- Set listening address, port, database file and debug mode through config file by @Mega-Cookie
- Save GO version in database by @Mega-Cookie
- Add logging by @Mega-Cookie
- Add edit and back buttons by @Mega-Cookie
