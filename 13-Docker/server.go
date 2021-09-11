package main

import (
	`fmt`
	`html`
	`log`
	`net/http`
)

/*
Dockerfile ---> (build) ---> Image ---> (launch) ---> Container
 */

func main() {
	fmt.Println("Launch server at port 8080")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprintf(w, "hello, %q", html.EscapeString(r.URL.Path)); err != nil {
			log.Fatalln(err)
		}
	})

	log.Fatalln(http.ListenAndServe("localhost:8080", nil))
}
