# aws-nginx-ha-manager

This utility is designed to run beside an nginx HA proxy and monitor an AWS ELB and update the list of upstream servers to match the ELB.

The purpose of this utility is to solve some issues with ELBs, such as support for long-polling, proper websocket support and other matters.

The utility simply resolves the ELBs dns and ensures the list of ips in the upstream list matches them exactly.
