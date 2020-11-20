# gdoor

Gdoor is a red team emulation tool deveoped by CyCraft Technology. We use it to construct a macOS cyber range for red team and blue team. We published our cyber range on ASIACCS'20 and HITB CyberWeek 2020.

## Build

### Client
Gdoor can be compiled on both Linux and macOS, while Linux version is only for debugging since the root module and keeping persistence approaches are all environment dependent.

For general macOS version, please use the command below on macOS.
```
make client_mac
```

For Linux debug version, please use the command below on Linux.
```
make client_linux_debug
```

### Server
In our scenario, C2 server of gdoor is on a Linux server, while it can actually be deployed on any platform that golang compiler supports.

Please use the command below to compile.
```
make server
```
## Reference
For more detail, please refer to our paper and talk.
- [ASIACCS'20 - POSTER: Construct macOS Cyber Range for Red/Blue Teams](https://dl.acm.org/doi/abs/10.1145/3320269.3405449)
- [HITB CyberWeek 2020 - Constructing An OSX Cyber Range For Red And Blue Teams](https://cyberweek.ae/2020/constructing-an-os-x-cyber-range-for-red-blue-teams/)
