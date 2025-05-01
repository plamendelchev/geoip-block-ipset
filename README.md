# GeoIP Block with IPSet

`geoip-block-ipset` is a Go program for creating GeoIP-based whitelists using
IPSet. It retrieves GeoIP data from
[RIPEstat](https://stat.ripe.net/docs/02.data-api/country-resource-list.html).

## Installation

### Using RPM Package

1. Download the latest RPM from
   [Releases](https://github.com/plamendelchev/geoip-block-ipset/releases)

2. Install with dnf:

```shell
dnf install https://github.com/plamendelchev/geoip-block-ipset/releases/
download/v0.1.0/geoip-block-ipset-0.1.0-1.el9.x86_64.rpm
```

### Using Binary

1. Get the latest binary from
   [Releases](https://github.com/plamendelchev/geoip-block-ipset/releases)

2. Install to your $PATH:

```shell
wget https://github.com/plamendelchev/geoip-block-ipset/releases/download/
v0.1.0/geoip-block-ipset -O /usr/local/bin/geoip-block-ipset
```

3. Set execute permissions:

```shell
chmod u+x /usr/local/bin/geoip-block-ipset
```

## Usage

```shell
Usage: geoip-block-ipset <command> [flags]  

Tool for managing country-based IP whitelisting with IPSet  

Flags:  
  -h, --help                              Display help  
  -c, --config="/etc/geoip-block.conf"    Specify config file path  
  -d, --debug                             Enable debug logging  

Commands:  
  create [flags]    Generate GeoIP block rules  
  delete [flags]    Remove GeoIP block rules  

Run "geoip-block-ipset <command> --help" for details. 

```

## Configuration

After installation, create a config file (default: `/etc/geoip-block.conf`) with
an `allowed_countries` key. This accepts ISO-3166 Alpha-2 country codes, and the
tool will create an IPSet for each.

Example:

```shell
allowed_countries=BG,US
```

## Running the Service

- RPM users: Enable the systemd timer:

  ```shell
  systemctl enable --now geoip-block-ipset.timer
  ```

- Binary users: Run manually and set up a cron job for updates:

  ```shell
  geoip-block-ipset create  
  echo "0 0 * * * root /usr/local/bin/geoip-block-ipset create" >> /etc/crontab
  ```

- Verifying IPSet

  Check created sets with:

  ```shell
  ipset list --name | grep ^geoip_
  ```

- Firewall Rules

  Add IPTables rules to restrict traffic. Example (block all except whitelisted
  countries):

  ```shell
  -P INPUT DROP  
  -A INPUT -m conntrack --ctstate RELATED,ESTABLISHED -j ACCEPT  
  -A INPUT -i lo -j ACCEPT  
  -A INPUT -m conntrack --ctstate INVALID -j DROP  
  -A INPUT -m set ! --match-set geoip_allow_bg src -j DROP  
  -A INPUT -p icmp -m conntrack --ctstate NEW -m icmp --icmp-type 8 -j ACCEPT  
  -A INPUT -p tcp -m tcp --dport 443 -j ACCEPT  
  -A INPUT -p udp -j REJECT --reject-with icmp-port-unreachable  
  -A INPUT -p tcp -j REJECT --reject-with tcp-reset  
  -A INPUT -j REJECT --reject-with icmp-proto-unreachable  
  ```
