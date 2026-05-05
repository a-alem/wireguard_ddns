# Performance collection between WireGuard and IPSec
The goal is to establish a baseline of unified testing that measures the tunnel's activation point in time (delay) for both WireGuard and IPSec, and performs empirical comparisons between both delays.

One benefit of doing so is to contextualize what kind of traffic each of these VPN technologies is best suited for, does a longer tunnel establishment time delay that guarantees better traffic condition is more suitable than a shorter tunnel establishment time delay but for the cost of link's traffic stability? 

By testing both, we can reach to that conclusion in a theoretical way, it is nigh impossible to guarantee result consistency across the internet or unknown intranet backbones, so these results are to be taken with a grain of salt, it shows a general outline or trajectory, not a guaranteed axiom of truth.

## Cold start measurement
The objective here is to test both VPN protocols at their initial phase of connection up to the point in time of the successful receival of the first encapsulated(encrypted) packet in the tunnel.

The benchmarking is used to compare both VPN protocols of their unique establishing properties, and how fast they converge into a working tunnel from ground zero.

One crucial element is the elimination of state, the slate must be clean and no state of either protocol should be present in order to make sure that the comparison is fair, so in the start of each test, all configurations and interfaces are removed, then re-implemented.

The process is to select a peer as the testing ground (peer A), and the other as the awaiting peer (peer B), this peer's state must be manually removed by taking down the interfaces and configurations, then initialize it again manually.
With the state removed from peer B, this peer would be listening to any incoming attempts of tunnel establishment from peer A (testing peer)

### WireGuard
Place the following in a file after interpolating the template values with real ones, and execute it
```bash
#!/usr/bin/env bash
set -euo pipefail

cleanup() {
  echo
  echo "Interrupted. Exiting WireGuard cold ping test..."
  exit 130
}

trap cleanup INT TERM

sudo wg-quick down <path_to_wg_conf_file> 2>/dev/null || true
sleep 2

START=$(date +%s%3N)

sudo wg-quick up <path_to_wg_conf_file>

until ping -c 1 -W 1 <other_peer_wireguard_private_ip> >/dev/null 2>&1; do
  :
done

END=$(date +%s%3N)

echo "WireGuard cold ping establishment delay: $((END-START)) ms"
```

### IPSec
Note, make sure that the security association files (or contents) are copied to `/etc/swanctl/swanctl.conf` in each peer respectively prior to running the testing script.
Place the following in a file after interpolating the template values with real ones, and execute it
```bash
#!/usr/bin/env bash
set -euo pipefail

cleanup() {
  echo
  echo "Interrupted. Exiting WireGuard cold ping test..."
  exit 130
}

trap cleanup INT TERM

sudo swanctl --terminate --ike ipsec-bench 2>/dev/null || true
sleep 2

START=$(date +%s%3N)

sudo swanctl --initiate --child bench

until ping -I <local_loopback_ipsec_bind_address> -c 1 -W 1 <other_peer_loopback_ipsec_bind_address> >/dev/null 2>&1; do
  :
done

END=$(date +%s%3N)

echo "IPSec cold ping establishment delay: $((END-START)) ms"
```

## New Installation - TCP session establishment
This measures the time delay from a new installation with no state up to the point of TCP session establishment, so we can think of it as: tunnel establishment delay + tcp session establishment delay.

### WireGuard
In the receiver node, open an iperf3 session
```bash
iperf3 -s -1 -B <current_node_wg_local_ip>
```

Then in the testing(sender) node
```bash
#!/usr/bin/env bash
set -euo pipefail

cleanup() {
  echo
  echo "Interrupted. Exiting WireGuard cold ping test..."
  exit 130
}

trap cleanup INT TERM

sudo wg-quick down <path_to_wg_conf_file> 2>/dev/null || true
sleep 2

START=$(date +%s%3N)

sudo wg-quick up <path_to_wg_conf_file>

until iperf3 -c <other_peer_wireguard_private_ip> -B <current_peer_wireguard_private_ip> -t 1 >/dev/null 2>&1; do
  :
done

END=$(date +%s%3N)

echo "WireGuard new installation - iperf3 TCP session establishment delay: $((END-START)) ms"
```

### IPSec
In the receiver node, open an iperf3 session
```bash
iperf3 -s -1 -B <current_node_ipsec_local_loopback_ip>
```

Then in the sender node
```bash
#!/usr/bin/env bash
set -euo pipefail

cleanup() {
  echo
  echo "Interrupted. Exiting WireGuard cold ping test..."
  exit 130
}

trap cleanup INT TERM

sudo swanctl --terminate --ike ipsec-bench 2>/dev/null || true
sleep 2

START=$(date +%s%3N)

sudo swanctl --initiate --child bench

until iperf3 -c <other_peer_loopback_ipsec_bind_address> -B <local_loopback_ipsec_bind_address> -t 1 >/dev/null 2>&1; do
  :
done

END=$(date +%s%3N)

echo "IPSec cold iperf3 establishment delay: $((END-START)) ms"
```