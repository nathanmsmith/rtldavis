According to the [Davis Vantage Vue manual](https://support.davisinstruments.com/article/0ics9tab6w-manual-vantage-vue-integrated-sensor-suite-manual-6250-6357),
the update intervals for each sensor are as follows:

Barometric Pressure - 1 min.
Inside Humidity - 1 min.
Outside Humidity - 50 sec.
Dew Point - 10 sec.
Rainfall Amount - 20 sec.
Rain Storm Amount - 20 sec.
Rain Rate - 20 sec.
Inside Temperature - 1 min.
Outside Temperature - 10 sec.
Heat Index - 10 sec.
Wind Chill - 10 sec.
Wind Speed - 2.5 sec.
Wind Direction - 2.5 sec.
Direction of High Speed - 2.5 sec.

An update interval of every 5 seconds seems like a nice balance between
real-time and pragmatism.

```
{
  temperature: {
    temperature: 50.1,
    received_at: 1742769011
  },
  wind: {
    speed: XXX,
    direction: XXX,
    received_at: 1742769011
  },
  sent_at: 1742769011
}
```
