# blog

## リリース

AWS CDK でデプロイする。

```
make deploy-dev
```

## DB のマイグレーション

sql-migrate にて実施する。

- マイグレーションファイルの作成

```
cd _tools
sql-migrate new <ファイル名>
```

- マイグレーションの実施

```
cd _tools
sql-migrate up -config=dbconfig.yml -env=<env> [-dryrun]
```

- マイグレーションの巻き戻し

```
cd _tools
sql-migrate down -config=dbconfig.yml -env=<env> [-limit=n]
```
