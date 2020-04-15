# acispanctl
**ACI Commander based tool for ERSPAN sessions**

1. Prerequisites
    i. install go
    ii. install go-dep
2. clone the repo
```
mkdir -p $HOME/work/src
cd $HOME/work/src
git clone https://github.com/noironetworks/acispanctl.git
```
3. Install dependencies
```
cd acispanctl
export GOPATH="$HOME/work"
dep ensure -v
```
4. Install acispanctl
```
cd cmd/acispanctl
go install .
```
5. Running the command
```
PATH=$PATH:$GOPATH/bin
acispanctl --help
```
6. Providing credentials to the tool
use $HOME/work/src/acispanctl/sampleconfig to create the credentials file
```
vi $HOME/work/src/acispanctl/sampleconfig
```

Known Issues:
1. The command does not retrieve or configure the SPAN sessions on ACI. Root cause: REST API calls to APIC timesout
Workaround: Rerun the command
2. The POST to bind the src/dest grp to the associated channels in the AEP has not been implemented
