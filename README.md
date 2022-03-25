# gdoor

Gdoor is a macOS red team emulation tool developed by CyCraft Technology. It provides a control panel to perform attacks on every client connected to it. It can perform advanced persistence threats with other initial access techniques such as [CVE-2018-6574](https://github.com/frozenkp/CVE-2018-6574) which we used to demo in our presentation. We utilized it to construct a macOS cyber range to evaluate the blue team tools. With our MITRE ATT&CK logging system, the timing and operations can be a reference answer for detection and make the evaluation progress easier.

![Demo](/demo.gif)

## Build

### Client

Since the persistence mechanism varies on different operating systems, we suggest compiling it for macOS, while for debugging or testing usage, it's possible to use it on Linux. There are some operating system-dependent functions. Thus, it doesn't support Golang native cross-compilation. Please compile it on macOS with the makefile.

```
make client
```

### Server

In our scenario, the C2 server of gdoor is on a Linux server, while it can actually be developed on any platform that the Golang compiler supports. Please use the command below to compile it.

```
make server
```

## Configuration

In default, the client connects to `newton.cycarrier`, which we hard-coded in `/etc/hosts`, and the server-side provides two ports to connect. One is for socket connection, and the other is for the file server. 
To alter the configuration, please refer to the two config files named `config_release.go` and `config_debug.go` under `config` directory, and they are for client and server in each. Please modify the three variables `ServerIP`, `ServerPort`, `FServerPort` to your preferred value.

## Encrypted Communication

We used RSA for key exchange and ChaCha20 for the following encrypted communication. The RSA key pair is embedded in the two binaries during compilation. Please use `keypair_generator.go` to generate the key pair, and then put them into the `config` directory before compiling.

```
go run keypair_generator.go
mv private.key config
mv public.key config
```

## Persistence

Upon executing, Gdoor moves itself to `$HOME/.default` as the name `Dropbox.app`. Then, it registers itself as a Launch Agent by adding a plist file named `com.mac.host` under `$HOME/Library/LaunchAgents`. Launch Agent will start the client automatically after boot. If the current user is root, it registers itself as a Launch Daemon by adding the file under `/Library/LaunchDaemons`, which will start the client with root privilege on boot.

Therefore, please delete these files and reboot if you want to remove the gdoor client entirely.

```
/Users/$USER/.default/Dropbox.app
/Users/$USER/Library/LaunchAgents/com.mac.host.plist
/var/root/.default/Dropbox.app
/Library/LaunchDaemons/com.mac.host.plist
```

## SSH Infection

Numerous users use the ssh configuration to enhance their efficiency and convenience, while they don't notice the security problem of keeping their keys on the client. The best practice is to apply a passphrase on the private key, and it will ensure that the key is useless even begin leaked. 

Gdoor utilizes the config files and those unencrypted keys, and it can log in and infect their remote servers easily with the keys. If the remote server is also macOS, it will copy and execute itself on the remote server. 

## Reference
For more detail, please refer to our paper and talk.
- [ASIACCS'20 - POSTER: Construct macOS Cyber Range for Red/Blue Teams](https://dl.acm.org/doi/abs/10.1145/3320269.3405449)
- [HITB CyberWeek 2020 - Constructing An OSX Cyber Range For Red And Blue Teams](https://cyberweek.ae/2020/constructing-an-os-x-cyber-range-for-red-blue-teams/)
