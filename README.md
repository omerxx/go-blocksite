# Go Block Site
Go blocksite is a minimal, configurable, FREE website blocker.

It's been designed as a small app to block unwanted websites on the system level.
GBS was born after a long usage of half-baked browser extensions that were limited, unextendable and quite pricy for their very simple job.


## Usage
The way GBS is working is by changing the hosts file on the machine, rerouting targets to localhost instead of the real site.
For that, it requires that the hosts file has write permissions. The easiest way to achieve that is by running as `sudo` (`sudo go-blocksite`).

The app is configurable by a `config.yml` where you can set a list of block targets, and/or by command line arguments such as `--block` or `--unblock`.

```
Usage of go-blocksite:
  -b, --block string     blocks a given url
  -u, --unblock string   unblocks a given url
  -U, --unblock-all      unblocks all urls
```

## Installation
Download the latest binary compiled to the desired architecture, and put it in the bin folder so it's part of your user path:
```bash
wget https://github.com/omerxx/go-blocksite/downloads... && chmod +x go-blocksite && mv go-blocksite /usr/local/bin/go-blocksite
```

## State
In order to keep things simple and light, a statefile named `/etc/gbs-state.json` is created and maintained to track block / unblocked websites. The file's name and path are also configurable through the optional

