# aws-nginx-ha-manager

[![Build Status](https://travis-ci.org/ccpgames/aws-nginx-ha-manager.svg?branch=master)](https://travis-ci.org/ccpgames/aws-nginx-ha-manager)

This utility is designed to run beside an nginx HA proxy and monitor an AWS ELB and update the list of upstream servers to match the ELB.

The purpose of this utility is to solve some issues with ELBs, such as support for long-polling, proper websocket support and other matters.

The utility queries AWS APIs and ensures the list of ips in the upstream list matches them exactly.

Note: It depends on systemd, i.e. CoreOS, Redhat et al, or newer debian based systems
