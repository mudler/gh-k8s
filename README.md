# k8s cluster on gh runners


Ever wondered if it's possible to have disposable multi-node Kubernetes cluster from GitHub Action runners? Yes it is possible!

## What

This repository is just a template that you can fork on. There is a GitHub pipeline configured to run a multi-node cluster on GitHub action runners with [k3s](https://k3s.io/). The Cluster after bootstrap is accessible via a unique VPN: it's a complete access, not just few ports exposed (like by using e.g. ngrok).
For convenience there is also a "jumpbox" job which is waiting for a tmate session. It can login only the event author by default.

You can tweak several things in the workflow file (which is commented): disable fleet installation, adding more worker, etc.

## Why?

![](https://media.giphy.com/media/2xwWl4iiaR0UKIiiRQ/source.gif)

Because Why not? Seriously, the reason is simple: Testing and Developing. If you are an Engineer working in the cloud or either a devops, you name it how many times you find yourself needing a spare k8s cluster for testing or proofing a concept. I hear you, and that's why this repository exists.

## Usage

### 1) Fork this repository

Fork and star (if you like it) this repository on Github! with a copy of this repo you will get all the necessary pieces to run the k8s cluster on Github Action runners.

### 2) Create a VPN configuration

For accessing the cluster we will use a VPN - no worries, no need to setup any server, or need of any box for routing traffic. We will use [edgevpn](https://github.com/mudler/edgevpn) as it is decentralized, and runs behind NAT. It is also static so it's handy to install it locally.

Now let's create our vpn configuration file, we will encode it in base64 and add it as a repository secrets so the node can connect between each other:

```bash
edgevpn -g | tee config.yaml | base64 -w0 > config_encoded
```

now set the content of `config_encoded` to the repository as `EDGEVPN` secret.

### 3) Run the workflow

Commit something, it will run automatically. Or just check out the latest run and click on the "Re-run job" button.


That's it, now go into the "jumpbox" job and wait for it to display an ssh connection string.

Once inside, you can use kubectl with:

```bash
KUBECONFIG=k3s.yaml kubectl get pods -A
```

## More

### Persistency

The cluster is not persistent. Github Action has a limit timeout of 6h. Nothing forbids to reschedule another job each 6h our, and you can decomment that in the workflow file.

### Access the cluster

To be able to access the cluster from your system, you need to connect via VPN.

In the terminal, run:

```bash
sudo IFACE=edgevpn0 ADDRESS=10.1.0.2/24 EDGEVPNCONFIG=$PWD/config.yaml edgevpn
```

_Note_ that the `ADDRESS` we are setting here is the one we will have in the VPN. We are not setting any public IP. Addresses are internal.

Open another terminal, and wait for connection to be available, monitor ```ping 10.1.0.20```.

The setup is fixed, you will find the master node on the `10.1.0.20` ip over VPN.

### Grab cluster kube config

From inside the VPN you can always get the latest kubeconfig with:

```bash
curl http://10.1.0.20:9091/k3s.yaml | sed 's/127\.0\.0\.1/10.1.0.20/g' > k3s.yaml
```

### `manifests`

Each kubernetes resource in the manifests folder gets automatically deployed with fleet. You can disable this behavior by commenting the `install_helm` and `install_fleet` line in the cluster workflow file.