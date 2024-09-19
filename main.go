package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/zhnt/connector/utils"
)

func main() {

	// 定义运行上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	/**
	响应节点，即发送节点不需要添加此参数
	请求节点通过该参数访问定义需要访问的响应节点的地址
	*/
	dest := flag.String("d", "", "Destination multiaddr string")
	flag.Parse()

	// 通过标识对身份进行判断
	if *dest == "" {
		// 是响应节点 需要添加协议处理
		host, _ := utils.CreateHostWithPort(7001)
		defer host.Close()

		host.SetStreamHandler(utils.PROTOCAL_CONNECTOR, handleStream)
		<-ctx.Done()
	} else {
		// 是请求节点
		host, _ := utils.CreateHostWithPort(0)
		defer host.Close()
		stream, err := connect(host, *dest)
		if err != nil {
			fmt.Println("Connect failed", err)
			return
		}
		testStream(stream)
		//err := downloadFile(ctx, h, *dest, fileName)
		//if err != nil {
		//return
		//}
		return
	}

	// 阻塞等待
	<-make(chan struct{})
}

func handleStream(s network.Stream) {
	// 打印远程对等节点信息
	fmt.Println("Got a new stream from", s.Conn().RemotePeer())

	// 创建读写器
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// 写入消息
	_, err := rw.WriteString("你好，世界！\n")
	if err != nil {
		fmt.Println("写入消息失败:", err)
		return
	}
	err = rw.Flush()
	if err != nil {
		fmt.Println("刷新缓冲区失败:", err)
		return
	}

	// 读取消息
	message, err := rw.ReadString('\n')
	if err != nil {
		fmt.Println("读取消息失败:", err)
		return
	}

	// 打印接收到的消息
	fmt.Printf("收到消息: %s", message)

	s.Close()
}

func testStream(s network.Stream) {

	// 创建读写器
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	// 写入消息
	_, err := rw.WriteString("你好，服务器！\n")
	if err != nil {
		fmt.Println("写入消息失败:", err)
		return
	}
	err = rw.Flush()
	if err != nil {
		fmt.Println("刷新缓冲区失败:", err)
		return
	}

	// 读取消息
	message, err := rw.ReadString('\n')
	if err != nil {
		fmt.Println("读取消息失败:", err)
		return
	}

	// 打印接收到的消息
	fmt.Printf("收到消息: %s", message)

	s.Close()
}

func connect(host host.Host, addr string) (network.Stream, error) {
	// 解析多地址
	maddr, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}

	// 从多地址中提取对等节点信息
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		return nil, err
	}

	// 连接对等节点
	fmt.Println("Connecting to", *info)
	err = host.Connect(context.Background(), *info)
	if err != nil {
		return nil, err
	}

	fmt.Println(host.Network().Connectedness(info.ID))
	// 打开新流
	return host.NewStream(context.Background(), info.ID, utils.PROTOCAL_CONNECTOR)
}
