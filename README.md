# aws-nginx-ha-manager (0.1.4)

[![Build Status](https://travis-ci.org/ccpgames/aws-nginx-ha-manager.svg?branch=master)](https://travis-ci.org/ccpgames/aws-nginx-ha-manager)

This utility is designed to run beside an nginx HA proxy and monitor an AWS ELB and update the list of upstream servers to match the ELB.

The purpose of this utility is to solve some issues with ELBs, such as support for long-polling, proper websocket support and other matters.

The utility queries AWS APIs and ensures the list of ips in the upstream list matches them exactly.

Note: It depends on systemd, i.e. CoreOS, Redhat et al, or newer debian based systems

## Example usage

To monitor an ELB named CertVKPil-ServiceL-15ERWL9O3YIZI issue the following command:

```bash
$ aws-nginx-ha-manager monitor CertVKPil-ServiceL-15ERWL9O3YIZI --upstream-file /etc/nginx/conf.d/pilot.upstream.conf --upstream-name pilot --port 8000 --interval 5
```

This will query the AWS API every ```5``` seconds and write an upstream file to ```/etc/nginx/conf.d/pilot.upstream.conf``` with upstream name ```pilot``` calling on port ```8000```

You can also type ```aws-nginx-ha-manager --help``` for a list of options and commands.

## Access requirements

The command requires access to the AWS ELB API, specifically the action ```elasticloadbalancing:DescribeListeners```.

Example policy below:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Stmt1473172616000",
            "Effect": "Allow",
            "Action": [
                "elasticloadbalancing:DescribeListeners"
            ],
            "Resource": [
                "arn:aws:elasticloadbalancing:eu-west-1:577470089438:loadbalancer/CertVKPil-ServiceL-15ERWL9O3YIZI"
            ]
        }
    ]
}
```
