# Generated by go2rpm 1.14.0.post0
%bcond check 1
%bcond bootstrap 0

%if %{with bootstrap}
%global debug_package %{nil}
%endif

%if %{with bootstrap}
%global __requires_exclude %{?__requires_exclude:%{__requires_exclude}|}^golang\\(.*\\)$
%endif

# https://github.com/plamendelchev/geoip-block-ipset
%global goipath         github.com/plamendelchev/geoip-block-ipset
Version:                0.1.0

%gometa -L -f

%global common_description %{expand:
Simple tool to whitelist countries in ipset.}

%global golicenses      LICENSE
%global godocs          README.md

Name:           golang-github-plamendelchev-geoip-block-ipset
Version:        0
Release:        %autorelease
Summary:        Simple tool to whitelist countries in ipset

License:        GPL-3.0-only
URL:            %{gourl}
Source:         %{gosource}

%description %{common_description}

%gopkg

%prep
%goprep -A
%autopatch -p1

%if %{without bootstrap}
%generate_buildrequires
%go_generate_buildrequires
BuildRequires:	systemd-rpm-macros
%endif

%if %{without bootstrap}
%build
for cmd in cmd/* ; do
  %gobuild -o %{gobuilddir}/bin/$(basename $cmd) %{goipath}/$cmd
done
%endif

%install
%gopkginstall
%if %{without bootstrap}
install -m 0755 -vd                     %{buildroot}%{_libexecdir}
install -m 0755 -vp %{gobuilddir}/bin/* %{buildroot}%{_libexecdir}/

install -m 0755 -vd			%{buildroot}%{_unitdir}
install -m 0644 -vp systemd/*		%{buildroot}%{_unitdir}/
%endif

%if %{without bootstrap}
%if %{with check}
%check
%gocheck
%endif
%endif

%post
%systemd_post geoip-block-ipset.timer

%preun
%systemd_preun geoip-block-ipset.timer

%postun
%systemd_postun geoip-block-ipset.timer


%if %{without bootstrap}
%files
%license LICENSE
%doc README.md
%{_libexecdir}/geoip-block-ipset
%{_unitdir}/geoip-block-ipset.{service,timer}
%endif

%gopkgfiles

%changelog
%autochangelog
