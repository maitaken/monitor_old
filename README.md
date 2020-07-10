# Monitor
![demo](https://raw.github.com/wiki/maitaken/monitor/monitor.gif)

# Install
```sh
make install
```

## 使い方

```sh
monitor [監視するファイル名] ["実行するシェル"]
ex)
monitor main.py "python3 test.py"
```

## オプション
| オプション |   引数   |                        概要                       |
|:----------:|:--------:|:-------------------------------------------------:|
|      f     | string | 指定したファイルを監視 (複数ファイル監視時に使用) |
|      s     | - | ターミナルに収まるように出力を省略 |
|      t     | int | 実行のタイムアウトを設定(秒) |

### f
`-f string`で使用．
```
ex)
monitor -f in -f main.py "python main.py < in"
```
### s
`-s`で使用．
```
ex)
monitor -f in -f main.py -s "python main.py < in"
```

### t
`-t int`で使用
```
ex)
monitor -t 2 -f in -f main.py -s "python main.py < in"
```

## TODO
* Windows対応
* ワイルドカードでのファイル指定
* 詳細(実行時間など)を表示するオプションの追加
* 出力をリアルタイムに表示するオプションの追加
