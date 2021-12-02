# OCP CoreOS Layering

## Bootstrap PoC

### General Ideas

There are a few ways we can approach this topic. The central questions to this is how we:

1. How much of the config can we build into the layered image
1. Divide responsibilities between ignition and MCD
1. Whether we can do builds in the bootstrap node

We could go so far as to have ignition still provision the full system. The MCD firstboot service would remove all manual ignition changes and pivot directly to the layered image. Generally though I think we should split the responsibilities more cleanly and have ignition only provision things that are required during node provisioning, and have the MCD clean up that subset of configs before pivot, which is what this PoC will attempt.

There is also the general question of what we can do to track the diffs/metadata of configs we own. The MCD will likely have to revert a node to pre-machineconfig state (remove MC changes but keep other manual customizations owned outside the MCO), since this will be required to upgrade from existing configs + perform bootstrap-in-place. In that sense upgrade + boostrap-in-place is the same general problem.

### Steps

1. Have the MCS understand how to serve a "layered image"

Currently we would serve on a different endpoint for PoC purposes

2. Have the MCS render the necessary information to facilitate node bootstrap

This would include: pull secret, desired image, a service to invoke MCD, a service to pull MCD binary, kubeconfig(?), metadata, other non-ignition customizations

3. The service runs between bootstrap and before kubelet runs to consume those changes, namely:

Add non-ignition customizations
Prepare pivot to layered image
Remove manually written files by ignition - maybe store all metadata to /etc/machine-config-daemon, and have that be manually write-able?
Reboots the node

4. The regular MCD takes over. Today we probably need to add some node label to notify MCD to not do anything

### Diffs today

Running an ostree admin config-diff on a 4.9 worker node:

