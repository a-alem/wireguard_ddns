# 📦 WireGuard DDNS
Affordable WireGuard tunnels made easy.

## 🌟 Highlights
- Establish seamless ddns resolution services
- Extensible and modular design written in Go
- Non-negotiable small package size, and limited number of 3rd party packages used by design.
- Services are created with systemd and systemd timers
- Wireguard configuration ready template
- Tagged releases controlled by protected branches and CI/CD deployments

## ❗️ What this project IS NOT
This project does not configure wireguard on your behalf, it assumes that you already have a proper wireguard interfaces configured, and the peers have already generated their respective keypairs and have exchanged their public keys through a proper PKI channel.

## ℹ️ Overview
This project aims to mitigate the complexities of establishing affordable WireGuard tunnels using dynamic DNS (DDNS) to avoid static public IP costs with minimal downtime due to DNS resolution and propagation delays.

The structure of the project is compartmentalized into the following, with each directory having a separate internal readme file for further documentation:
- `/ddns_updated`: The core DDNS resolution service written in Go
- `/scripts`: Ready made scripts for systemd service, systemd timer, and the installations script that abstracts away all the complexities of configuring the service
- `/terraform`: A ready template to test the project in a plug and play manner in AWS by creating the necessary VPCs, EC2 instances and security groups.
- `/wireguard`: a pseudo template for WireGuard's configurations

WireGuard configuration is not heavily involved or restricted here by design, this project acts as an abstraction layer to any tunneling technology that can be configured with DNS, so an interested patry, in theory, can use this project to establish an IPSec tunnel by resolving domain names.

The reason WireGuard was chosen is because of its brilliant simplicity and opinionated design, most crucially, its fail-safe design and the assurance that if the tunnels are misconfigured in some way, a connection is impossible to maintain, unlike some other VPN technologies such as IPSec where the tunnel could still operate under false security pretentions. This fact perfectly matches the kind of traffic intended for this project

## 🚀 Usage
### Prerequisites
Install `openresolv` package or similar alternatives in order for WireGuard to resolve domain names

For example, on ubuntu:
```bash
sudo apt install openresolv
```
### ⬇️ Installation
```bash
git clone <repo>

sudo setup_ddns_updater.sh \ 
   <release_version_number> \
   <config.yaml> \
   <ddns-updater.service> \
   <ddns-updater.timer>
```
## 💭 Feedback and Contributing
All contributions are welcomed!
Please feel free to read the contribution document 

Feel free to also reach out and open discussions/raise issues anytime

For issue reporting, please use GitHub's issues from the issues tab
Similarly, to open discussions, please use GitHub's discussions from the discussions tab