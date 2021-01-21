reference

<https://github.com/WireGuard/wgctrl-go>

<https://godoc.org/golang.zx2c4.com/wireguard/wgctrl>

<https://wiki.archlinux.jp/index.php/WireGuard>

# 環境構成

Peer A 192.168.11.20
Peer B 192.168.11.21

Peer A 側で本プログラムを使用する。

# 手順 

## (1) Peer B

秘密鍵、公開鍵ペアを作成しておく。

```shell
$ wg genkey > privatekey
$ wg pubkey < privatekey > publickey
```

## (2) Peer A

wg0インターフェースの作成と、IPアドレス付与、UPを事前にしておく。

```shell
$ ip link add dev wg0 type wireguard
$ ip addr add 10.0.0.1/24 dev wg0
$ ip link set wg0 up
```

```shell
$ sudo ./wgctrl_go_demo wg0 192.168.11.21 39814 ATaPUw/JCZd9Kn29JJ2ztsraZvqGZ86AMBAw5Tgt5nk= 10.0.0.2/32
```
Peer BのIPアドレス、ポート、公開鍵、許可IPを設定する。

出力例

```
Name: wg0
Type: Linux kernel
Private key: KKoG31L4M1I9707s3Y9wAy9/DqsoCnijD4ZpnVcHQVA=
Public Key: KOSmJlCkdtYuvrnnpO5woy0seCAr/LAVOKzv1SzkTy4=
Listen Port: 48574
Peers: [{ATaPUw/JCZd9Kn29JJ2ztsraZvqGZ86AMBAw5Tgt5nk= AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA= 192.168.11.21:39814 0s 0001-01-01 00:00:00 +0000 UTC 0 0 [{10.0.0.2 ffffffff}] 1}]
```

## (3) Peer B

```shell
$ ip link add dev wg0 type wireguard
$ ip addr add 10.0.0.2/24 dev wg0
$ wg set wg0 listen-port 39814 private-key ./privatekey
$ wg set wg0 peer KOSmJlCkdtYuvrnnpO5woy0seCAr/LAVOKzv1SzkTy4= persistent-keepalive 25 allowed-ips 10.0.0.1/32 endpoint 192.168.11.20:48574
$ ip link set wg0 up
```
