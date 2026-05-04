# 🛜 IPSec Comparison
The main goal here is to perform a comparison between WireGuard and IPSec and collected performance measurements in their respective tunnels.

## 🛠️ Infrastructure
IPSec was deployed in AWS, similar to how WireGuard was deployed, and the infrastructure code can be found in the `/terraform` directory.

## 🔐 Security Associations
The template can be found here in `security_association_example.conf`, please replace the public IPs with the actual ones, and the PSK with what you generated.
Also, this requires installations at both sites, so for the other site, just reverse the order of the IPs -> A becomes B, and B becomes A

## 🧑‍💻 Steps
First, we need to install `strongswan` for IPSec and `charon`, which is part of `strongswan`, for IKEv2

```bash
sudo apt update
sudo apt install -y charon-systemd strongswan-swanctl iperf3
```

Then, enable the daemons:
```bash
sudo systemctl enable --now strongswan
sudo systemctl status strongswan
```

We need to create the loopback addresses in order to bind the child security association with, so in each peer:
```bash
sudo ip addr add <any_private_ip> dev lo
```

Generate a strong PSK, either by your own or through the following:
```bash
openssl rand -base64 32
```

Then we need to copy the security associations, please take a look in the security association section of this document.
After populating the file, please copy it here (edit the file directly or replace with your own)(to be done in each peer):
`/etc/swanctl/swanctl.conf`

Then, in each host, run the following to install the security associations:
```bash
sudo swanctl --load-all
sudo systemctl restart strongswan
sudo swanctl --list-conns
sudo swanctl --list-sas
```

After installing in both hosts, please run the following in both to check that the connection has been established:
```bash
sudo swanctl --list-sas
```

To test connection, just ping through the loopback address like this:
```bash
ping -I <your_current_host_ipsec_loopback_addr> <other_host_ipsec_loopback_addr>
```

## 📝 Tests
I recommend the following tests, for TCP (Single vs Multi Streams) and UDP
I used `iperf3`, and it's what I will be basing my tests upon, you are free to pick any testing framework you need.

### TCP
On one host, bind the address to create a socket for testing:
```bash
iperf3 -s -B <this_hosts_ipsec_loopback_addr>
```

Then in the other host, run the tests from it:
```bash
iperf3 -c <other_hosts_ipsec_loopback_addr> -B <this_hosts_ipsec_loopback_addr> -t <test_duration_time_in_seconds>
```

For TCP multi-stream tests:
```bash
iperf3 -c <other_hosts_ipsec_loopback_addr> -B <this_hosts_ipsec_loopback_addr> -t <test_duration_time_in_seconds> -P <number_of_streams>
```

### UDP
For UDP, setting `-b 0` means send at best effort, or you can select a specific bandwidth like `-b 10M` which means bandwidth of 10Mbps:
```bash
iperf3 -c <other_hosts_ipsec_loopback_addr> -B <this_hosts_ipsec_loopback_addr> -u -b <bandwidth> -t <test_duration_time_in_seconds>
```