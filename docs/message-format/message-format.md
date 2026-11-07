# Message Format

https://github.com/dekay/DavisRFM69/wiki/Message-Protocol


## Byte 0

```
20 01 ba d4 c3 8a cc b2
^
The header. Values 1-E. Many unknown.
2 = Supercap voltage (Vue only)
3 = ?
4 = UV Index
5 = Rain rate
6 = Solar radiation
7 = Solar Cell output (Vue only)
8 = Temperature
9 = Wind gust
a = Humidity
e = Rain
```

```
50 02 c5 ff 71 0b 7e 96
 ^
0000 = 0

0000
^^^
Bits 0-2 identify which transmitter is sending. Mine is 0.

0000
   ^
Bit 3 flags if the battery is low.
```

## Byte 1

```
a0 01 d6 2f 29 07 a2 8d
   ^^
Wind speed, although not for some messages. ("observed" wind speeds over 65mph
in Jan-Feb 2026:

b8 a0 80 a4 f8 42 92 00
30 86 e3 de 51 3b 1b cc
20 79 7d 2b a3 81 44 b2
a8 9a c0 32 23 2f 12 17
30 86 96 f1 3a 1c cc 5e
80 63 37 e2 66 45 29 56
18 86 81 65 05 f8 b2 66
98 64 88 da da e1 be 71
d8 66 ce d4 0d d0 68 24
10 8b d2 59 27 5b 25 70
50 a0 c2 ad 50 a6 20 62

Message types:
1 - unknown
2 - known (supercap)
3 - unknown
5 - known (rain rate)
8 - known (temperature)
9 - known, but not used (wind gust)
a - known, humidity
b - unknown
d - unknown

```

256 bits, could wind speed be a nibble? To 127mph?

* No. `30 86 e3 de 51 3b 1b cc`, `20 79 7d 2b a3 81 44 b2` occur around 0mph
  wind readings otherwise.

Claude's idea (which I buy):
> **Most likely cause:** Your station is picking up a **co-channel interferer** — another device transmitting on the same frequency (common in the 433 MHz or 915 MHz ISM bands that weather stations use). The collision produces a garbled packet that your decoder interprets as a valid frame with wild values, and the resulting desynchronization causes the gap. Alternatively, the weather station itself could be sending a different message type periodically (like a battery status or sensor ID broadcast) that your decoder is misinterpreting as a wind packet.

Looking at the messages, ~half have `8` in the lower nibble of byte 1, ~half
have `0`. The 8s are definitely interference, the 0s could be.

TODO: add station check and sanity validations to wind speed
