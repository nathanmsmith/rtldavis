# rtldavis

## About this Repository

This fork (mdickers47/rtldavis) is a fork of
[http://github.com/lheijst/rtldavis].  The way the lheijst and bemasher
versions work is that they implement the Davis frequency hopping,
demodulate the packets, and dump them to the log file as hex bytes.
Something else has to decode the messages to create weather data.
Luc wrote [https://github.com/lheijst/weewx-rtldavis][weewx-rtldavis]
which is in Python and operates as a weewx driver to update a weewx
database.

My change is that I have adapted the packet parsing code from
weewx-rtldavis into the go binary.  You can now write the decoded weather
data to stdout with the option `-g -`.  It will be written as one data
point per line, where each data point is three columns, which are metric
name, value, and UNIX timestamp:

```
wx.davis.windspeed_raw 0 1703133526
wx.davis.winddir 143.74 1703133526
wx.davis.windspeed 0.00 1703133526
wx.davis.temp 50.20 1703133526
```

This is the Graphite/Carbon "line receiver" format, so, by supplying
`-w server:udpport`, this rtldavis binary will send the data directly
to Graphite via UDP packets.

*Beware* that I have only implemented the parsing for the Vantage Vue
ISS sensor package with the 0.2 inch rain bucket, because that is what
I have.  If the `-w` option is used with a different weather station,
it will be wrong.

- 2023-12-19 M. Dickerson

From here on is the README from lheijst/rtldavis.

This repository is a fork of [https://github.com/bemasher/rtldavis](https://github.com/bemasher/rtldavis) for use with custom receiver modules. It has been modified in numerous ways.
1) Added EU frequencies
2) Handling of more than one concurrent transmitters.
3) Output format is changed for use with the weewx-rtldavis driver which does the data parsing. 


## Installation

#### Install packages needed for rtldavis

    sudo apt-get install golang git cmake librtlsdr-dev

#### Setup Udev Rules

Next, you need to add some udev rules to make the dongle available for the non-root users. First you want to find the vendor id and product id for your dongle.
The way I did this was to run:

    lsusb

The last line was the Realtek dongle:
    Bus 001 Device 008: ID 0bda:2838 Realtek Semiconductor Corp.
    Bus 001 Device 005: ID 0bda:2838 Realtek Semiconductor Corp. RTL2838 DVB-T

The important parts are "0bda" (the vendor id) and "2838" (the product id).

Create a new file as root named /etc/udev/rules.d/20.rtlsdr.rules that contains the following line:
    nano /etc/udev/rules.d/20.rtlsdr.rules
    SUBSYSTEM=="usb", ATTRS{idVendor}=="0bda", ATTRS{idProduct}=="2838", GROUP="adm", MODE="0666", SYMLINK+="rtl_sdr"

With the vendor and product ids for your particular dongle. This should make the dongle accessible to any user in the adm group. and add a /dev/rtl_sdr symlink when the dongle is attached.

#### Get librtlsdr

    cd /home/pi
    git clone https://github.com/steve-m/librtlsdr.git
    cd librtlsdr
    mkdir build
    cd build
    cmake ../ -DINSTALL_UDEV_RULES=ON
    make
    sudo make install
    sudo ldconfig

#### Change / create ~/profile

    sudo nano ~/.profile
    add at the end of the file:
    
    export GOROOT=/usr/lib/go
    export GOPATH=$HOME/work
    export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
    
#### Activate changed / new ~/profile
    source ~/.profile

#### Get the rtldavis package

    cd /home/pi
    go get -v github.com/lheijst/rtldavis

#### Compiling GO sources

    cd $GOPATH/src/github.com/lheijst/rtldavis
    git submodule init
    git submodule update
    go install -v .

#### Start program rtldavis

    $GOPATH/bin/rtldavis

#### Usage

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

The source of this project is licensed under GPL v3.0. According to [http://choosealicense.com/licenses/gpl-3.0/](http://choosealicense.com/licenses/gpl-3.0/) you may:

#### Required:

 * **Disclose Source:** Source code must be made available when distributing the software. In the case of LGPL and 
 OSL 3.0, the source for the library (and not the entire program) must be made available.
 * **License and copyright notice:** Include a copy of the license and copyright notice with the code.
 * **State Changes:** Indicate significant changes made to the code.

#### Permitted:

 * **Commercial Use:** This software and derivatives may be used for commercial purposes.
 * **Distribution:** You may distribute this software.
 * **Modification:** This software may be modified.
 * **Patent Use:** This license provides an express grant of patent rights from the contributor to the recipient.
 * **Private Use:** You may use and modify the software without distributing it.

#### Forbidden:

 * **Hold Liable:** Software is provided without warranty and the software author\/license owner cannot be held liable for damages.

### Feedback
If you have any general questions or feedback leave a comment below. For bugs, feature suggestions and anything 
directly relating to the program itself, submit an issue.
