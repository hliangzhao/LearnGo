package main

import (
	`context`
	`errors`
	`google.golang.org/protobuf/proto`
	`io`
	`math`
	`time`

	// 注意，go.mod里面本module定义为github.com/hliangzhao/LearnGo/10-gRPC，
	// 因此本module中的package route自然就是github.com/hliangzhao/LearnGo/10-gRP/route
	pb `github.com/hliangzhao/LearnGo/10-gRPC/route`
	`google.golang.org/grpc`
	`log`
	`net`
)

/*
本文件实现gRPC server，用于响应gRPC client对本server提供的函数的调用。
*/

// 对接口pb.RouteGuideServer的实现
// 首先，嵌入UnimplementedRouteGuideServer这个结构体
// 然后依次实现该接口中声明的方法
type routeGuideServer struct {
	features                         []*pb.Feature // 模拟一个数据库
	pb.UnimplementedRouteGuideServer               // 本接口必然要内嵌这个方法
}

func (s *routeGuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	for _, feature := range s.features {
		if proto.Equal(feature.Location, point) {
			return feature, nil
		}
	}
	return nil, errors.New("point not found in db")
}

func inRange(point *pb.Point, rect *pb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left && float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom && float64(point.Latitude) <= top {
		return true
	}
	return false
}

func (s *routeGuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	for _, feature := range s.features {
		if inRange(feature.Location, rect) {
			err := stream.Send(feature)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

func calcDis(p1, p2 *pb.Point) int32 {
	const CordFactor float64 = 1e7
	const R = float64(6371000)
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) + math.Cos(lat1)*math.Cos(lat2)*math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func (s *routeGuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	startTime := time.Now()
	var pointCnt, dis int32
	var prevPoint *pb.Point
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			// make summary
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:  pointCnt,
				Distance:    dis,
				ElapsedTime: int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCnt++
		if prevPoint != nil {
			dis += calcDis(prevPoint, point)
		}
		prevPoint = point
	}
}

func (s *routeGuideServer) recommendOnce(request *pb.RecommendationRequest) (*pb.Feature, error) {
	var nearest, farthest *pb.Feature
	var nearestDistance, farthestDistance int32

	for _, feature := range s.features {
		distance := calcDis(feature.Location, request.Point)
		if nearest == nil || distance < nearestDistance {
			nearestDistance = distance
			nearest = feature
		}
		if farthest == nil || distance > farthestDistance {
			farthestDistance = distance
			farthest = feature
		}
	}
	if request.Mode == pb.RecommendationMode_GetFarthest {
		return farthest, nil
	} else {
		return nearest, nil
	}
}

func (s *routeGuideServer) Recommend(biStream pb.RouteGuide_RecommendServer) error {
	for {
		request, err := biStream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		recommended, err := s.recommendOnce(request)
		if err != nil {
			return err
		}

		return biStream.Send(recommended)
	}
}

func newRouteGuideServer() *routeGuideServer {
	return &routeGuideServer{features: []*pb.Feature{
		{
			Name: "浙江大学玉泉校区 浙江省杭州市西湖区浙大路38号",
			Location: &pb.Point{
				Latitude:  30306202,
				Longitude: 120084879,
			},
		},
		{
			Name: "西溪国家湿地公园 浙江省杭州市西溪湿地",
			Location: &pb.Point{
				Latitude:  30268839,
				Longitude: 120063346,
			},
		},
		{
			Name: "三潭映月 浙江省杭州市西湖区西湖景点",
			Location: &pb.Point{
				Latitude:  30238665,
				Longitude: 120144978,
			},
		},
	}}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:5000")
	if err != nil {
		log.Fatalln("cannot create a listener at this address")
	}
	// 起一个gPRC server并将proto文件中定义的service（routeGuide server）实例添加/注册到本gPRC server，
	// 然后绑定对应的监听端口并启动即可
	grpcServer := grpc.NewServer()
	// pb: protoc buffer
	pb.RegisterRouteGuideServer(grpcServer, newRouteGuideServer())
	log.Fatalln(grpcServer.Serve(lis))
}
