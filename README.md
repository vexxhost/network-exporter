# Network Exporter

Prometheus exporter for a variety of network devices using their respective
APIs.  The exporter implements a generic standarized way of delivering
information regardless of the device used, which makes it useful for
environments with heterogeous network infrastructure.

## Supported Devices

- Arista
- MikroTik

## Supported Collectors

- BGP
- Device Info
- Interface statistics

## SNMP

There is an extremely useful and powerful SNMP exporter which already exists
however we've found that we could never consistently retrieve information
across several device manufactureres (for example, some do not expose BGP
peer information over SNMP).

The library embedded in here does the majority of the work of generalizing
that information retrieval so you can have very clear and clean metrics
regardless of the devices used in your environment.
