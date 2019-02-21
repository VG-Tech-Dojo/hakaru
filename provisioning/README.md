# provisioning

- ami/ ではAMIを作る
  - packer を利用している
- instance/ にはアプリケーションのインストール/セットアップ処理を置く
  - hakaru
  - amazon-cloudwatch-agent

## on Instance

```
$ cd /root/hakaru
$ tree .
.
├── Makefile # ami/scripts/deploy/Makefile が配置される。 artifacts.tgz を持って来てデプロイを実行する
└── app # artifacts.tgz の展開先。instance/配下とほぼ同じ
    ├── Makefile # インストール/セットアップ処理
    ├── amazon-cloudwatch-agent
    │   └── amazon-cloudwatch-agent.json
    ├── hakaru # hakaruバイナリ
    ├── sysconfig
    │   └── hakaru
    └── systemd
        └── hakaru.service
```
