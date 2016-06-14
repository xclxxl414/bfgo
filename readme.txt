QQ交流群340195342，点击加入：http://jq.qq.com/?_wv=1027&k=2ADNTk3
======

快速上手
======
1. 下载bftrader发布包
下载地址: https://github.com/sunwangme/bftrader/releases
下载地址: http://pan.baidu.com/s/1nvgrNst

2. 安装golang编译器和IDE
   2.1 安装 golang1.6.2 windows x86
   2.2 安装 liteide x29 windows x86
   2.3 安装git for windows
   
3. 下载bygo源代码
   3.1 go get github.com/sunwangme/bfgo

4. 写策略，调试策略  
   4.1 运行ctpgateway.exe,datafeed.exe
   4.2 点击ctpgateway的net/netStart,点击datafeed的net/netStart
   4.3 运行datarecorder/main.go，以连接ctpgateway datafeed
   4.4 点击ctpgateway的ctp/ctpStart
   4.5 可以看到datarecorder跑起来啦

网友策略列表
======
darecorder/：tick收集器，演示BfTraderClient+BfRun的用法
kvclient/ & kvserver/：tick收集器，演示BfTraderClient+BfRun的用法

（完）
