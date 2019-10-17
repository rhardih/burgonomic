package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Data struct {
	Date          time.Time
	Iso_a2        string
	Iso_a3        string
	Currency_code string
	Name          string
	Local_price   *float64
	Dollar_ex     *float64
	Dollar_price  *float64
	//USD_raw       *float64
	//EUR_raw       *float64
	//GBP_raw       *float64
	//JPY_raw       *float64
	//CNY_raw       *float64
	GDP_dollar   *float64
	Adj_price    *float64
	USD_adjusted *float64
	EUR_adjusted *float64
	GBP_adjusted *float64
	JPY_adjusted *float64
	CNY_adjusted *float64
}

type PageData struct {
	Headers  []string
	Records  []Data
	JsonData string
}

var header []string
var data []Data
var sdata [][]string

// ENV vars
var port = mustGetenv("APP_PORT")

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s environment variable not set.", k)
	}
	return v
}

func readData() {
	f, err := os.Open("big-mac-adjusted-index.csv")
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(bufio.NewReader(f))
	var lines [][]string

	header, err = reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	lines, err = reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	sdata = lines

	for _, line := range lines {
		date, err := time.Parse("2006-01-02", line[0])
		if err != nil {
			log.Fatal(err)
		}

		toFloat := func(v string) *float64 {
			if s, err := strconv.ParseFloat(v, 64); err == nil {
				return &s
			}

			return nil
		}

		local_price := toFloat(line[4])
		dollar_ex := toFloat(line[5])
		dollar_price := toFloat(line[6])
		//usd_raw := toFloat(line[7])
		//eur_raw := toFloat(line[8])
		//gbp_raw := toFloat(line[9])
		//jpy_raw := toFloat(line[10])
		//cny_raw := toFloat(line[11])
		gdp_dollar := toFloat(line[7])
		adj_price := toFloat(line[8])
		usd_adjusted := toFloat(line[9])
		eur_adjusted := toFloat(line[10])
		gbp_adjusted := toFloat(line[11])
		jpy_adjusted := toFloat(line[12])
		cny_adjusted := toFloat(line[13])

		if date.Year() == 2018 && date.Month() == 7 {
			d := Data{
				Date:          date,
				Iso_a3:        line[1],
				Iso_a2:        A3ToA2[line[1]],
				Currency_code: line[2],
				Name:          line[3],
				Local_price:   local_price,
				Dollar_ex:     dollar_ex,
				Dollar_price:  dollar_price,
				//USD_raw:       usd_raw,
				//EUR_raw:       eur_raw,
				//GBP_raw:       gbp_raw,
				//JPY_raw:       jpy_raw,
				//CNY_raw:       cny_raw,
				GDP_dollar:   gdp_dollar,
				Adj_price:    adj_price,
				USD_adjusted: usd_adjusted,
				EUR_adjusted: eur_adjusted,
				GBP_adjusted: gbp_adjusted,
				JPY_adjusted: jpy_adjusted,
				CNY_adjusted: cny_adjusted,
			}

			data = append(data, d)
		}
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	log.Println("Listening on port:" + port)
	srv := http.Server{Addr: ":" + port}
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case sig := <-sigs:
				log.Println(sig)
				log.Println("Shutting server down... ")
				err := srv.Shutdown(nil)
				log.Println("Done!")
				if err != nil {
					log.Fatal(err)
				}
				return
			}
		}
	}()

	http.HandleFunc("/", handleHtml)
	http.Handle("/static/", http.StripPrefix("/static/",
		http.FileServer(http.Dir("static"))))

	readData()

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func handleHtml(w http.ResponseWriter, r *http.Request) {
	var err error
	var t *template.Template
	log.Println(r.Method, r.RequestURI)

	if r.RequestURI != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	t, err = template.New("index.html.tmpl").Funcs(funcMap).
		Delims("[[", "]]").
		ParseFiles("html/index.html.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	var tmp []byte
	tmp, err = json.Marshal(data)
	if err != nil {
		log.Println(err)
	}

	data := PageData{
		Headers:  header,
		Records:  data,
		JsonData: string(tmp),
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
