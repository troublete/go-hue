# go-hue
> Control Philips Hue Lamps via Lua

## About

This is just a throw-away tooling to control Hue Lamps with a simple Lua script. 

```bash
hue script.lua -interval 2000 # run script.lua every 2 seconds (leave empty to run one-time)
```

## Install

```bash
# install Lua (i.e. checkout, make generic, sudo make install)
make build
sudo make install
```

## Related

- [Bridge Device Limits](https://developers.meethue.com/develop/application-design-guidance/bridge-maximum-settings/)
- [Querying Limits](https://developers.meethue.com/develop/application-design-guidance/hue-system-performance/)


