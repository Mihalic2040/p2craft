# P2Craft

P2Craft lets a Minecraft Java client join a server through a [libp2p](https://libp2p.io) connection instead of a public IP and port forwarding. The server side is a small Go daemon (`hole`) that sits next to a Minecraft server. The client side is a Fabric mod with an embedded Go library. Players enter `<peer-id>.p2p` as the server address.

## Features

Everything below is implemented in `core/` 

- **Peer addressing by libp2p peer ID.** Servers use a persistent Ed25519 identity. Clients generate a new ephemeral identity on every launch.
- **Transports:** TCP and QUIC, secured with Noise.
- **Peer discovery:**
  - Kademlia DHT, using libp2p's default bootstrap peers. Clients run in DHT client mode and servers in `ModeAuto`.
  - mDNS on the local network, with service tag `p2craft`.
- **NAT traversal:**
  - UPnP/NAT-PMP port mapping
  - AutoNAT v1 / v2
  - Circuit relay v2 client with AutoRelay (1 relay wanted, at least 2 candidates)
  - Hole punching DCUtR 
- **TCP tunnelling.** Each local TCP connection becomes one libp2p stream on protocol `/p2craft/minecraft/1.0.0`.
- **Client integration.** Fabric mixins intercept both the multiplayer server-list ping and the join action for `*.p2p` addresses.
- **Native platforms bundled in the mod JAR:** Linux x86-64, Linux aarch64, Windows x86-64.

## How it works

```
 Minecraft client (Fabric + P2Craft mod)                      Server host
┌──────────────────────────────────────┐                ┌────────────────────────────┐
│ Minecraft ──TCP──► 127.0.0.1:<port>  │                │ hole (Go daemon)           │
│ (mixins rewrite *.p2p addresses)     │                │  libp2p host               │
│        ▲                             │                │  stream handler ──TCP──►   │
│        │ JNA                         │   libp2p       │      Minecraft server      │
│  libp2pcraft native lib (Go) ◄───────┼── stream ──────┼► 0.0.0.0:25565             │
│  TCPProxy per remote peer            │ direct / holepunched / via relay            │
└──────────────────────────────────────┘                └────────────────────────────┘
```

1. A player enters `<peer-id>.p2p` when joining a Minecraft server.
2. The mod detects the `.p2p` address and starts setting up a connection in the background.
3. The Go library creates a local TCP proxy on a free port. If one already exists for that server, it reuses it.
4. Minecraft connects to that local address instead. Server list pings are also sent through the same connection.
5. The proxy connects to the other player's peer using libp2p. If it doesn't know where that peer is, it uses the DHT to find it.
6. On the other end, the `hole` program connects to the actual Minecraft server and passes data back and forth.

The connection between players is encrypted by libp2p. The connection between Minecraft and the local proxy uses regular localhost TCP.

### Connectivity

Both nodes try to become reachable: port mapping, AutoNAT reachability checks, relay reservations (AutoRelay), and hole punching. The relay candidates come from a custom peer source (`core/node/peer_source_relay.go`) that periodically feeds in DHT-closest peers and currently connected peers. The feeder first fires after 15 s and then backs off ×3 each round, capped at 1 h.

P2Craft nodes do not run a relay service themselves (`EnableRelayService` is not used). They use relays offered by other libp2p peers. The code does not opt in to data over limited relay connections (no `network.WithAllowLimitedConn` in `core/`), so treat relays as a means of coordinating hole punching. Whether a given NAT combination ends up with a direct connection is not guaranteed.

## Architecture

Native API exported from `core/lib/main.go`:

| Function | Behaviour |
|---|---|
| `Connect()` | Starts a client-mode node (ephemeral identity). Returns `0` on success and `-1` on error. Does nothing if already connected. |
| `GetStatus()` | `0` disconnected, `1` connected, `2` error |
| `CreateTunnel(remoteID)` | Starts or reuses a loopback proxy for the peer; returns `"127.0.0.1:<port>"`, or `NULL` on failure |
| `FreeString(ptr)` | Frees a string returned by the library |
| `Disconnect()` | Closes the tunnels and resets the status |

The mod calls `Connect()` once at client initialisation (`P2craftClient`). The main (server-side) mod entrypoint only logs a message; hosting is done by `hole`, not by the mod.
Maybe maybe, some day. There will be server side mod. 

## Requirements

### To play
- Minecraft 1.20.6
- Fabric Loader 0.20.1 or newer
- Fabric API 0.100.8+1.20.6

### To build from source
- Go 1.27.2
- JDK 21 or newer
- Gradle 9.8.2 (included in the project)
- C compiler and Make:
  - `gcc`
  - `aarch64-linux-gnu-gcc`
  - `x86_64-w64-mingw32-gcc`
  - `make`

On Debian/Ubuntu the equivalents are likely `build-essential gcc-aarch64-linux-gnu gcc-mingw-w64-x86-64`. 

## Building from source


### Full build (all three native libraries plus the mod) *CLIENT*

```sh
cd p2craft
./gradlew build
```

`processResources` depends on the `buildNative` task, which runs `make install` in `../core`. That builds `libp2craft.so` (linux-x86-64 and linux-aarch64) and `p2craft.dll` (windows-x86-64) with `go build -buildmode=c-shared`, and installs them into `p2craft/src/main/resources/natives/<platform>/`. That directory is git-ignored. All three cross-compilers must therefore be installed.

The JAR is written to `p2craft/build/libs/p2craft-*.jar`.

### Building natives separately

```sh
cd core
make                 # same as `make install`: all targets, copied into the mod resources
make linux-amd64     # one target only, output in core/build/linux-x86-64/
make clean
```

Other targets: `linux-arm64`, `windows-amd64`, `build`. 


### Server daemon (`hole`)

```sh
cd core
go build -o build/hole ./cmd/hole
```

## Usage

### Hosting

1. Start your Minecraft server so it listens on **port 25565** on the same machine.
2. Start `hole`:

   ```sh
   cd core
   HOLE_PRIV_KEY=../data/private.key HOLE_TARGET_ADDR=127.0.0.1:25565 ./build/hole
   ```

   The key file is created on first run (mode 0600, parent directory 0700). Keep it: it defines your peer ID.
3. Read the log line `Hole ID: <peer-id>.p2p` and share that value to the players.


| Variable | Used by | Meaning |
|---|---|---|
| `HOLE_PRIV_KEY` | `hole` | libp2p private-key file; created if missing.|
| `HOLE_TARGET_ADDR` | `hole` | Deafult is 127.0.0.1:25565 |


### Joining

The mod has no UI of its own. So i think you know how to add a server? O_o


### Docker (server)

`compose.yml` builds `core/Dockerfile` into a `hole` image, mounts `./data` as `/app`, and sets `HOLE_PRIV_KEY=/app/private.key`:

```sh
docker compose up -d --build
docker compose logs hole      # shows the Hole ID
```

## Troubleshooting

Good luck. 

## Known problems

- `Disconnect()` does not close the libp2p host, and the mod never calls it.

## License

Copyright (c) 2026 All rights reserved
Mihalic2040.
