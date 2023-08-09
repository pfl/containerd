# WIP: containerd with NRI networking resource messages

This is work in progress for containerd with NRI network resource control
messages added. 

## Added messages

When the pod network is started, an 'AdjustPodSandboxNetwork' message is sent
to NRI. This message contains structured information about port mappings,
bandwidth settings, CNI labels and DNS settings, which can be modified and
returned in the reply.

When network configuration is read from disk, all network configurations are
sent to NRI via an 'CreatePodSandboxNetworkConf' message. Eventually returned
networks will be processed and changes applied - this is TBD at the moment.

A version of NRI supporting network messages is in the
[wip_adjustpodsandboxnetwork NRI branch](https://github.com/pfl/nri/tree/wip_adjustpodsandboxnetwork).
Tag 'v0.0.5-adjustpodsandboxnetwork' is pointing to the current work in progress
tip of the branch.

## Compiling

Compiling containerd will pull in the above required dependency.

## Using

A pod/container with the above NRI branch needs to be running for
'AdjustPodSandboxNetwork' message support. See
[work in progress NRI](https://github.com/pfl/nri/tree/wip_adjustpodsandboxnetwork)
for more information.

## Bugs

Probably does not work properly yet.

## Original README

The original [README](/README.md).
