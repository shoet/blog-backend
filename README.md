# blog

## リリース

ServerlessFramework でデプロイする。

```
make deploy
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
cd_tools
sql-migrate up [-dryrun]
```

- マイグレーションの巻き戻し

```
cd_tools
sql-migrate down [-limit=n]
```
