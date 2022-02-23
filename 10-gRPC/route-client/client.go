package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/hliangzhao/LearnGo/10-gRPC/route"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

/*
grpc client模拟对grpc server提供的方法发起调用。
本文件代码可以放在任何位置，通过dial up连接grpc server。
*/

func testGetFeature(client pb.RouteGuideClient) {
	feature, err := client.GetFeature(context.Background(), &pb.Point{
		Latitude:  30306202,
		Longitude: 120084879,
	})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(feature)
}

func testListFeatures(client pb.RouteGuideClient) {
	receivedStream, err := client.ListFeatures(context.Background(), &pb.Rectangle{
		Lo: &pb.Point{
			Latitude:  30306202,
			Longitude: 120084879,
		},
		Hi: &pb.Point{
			Latitude:  30238665,
			Longitude: 120144978,
		},
	})
	if err != nil {
		log.Fatalln(err)
	}

	// 不断从流中读取数据
	for {
		feature, err := receivedStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(feature)
	}
}

func testRecordRoute(client pb.RouteGuideClient) {
	points := []*pb.Point{
		{
			Latitude:  30306202,
			Longitude: 120084879,
		},
		{
			Latitude:  30268839,
			Longitude: 120063346,
		},
		{
			Latitude:  30238665,
			Longitude: 120144978,
		},
	}

	// 发送流
	sendStream, err := client.RecordRoute(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	for idx, point := range points {
		fmt.Printf("sending point #%d\n", idx)
		if err := sendStream.Send(point); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(time.Second)
	}

	// 接收grpc server发来的单次结果
	summary, err := sendStream.CloseAndRecv()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(summary)
}

func readIntFromCMD(reader *bufio.Reader, target *int32) {
	_, err := fmt.Fscanf(reader, "%d\n", target)
	if err != nil {
		log.Fatalln("Cannot scan", err)
	}
}

func testRecommend(client pb.RouteGuideClient) {
	biStream, err := client.Recommend(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	// 从stream中拿值
	go func() {
		feature, err2 := biStream.Recv()
		if err2 != nil {
			log.Fatalln(err2)
		}
		fmt.Println("Recommended:", feature)
	}()

	reader := bufio.NewReader(os.Stdin)

	for {
		// 这里request的Point需要分配空间，否则为nil会报错
		request := pb.RecommendationRequest{Point: new(pb.Point)}
		var mode int32
		fmt.Println("Enter recommendation mode (0 for the farthest, 1 for the nearest) ")
		readIntFromCMD(reader, &mode)

		fmt.Println("Enter Latitude: ")
		readIntFromCMD(reader, &request.Point.Latitude)

		fmt.Println("Enter Longitude: ")
		readIntFromCMD(reader, &request.Point.Longitude)

		request.Mode = pb.RecommendationMode(mode)
		// if mode == 0 {
		// 	request.Mode = pb.RecommendationMode_GetFarthest
		//
		// } else if mode == 1 {
		// 	request.Mode = pb.RecommendationMode_GetNearest
		// } else {
		// 	log.Fatalln("Unsupported mode")
		// }

		if err := biStream.Send(&request); err != nil {
			log.Fatalln(err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {
	// WithInsecure：不要求证书验证
	// WithBlock：阻塞式
	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln("client cannot dial grpc server")
	}
	defer conn.Close()

	// 创建一个client stub，现在已经可以直接像访问本地函数一样直接访问grpc server提供的函数了
	client := pb.NewRouteGuideClient(conn)
	// testGetFeature(client)
	// testListFeatures(client)
	// testRecordRoute(client)
	testRecommend(client)
}
