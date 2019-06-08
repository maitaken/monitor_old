# Monitor
![demo](https://raw.github.com/wiki/maitaken/monitor/monitor.gif)

## 使い方

```sh
monitor [監視するファイル名] ["実行するシェル"]
ex)
monitor main.py "python3 test.py"
```

## オプション
| オプション |   引数   |                        概要                       |
|:----------:|:--------:|:-------------------------------------------------:|
|      f     | filename | 指定したファイルを監視 (複数ファイル監視時に使用) |

### f
`-f value`で使用．
```
ex)
monitor -f in -f main.py "python main.py < in"
```

## TODO
* Windows対応
* ワイルドカードでのファイル指定
* 詳細(実行時間など)を表示するオプションの追加
* 出力をリアルタイムに表示するオプションの追加
