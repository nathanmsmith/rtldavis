# rtldavis

## About this Repository

This is a fork of [mdickers47/rtldavis](https://github.com/mdickers47/rtldavis) which is a fork in turn of [lheijst/rtldavis](https://github.com/lheijst/rtldavis) and [bemasher/rtldavis](https://github.com/bemasher/rtldavis).

The way the lheijst and bemasher versions work is that they implement the Davis frequency hopping,
demodulate the packets, and dump them to the log file as hex bytes.
Something else has to decode the messages to create weather data.
Luc wrote [https://github.com/lheijst/weewx-rtldavis](weewx-rtldavis)
which is in Python and operates as a weewx driver to update a weewx
database.

mdickers47's change was to adapt the packet parsing code from
weewx-rtldavis into the go binary. He also has code to upload data to a Graphite server.
It will be written as one data point per line, where each data point is three columns, which are metric name, value, and UNIX timestamp:

```
wx.davis.windspeed_raw 0 1703133526
wx.davis.winddir 143.74 1703133526
wx.davis.windspeed 0.00 1703133526
wx.davis.temp 50.20 1703133526
```

You can also write the decoded weather data to stdout with the option `-gs -`.  

My change will to be to upload data to a server using HTTP and JSON. For my uses, this will be a Rails server hosted on a VPS.

Tested and used on a Davis Vantage Vue. Installed on macOS, using go 1.24.0.

## Installation

### Debian (Trixie)

Download the `.deb` for your architecture (amd64 or arm64) from the
[latest release](https://github.com/nathanmsmith/rtldavis/releases/latest), then:

```
sudo apt install ./rtldavis_*.deb
```

This installs a `rtldavis` systemd service, which is not enabled by default. To set it up:

1. Put your flags and API key in `/etc/default/rtldavis`.
2. Reboot so the blacklisted DVB-T kernel driver releases the dongle.
3. Run `sudo systemctl enable --now rtldavis`.

Logs are available with `journalctl -u rtldavis -f`.

### Releasing

Push a semver tag such as `v0.1.0`. The workflow in `.github/workflows/release.yml` builds and
smoke-tests the `.deb`s, then publishes them to a GitHub release.

### macOS

```
brew install librtlsdr
git clone https://github.com/nathanmsmith/rtldavis/
env CGO_CFLAGS="-I/opt/homebrew/include" CGO_LDFLAGS="-L/opt/homebrew/lib" go install -v .
$GOPATH/bin/rtldavis
```


## Usage

Available command-line flags are as follows:

```
Usage of rtldavis:

  -tr [transmitters]
    	code of the stations to listen for: 
        tr1=1 tr2=2 tr3=4 tr4=8 tr5=16 tr6=32 tr7=64 tr8=128
        or the Davis syntax (first transmitter ID has value 0):
        ID 0=1 ID 1=2 ID 2=4 ID 3=8 ID 4=16 ID 5=32 ID 6=64 ID 7=128
        When two or more transmitters are combined, add the numbers.
        Example: ID0 and ID2 combined is 1 + 4 => -tr 5
        
        Default = -tr 1 (ID 0)

  -tf [tranceiver frequencies]
        EU or US
        Default = -tf EU

  -ex [extra loop_delay in ms]
        In case a lot of messages are missed we might try to use the -ex parameter, like -ex 200
        Note: A negative value will probably lead to message loss
        Default = -ex 0
 
  -fc [frequency correction in Hz for all channels]
        Default = -fc 0
        
  -ppm [frequency correction of rtl dongle in ppm]
        Default = -ppm 0
        
  -maxmissed [max missed-packets-in-a-row before new init]
        Normally you should set this parameter to 4 (-maxmissed 4). 
        During testing of new hardware it may be handy (for US equipment) to leave the default value of 51. 
        The program hops along all channels and present information about each individual channel. 
        Default = -maxmissed 51
        
  -u [log undefined signals]
        The program can pick up (i.e. reveive) messages from undefined transmitters, e.g. from a weather 
        station near-by. De messages are discarded, but you may want to see on which channels they are 
        received and how many.
        Default = -u false
```

### License

The source of this project is licensed under GPL v3.0. See the LICENSE file for details.
