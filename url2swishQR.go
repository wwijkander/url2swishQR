package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	qrcode "github.com/skip2/go-qrcode"
)

func generateQR(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	refererPattern := `(^https://donera\.scouterstottar\.se/.*$|^https://.*\.adalo.com/.*$|^$)`
	var refererExp = regexp.MustCompile(refererPattern)

	referer := r.Referer()
	if !refererExp.MatchString(referer) {
		http.Error(w, "request denied", 403)
		return
	}

	amount := r.URL.Query().Get("amount")
	message := r.URL.Query().Get("message")

	amountPattern := `^[0-9]+$`
	var amountExp = regexp.MustCompile(amountPattern)

	messagePattern := `^.{1,50}$`
	var messageExp = regexp.MustCompile(messagePattern)

	if !amountExp.MatchString(amount) || !messageExp.MatchString(message) {
		http.Error(w, "invalid amount and/or message", 400)
		return
	}

	var out []byte
	out, err := qrcode.Encode("C1236145155;"+amount+";"+message+";0", qrcode.Medium, 300)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(out)

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/generateQR", generateQR)
	log.Fatal("Error: " + http.ListenAndServe(":"+port, mux).Error())
	log.Printf("Listening on port %s", port)
}
