# Subscription Generation

Generate v2ray configuration for V2Ray, Quantumult (X) and Clash.

## Build

```bash
go install github.com/alviezhang/subgen@v0.1.1
```

## Usage

### Start server

```bash
subgen --id <subscription-id> --config <filename> --port <port>
```

### API endpoint

```plain
/sub?id=<subscription-id>&type=<type>
```

`type` is one of `v2ray`, `quantumult` and `clash`.

### Configuration file example

```yaml
client_config:
  uuid: myuuid

nodes:
  - name: node1
    endpoint: node1:10086
    region: JPN
    relays:
      - name: SH-JP
        endpoint: node1-relay1:10086
      - name: SH-HK
        endpoint: node1-relay2:10086

  - name: node2
    endpoint: node2:10086
    region: USA
    relays:
      - name: SZ-HK
        endpoint: node2-relay1:10086
      - name: SH-HK
        endpoint: node2-relay2:10086

  - name: node3
    endpoint: node3:10086
    region: JPN
    relays:
      - name: SH-JP
        endpoint: node3-relay1:10086
      - name: SH-HK
        endpoint: node3-relay2:10086

```
