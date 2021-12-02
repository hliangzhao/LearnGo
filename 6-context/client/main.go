package main

import (
	`context`
	`fmt`
	`io/ioutil`
	`log`
	`net/http`
	`time`
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// TODO：cancel()用于释放与ctx关联的资源。因此cancel()的调用需要放在所有的operations结束之后
	//  对于本案例，直接放在最后即可。
	defer cancel()

	// 用新的、发起http request的方法发送请求，目的是让我们能够把context传进去
	request, err := http.NewRequest(http.MethodGet, "http://localhost:9000", nil)
	if err != nil {
		log.Fatalln(err)
	}
	request = request.WithContext(ctx)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	// 旧的、发起http request的方法
	// response, err := http.Get("http://localhost:9000")
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", resBytes)
	if err := response.Body.Close(); err != nil {
		log.Fatalln(err)
	}
}
