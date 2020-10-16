=====================================
 かけるさん 〜負荷を掛けるよ2020秋〜
=====================================

.. contents::

アプリケーションのパッケージング
--------------------------------

以下で artifacts バケットにアプリケーションが追加される。

::

   $ make upload


AMI焼き
-------

../provisioning/ami/kakeru.json の部分を環境に合わせて変更する。
その後以下のコマンドを実行するとAMIが焼かれる。

::

   $ cd ../provisioning/ami
   $ make TO=kakeru


実行
======================

起動テンプレート
~~~~~~~~~~~~~~~~~~~~~~

起動テンプレート名 ``kakeru-`` で始まるものを用意してある
これを変更して、AMIを先程焼いたkakeruのAMIを指定して新しいバージョンを作成する

Auto Scaling グループ
~~~~~~~~~~~~~~~~~~~~~~

すでに ``kakeru`` という名前で Auto Scaling グループを用意している

Auto Scaling グループからインスタンスを起動
~~~~~~~~~~~~~~~~~~~~~~

希望する容量を ``希望する台数+1`` にしてインスタンスを立てる。
複数立てたうちの一台はコントローラーノードなので、ワーカーも含めてマルチノード処理したい場合は **最低3台** 必要。

負荷を掛ける
------------

立てたインスタンスに ssh して、以下のようにコマンドを実行すると hakaru に負荷がかかる。
シナリオは 1, 2, 3 のうちのどれかが使える。

::

   インスタンスへの接続 -> https://github.com/voyagegroup/sunrise2020/blob/master/docs/session.md
   # cd /opt/kakeru
   # make deploy # artifactsからアプリケーション配備するやつ
   # make -C app kakeru upload HOST=${ELBエンドポイントのドメイン} SCENARIO=${1,2,3}

実行結果は {チーム名}-kakeru-report のバケットにアップロードされるので、見てみましょう。

アプリケーションの再デプロイ
============================

再パッケ

::

   $ make clean upload

インスタンスへのデプロイ

::

   インスタンスへの接続 -> https://github.com/voyagegroup/sunrise2020/blob/master/docs/session.md
   # cd /opt/kakeru
   # make deploy

kakeru を実行するインスタンスだけでデプロイすればOK