```
M    group
M    gshadow
M    hosts
M    passwd
M    openvswitch
M    openvswitch/default.conf
M    pki/ca-trust/extracted/java/cacerts
M    subgid
M    subuid
M    containers/policy.json
M    containers/registries.conf
M    containers/storage.conf
M    udev/hwdb.bin
M    shadow
A    modules-load.d/iptables.conf
A    openvswitch/.conf.db.~lock~
A    openvswitch/conf.db
A    openvswitch/system-id.conf
A    pki/ca-trust/source/anchors/openshift-config-user-ca-bundle.crt
A    NetworkManager/conf.d/sdn.conf
A    NetworkManager/conf.d/99-keyfiles.conf
A    NetworkManager/dispatcher.d/99-vsphere-disable-tx-udp-tnl
A    NetworkManager/systemConnectionsMerged
A    ssh/sshd_config.d
A    ssh/sshd_config.d/10-disable-ssh-key-dir.conf
A    ssh/ssh_host_ecdsa_key
A    ssh/ssh_host_ecdsa_key.pub
A    ssh/ssh_host_ed25519_key
A    ssh/ssh_host_ed25519_key.pub
A    ssh/ssh_host_rsa_key
A    ssh/ssh_host_rsa_key.pub
A    audit/rules.d/mco-audit-quiet-containers.rules
A    audit/audit.rules
A    sysconfig/network
A    sysconfig/orig_irq_banned_cpus
A    sysctl.d/inotify.conf
A    sysctl.d/forward.conf
A    systemd/system/multi-user.target.wants/kubelet.service
A    systemd/system/multi-user.target.wants/machine-config-daemon-firstboot.service
A    systemd/system/multi-user.target.wants/etc-NetworkManager-systemConnectionsMerged.mount
A    systemd/system/multi-user.target.wants/node-valid-hostname.service
A    systemd/system/multi-user.target.wants/openvswitch.service
A    systemd/system/network-online.target.wants/aws-kubelet-nodename.service
A    systemd/system/network-online.target.wants/ovs-configuration.service
A    systemd/system/kubelet.service.d
A    systemd/system/kubelet.service.d/20-logging.conf
A    systemd/system/kubelet.service.d/10-mco-default-env.conf
A    systemd/system/kubelet.service.d/10-mco-default-madv.conf
A    systemd/system/kubelet.service.d/20-aws-node-name.conf
A    systemd/system/crio.service.d
A    systemd/system/crio.service.d/10-mco-default-env.conf
A    systemd/system/crio.service.d/10-mco-profile-unix-socket.conf
A    systemd/system/crio.service.d/10-mco-default-madv.conf
A    systemd/system/docker.socket.d
A    systemd/system/docker.socket.d/mco-disabled.conf
A    systemd/system/ovs-vswitchd.service.d
A    systemd/system/ovs-vswitchd.service.d/10-ovs-vswitchd-restart.conf
A    systemd/system/ovsdb-server.service.d
A    systemd/system/ovsdb-server.service.d/10-ovsdb-restart.conf
A    systemd/system/pivot.service.d
A    systemd/system/pivot.service.d/10-mco-default-env.conf
A    systemd/system/zincati.service.d
A    systemd/system/zincati.service.d/mco-disabled.conf
A    systemd/system/kubelet.service.requires
A    systemd/system/kubelet.service.requires/kubelet-auto-node-size.service
A    systemd/system/kubelet.service.requires/machine-config-daemon-firstboot.service
A    systemd/system/crio.service.requires
A    systemd/system/crio.service.requires/machine-config-daemon-firstboot.service
A    systemd/system/machine-config-daemon-firstboot.service.requires
A    systemd/system/machine-config-daemon-firstboot.service.requires/machine-config-daemon-pull.service
A    systemd/system/network-online.target.requires
A    systemd/system/network-online.target.requires/node-valid-hostname.service
A    systemd/system/aws-kubelet-nodename.service
A    systemd/system/kubelet-auto-node-size.service
A    systemd/system/kubelet.service
A    systemd/system/machine-config-daemon-firstboot.service
A    systemd/system/machine-config-daemon-pull.service
A    systemd/system/etc-NetworkManager-systemConnectionsMerged.mount
A    systemd/system/node-valid-hostname.service
A    systemd/system/nodeip-configuration.service
A    systemd/system/ovs-configuration.service
A    systemd/system.conf.d
A    systemd/system.conf.d/kubelet-cgroups.conf
A    systemd/system.conf.d/10-default-env-godebug.conf
A    systemd/system-preset
A    crio/crio.conf.d
A    crio/crio.conf.d/00-default
A    tmpfiles.d/cleanup-cni.conf
A    tmpfiles.d/nm.conf
A    iscsi/initiatorname.iscsi
A    kubernetes/cni/net.d/00-multus.conf
A    kubernetes/cni/net.d/multus.d
A    kubernetes/cni/net.d/multus.d/multus.kubeconfig
A    kubernetes/cni/net.d/whereabouts.d
A    kubernetes/cni/net.d/whereabouts.d/whereabouts.kubeconfig
A    kubernetes/cni/net.d/whereabouts.d/whereabouts.conf
A    kubernetes/static-pod-resources
A    kubernetes/static-pod-resources/configmaps
A    kubernetes/static-pod-resources/configmaps/cloud-config
A    kubernetes/static-pod-resources/configmaps/cloud-config/ca-bundle.pem
A    kubernetes/kubelet-plugins
A    kubernetes/kubelet-plugins/volume
A    kubernetes/kubelet-plugins/volume/exec
A    kubernetes/kubelet-plugins/volume/exec/.dummy
A    kubernetes/kubeconfig
A    kubernetes/kubelet.conf
A    kubernetes/kubelet-ca.crt
A    kubernetes/ca.crt
A    kubernetes/cloud.conf
A    kubernetes/manifests
A    docker
A    docker/certs.d
A    docker/certs.d/image-registry.openshift-image-registry.svc:5000
A    docker/certs.d/image-registry.openshift-image-registry.svc:5000/ca.crt
A    docker/certs.d/image-registry.openshift-image-registry.svc.cluster.local:5000
A    docker/certs.d/image-registry.openshift-image-registry.svc.cluster.local:5000/ca.crt
A    mco
A    mco/proxy.env
A    machine-config-daemon
A    machine-config-daemon/node-annotation.json.bak
A    machine-config-daemon/currentconfig
A    .updated
A    gshadow-
A    subuid-
A    shadow-
A    group-
A    ignition-machine-config-encapsulated.json.bak
A    subgid-
A    .pwd.lock
A    mcs-machine-config-content.json
A    node-sizing-enabled.env
A    passwd-
A    machine-id
A    node-sizing.env
A    resolv.conf
```

### Bootstrap node and control plane provisioning

Do we want to run a build in the bootstrap node? Or do we provision control plane nodes differently?

### Bootstrap-in-place

Does that still use onceFrom? Since the bootstrap node likely will be full ignition, do we need to clean things up there? (We can consider this as upgrade place, basically, and not really bootstrap)

### Things to consider

1. How to differentiate requests for layered image provisioning vs what we have today (maybe not relevant in the long run)
1. What is needed to be served to the node initially
1. How do we specify what layered image to be on? Should this be a MCP field
1. What metadata do we need to save
1. How do we update other parts (extensions, kargs) - like we do today?
1. How about FIPS
1. Does https://github.com/openshift/enhancements/blob/master/enhancements/machine-config/custom-ignition-machineconfig.md affect anything
