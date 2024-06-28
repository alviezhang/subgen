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
    host: node1
    port: 10086
    region: JPN

  - name: node2
    host: node2
    port: 10086
    region: JPN
```
