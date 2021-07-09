# icfpc2021

## 評価機

./eval  (linux)
./eval.darwin (for mac)
./eval.exe (for windows)


### help

```
shiota@DESKTOP-5NQR1JN ~/ICFPC2021/icfpc2021/eval
 % ./eval --help                                                                                                                                                                                                                                                                           ✘ 1 
Usage of ./eval:
  -pose-file string
        pose file
  -problem-file string
        problem file
  -problem-id string
        problem id
  -server string
        http server mode
```

### cli mode

problem-idかproblem fileどちらかを指定
pose-fileが評価したい回答

ファイル指定
```
shiota@DESKTOP-5NQR1JN ~/ICFPC2021/icfpc2021/eval
 % ./eval --problem-file ../problems/1 --pose-file ../solutions/sample/1
3704
```

id指定
```
shiota@DESKTOP-5NQR1JN ~/ICFPC2021/icfpc2021/eval
 % ./eval --problem-id 1 --pose-file ../solutions/sample/1
3704
```

### server mode

server
```
./eval --server 0.0.0.0:8080
```

client
```
> curl -X POST --data "@../solutions/sample/1" localhost:8080/eval/1
{"dislike": 3704}
```