package main

import (
	`fmt`
	`log`
	`net/http`
	`strconv`
)

func main() {
	http.HandleFunc("/", doubleNumHandler)
	log.Fatalln(http.ListenAndServe(":9000", nil))
}

func doubleNumHandler(w http.ResponseWriter, r *http.Request) {
	text := r.FormValue("v")
	if text == "" {
		http.Error(w, "missing value v", http.StatusBadRequest)
		return
	}

	v, err := strconv.Atoi(text)
	if err != nil {
		http.Error(w, "NaN", http.StatusBadRequest)
		return
	}

	if _, err := fmt.Fprintf(w, strconv.Itoa(v * 2) + "\n"); err != nil {
		http.Error(w, "cannot double v", http.StatusBadRequest)
		return
	}
}
