# hakaru [![Build Status][travis-img]][travis-url]

[travis-img]: https://travis-ci.com/voyagegroup/hakaru.svg?token=iBCGFnZyWWvHWvMJnnx3&branch=master
[travis-url]: https://travis-ci.com/voyagegroup/hakaru

Sunrise2019: 素朴な計測サーバ

## manually initial setup checklist

- [ ] voyagegroup/sunrise2019 での `make apply` が完了している
- [ ] ./team_name.txt の1行目をチーム名に変更している
- [ ] ./provisioning/ami/packer.json を voyagegroup/sunrise2019 hakaru/README.md の通りに変更している
- [ ] ./provisioning/instance/sysconfig/hakaru を voyagegroup/sunrise2019 hakaru/README.md の通りに変更している

## build AMI

```bash
$ cd provisioning/ami
$ make
```

### launch EC2 instance

voyagegroup/sunrise2019 hakaru/README.md を参考にしてください

## deployment

1. ビルドを実施し、成果物をアップロードする

*ビルド/アップロードを自動化する場合は .travis.yml を参考に*

```bash
$ make install
$ make upload
```

1. blue/green or in-place のどちらかを実施する

### blue/green deployment

1. AMI をビルドする
1. AMIからEC2インスタンスを起動する
1. 起動するEC2インスタンスの User data に ./user_data.sh の内容をコピペする
1. EC2インスタンスをロードバランサーに紐付る
1. 古いEC2インスタンスを終了する

### in-place deployment

1. 既にEC2インスタンスを起動していること
1. インスタンス上でコマンドを実行する

```bash
$ sudo su -l
# cd /root/hakaru
# make ARTIFACTS_COMMIT=latest
```
