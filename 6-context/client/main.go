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
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second * 5)
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

	defer response.Body.Close()
	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%s\n", resBytes)
}
