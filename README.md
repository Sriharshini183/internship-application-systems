# Submission for Cloudflare Internship Application: Systems

## What is it?

A small Ping CLI application for MacOS or Linux.
The CLI app accepts a hostname or an IP address as its argument, then sends ICMP "echo requests" in a loop to the target while receiving "echo reply" messages.
It reports loss and RTT times for each sent message.
Language chosen : Go

## Useful Links

- [A Tour of Go](https://tour.golang.org/welcome/1)

## Requirements

### 1. Build a tool with a CLI interface

The tool should accept as a positional terminal argument a hostname or IP address.

### 2. Send ICMP "echo requests" in an infinite loop

As long as the program is running it should continue to emit requests with a periodic delay.

### 3. Report loss and RTT times for each message

Packet loss and latency should be reported as each message received.


## Additional feature

1. Added support for both IPv4 and IPv6
