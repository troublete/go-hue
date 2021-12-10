# go-hue
> Control Philips Hue Lamps via Lua

## About

This is just a throw-away tooling to control Hue Lamps with a simple Lua script. 

```bash
hue -interval=2000 script.lua # run script.lua every 2 seconds (omit *-interval=x* to run one-time)
```

The runtime exposes **loop** global variable when ran in interval mode, so one
is able to keep track of the iterations run.

The runtime exposes **bridges** which is a table of all bridges "found";
(format: {bridge_serial(string)=bridge(table)})

## Install

```bash
# install Lua (i.e. checkout, make generic, sudo make install)
make build
sudo make install
```

## Troubleshooting

- Sometimes SSDP does not discover a Hue bridge/a Hue bridge does not
  advertise itself; try re-running the script when this happens.

- When creating a user, it is necessary to first hit the 'link' button on the
  bridge; if you don't, an error will tell you to do so. 

## Related

- [Bridge Device Limits](https://developers.meethue.com/develop/application-design-guidance/bridge-maximum-settings/)
- [Querying Limits](https://developers.meethue.com/develop/application-design-guidance/hue-system-performance/)
