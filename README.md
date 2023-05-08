# nf-svc

This is a service that collects and summarizes network statistics.
It uses the [CloudFlare GoFlow](https://github.com/cloudflare/goflow) project as a library for receiving and pre-processing [NetFlow](https://en.wikipedia.org/wiki/NetFlow) packets.

The main idea is to have a network congestion analysis tool that does not require cloud-like infrastructure or maintenance.

The service is designed to work on a small server or workstation.
Also, although the service can be compiled and run on a multi-platform, the main operating system is Windows.

Because of these specific requirements:
- it has its own log management system
- its settings are set via ini-file

The motivation for creating this tool is to document excessive network load on some device such as a printer, POS, etc. Traffic statistics can be captured using a NetFlow sensor located next to the device. I use a two-port [Mikrotik mAP 2nD](https://mikrotik.com/product/RBmAP2nD) for this, through which the monitored device is turned on.

## How it works

The service receives NetFlow data from the sensor and saves log files of several types:
 - `nf-svc.log` - Information about the operation of service components, an error message, and panic dumps are stored here.
 - `nf-svc-netflow.log` - Here are NetFlow packet dumps, similar to [GoFlow](https://github.com/cloudflare/goflow) log dumps.
 - `nf-svc-summary.log` - This log file is what the service was made for. A summary of traffic statistics collected over a certain time is saved here.

<details><summary>nf-svc-netflow.log example</summary>
  
```
2023-04-29 01:38:54.503 Type:NETFLOW_V9 TimeRecv:1682721534 SequenceNum:172685 Sampler:192.168.255.10 TimeFlowStart:2203946 TimeFlowEnd:2203946 Bytes:110 Packets:2 SrcAddr:192.168.0.82 DstAddr:224.0.0.252 Etype:2048 Proto:17 SrcPort:61238 DstPort:5355 InIf:2 OutIf:1 SrcMac:xx:xx:xx:xx:xx:xx DstMac:aa:aa:aa:aa:aa:aa
2023-04-29 01:38:54.503 Type:NETFLOW_V9 TimeRecv:1682721534 SequenceNum:172685 Sampler:192.168.255.10 TimeFlowStart:2203894 TimeFlowEnd:2203947 Bytes:4212 Packets:54 SrcAddr:192.168.0.82 DstAddr:192.168.0.255 Etype:2048 Proto:17 SrcPort:137 DstPort:137 InIf:2 OutIf:1 SrcMac:xx:xx:xx:xx:xx:xx DstMac:aa:aa:aa:aa:aa:aa
2023-04-29 01:38:54.503 Type:NETFLOW_V9 TimeRecv:1682721534 SequenceNum:172685 Sampler:192.168.255.10 TimeFlowStart:2203947 TimeFlowEnd:2203947 Bytes:370 Packets:8 SrcAddr:192.168.255.80 DstAddr:192.168.255.1 Etype:2048 Proto:6 SrcPort:47141 DstPort:5007 InIf:2 OutIf:1 SrcMac:xx:xx:xx:xx:xx:xx DstMac:cc:cc:cc:cc:cc:cc TCPFlags:2
2023-04-29 01:38:54.503 Type:NETFLOW_V9 TimeRecv:1682721534 SequenceNum:172685 Sampler:192.168.255.10 TimeFlowStart:2203947 TimeFlowEnd:2203947 Bytes:478 Packets:10 SrcAddr:192.168.255.1 DstAddr:192.168.255.80 Etype:2048 Proto:6 SrcPort:5007 DstPort:47141 InIf:1 OutIf:2 SrcMac:xx:xx:xx:xx:xx:xx DstMac:dd:dd:dd:dd:dd:dd TCPFlags:18
2023-04-29 01:38:59.703 Type:NETFLOW_V9 TimeRecv:1682721539 SequenceNum:172688 Sampler:192.168.255.10 TimeFlowStart:2203945 TimeFlowEnd:2203951 Bytes:495 Packets:3 SrcAddr:192.168.0.20 DstAddr:239.255.255.250 Etype:2048 Proto:17 SrcPort:63710 DstPort:1900 InIf:2 OutIf:1 SrcMac:xx:xx:xx:xx:xx:xx DstMac:ee:ee:ee:ee:ee:ee
```

</details>

<details><summary>nf-svc-summary.log example</summary>
  
```
2023-04-29 02:22:14.207 *** Summary for every 5 minutes ***
NETFLOW_V9(192.168.255.10)
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=224.0.0.251, DstPort=5353, {Bytes:37321 Packets:611}
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=192.168.0.255, DstPort=137, {Bytes:18252 Packets:234}
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=224.0.0.252, DstPort=5355, {Bytes:13420 Packets:244}
  L3=IPv4, L4=UDP, Src=192.168.0.128, Dst=239.255.255.250, DstPort=1900, {Bytes:4752 Packets:24}
  L3=IPv4, L4=TCP, Src=192.168.255.1, Dst=192.168.255.80, SrcPort=5007, {Bytes:3095 Packets:63}
  L3=IPv4, L4=TCP, Src=192.168.255.80, Dst=192.168.255.1, DstPort=5007, {Bytes:2603 Packets:57}
  L3=IPv4, L4=UDP, Src=192.168.0.123, Dst=239.255.255.250, DstPort=1900, {Bytes:2340 Packets:12}
...
2023-04-29 02:37:14.222 *** Summary for every 20 minutes ***
NETFLOW_V9(192.168.255.10)
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=224.0.0.251, DstPort=5353, {Bytes:144461 Packets:2365}
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=192.168.0.255, DstPort=137, {Bytes:73944 Packets:948}
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=224.0.0.252, DstPort=5355, {Bytes:59840 Packets:1088}
  L3=IPv4, L4=UDP, Src=192.168.0.128, Dst=239.255.255.250, DstPort=1900, {Bytes:22196 Packets:112}
  L3=IPv4, L4=TCP, Src=192.168.255.1, Dst=192.168.255.80, SrcPort=5007, {Bytes:14672 Packets:299}
  L3=IPv4, L4=TCP, Src=192.168.255.80, Dst=192.168.255.1, DstPort=5007, {Bytes:12184 Packets:267}
  L3=IPv4, L4=UDP, Src=192.168.0.24, Dst=239.255.255.250, DstPort=1900, {Bytes:8120 Packets:40}
...
2023-04-29 10:17:14.233 *** Summary for every 480 minutes ***
NETFLOW_V9(192.168.255.10) top 100 of 146
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=224.0.0.251, DstPort=5353, {Bytes:3606097 Packets:59037}
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=192.168.0.255, DstPort=137, {Bytes:1770210 Packets:22695}
  L3=IPv4, L4=UDP, Src=192.168.0.82, Dst=224.0.0.252, DstPort=5355, {Bytes:1500893 Packets:27289}
  L3=IPv4, L4=UDP, Src=192.168.0.128, Dst=239.255.255.250, DstPort=1900, {Bytes:568676 Packets:2872}
  L3=IPv4, L4=UDP, Src=192.168.0.24, Dst=239.255.255.250, DstPort=1900, {Bytes:194880 Packets:960}
  L3=IPv4, L4=UDP, Src=192.168.0.122, Dst=239.255.255.250, DstPort=1900, {Bytes:189885 Packets:977}
...
```

</details>
