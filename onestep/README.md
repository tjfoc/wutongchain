这里将会提供一键式命令来启动两个节点和一个sdk的网关服务，可供最基本的区块链调用。当使用启动命令启动成功后，可在浏览器里查看是否启动成功http://<物理机的ip>:8888/getheight

使用说明:
启动命令:./start.sh
停止命令:./stop.sh

注意事项:
1.启动之前需要先安装docker和docker-compose
2.每个节点使用的都是物理机的网卡，所以peer，sdk使用的ip都是和物理机一样的ip(参看如下集群配置，不可使用localhost)，修改配置文件时请注意ip地址和端口，避免端口的重复监听
3.默认是使用守护进程模式启动，该模式下不会有日志输出，peer的日志会输出到对应文件夹下的logs中，如果不想以守护进程模式启动，请修改start.sh，将最后的-d去掉

如下：
Self:
    Id: "1"
    ShownName: "peer1"
    Addr: "192.168.x.xxx:60000"

Members:
    Peers:
        - Id: "1"
          ShownName: "peer1"
          InAddr: "192.168.x.xxx:60000"
          OutAddr: "192.168.x.xxx:60000"

        - Id: "2"
          ShownName: "peer2"
          InAddr: "192.168.x.yyy:60000"
          OutAddr: "192.168.x.yyy:60000"
        