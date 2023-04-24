# nf-svc

This is a service that collects and summarizes network statistics.
It uses the [CloudFlare GoFlow](https://github.com/cloudflare/goflow) project as a library for receiving and pre-processing [NetFlow](https://en.wikipedia.org/wiki/NetFlow) packets.

The main idea is to have a network congestion analysis tool that does not require any infrastructure or maintenance.

The service is designed to work on a small server or workstation.
Also, although the service can be compiled and run on a multi-platform, the main operating system is Windows.

Because of these specific requirements:
- it has its own log management system
- its settings are set via ini-file

