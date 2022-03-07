TODO:
- Check initial values
- Test CalculateDirection in elevator
- Test Merge in orders
- Possibly delay lights
- Add and document NTP
- Make persistant cs storage

To test locally with mocknet
```
sudo iptables -t nat -A OUTPUT -p udp -d 255.255.255.255 -j DNAT --to-destination 127.255.255.255
```
