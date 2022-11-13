# iploc
Track IP Countries (Offline) / Simple program to find country of IP ...

# Installation
- First clone ( or Download ZIP file and extract it) and build binary file by `go build .`.

```bash
git clone https://github.com/awolverp/iploc && cd iploc && go build .
```

## Usage
Use `iploc --help` to get usage:
```yaml
Usage:
  ./iploc [-h | -list | -all] [OPTIONS] QUERY

Required:
  QUERY              IP/Country name

Options:
  -ns                Don't show summery of results
  -silent            Silent output
  -offset N          Set offset for results
  -limit N           Set limit for results

Output:
  -format FORMAT     Set output format ( json, csv, default )

Other:
  -list             List all countries
  -all              Show all results
```

> **!NOTE: geoip.csv file must be in current directory which you run `./iploc`**

## Examples
**Can search IP**:
```bash
$ ./iploc 1.1.1.1
 ___   ____    _ 
|_ _| |  _ \  | |       ___     ___ 
 | |  | |_) | | |      / _ \   / __|
 | |  |  __/  | |___  | (_) | | (__ 
|___| |_|     |_____|  \___/   \___|
 

---------> 1.1.1.1
1.1.1.0-1.1.1.255		apnic	1313020800	Australia [AU / AUS]

- found 1 = (1) results / 169.816µs
```
**Or country name:**
```bash
$ ./iploc "TR"
 ___   ____    _ 
|_ _| |  _ \  | |       ___     ___ 
 | |  | |_) | | |      / _ \   / __|
 | |  |  __/  | |___  | (_) | | (__ 
|___| |_|     |_____|  \___/   \___|
 

---------> TR
2.56.60.0-2.56.63.255		ripencc	1552435200	Turkey [TR / TUR]
2.56.152.0-2.56.155.255		ripencc	1552521600	Turkey [TR / TUR]
2.57.188.0-2.57.191.255		ripencc	1553040000	Turkey [TR / TUR]
...

- found 1295 = (1295) results / 254.770107ms
```
