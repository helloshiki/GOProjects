package online

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Compasses/GOProjects/apiservice/db"
)

type ProxyRoute struct {
	client *http.Client
	url    string
	GrabIF string
	db     *db.ReplayDB
}

func NewProxyHandler(newurl, grabIF string) *ProxyRoute {
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	if grabIF != "" {
		go func() {
			for {
				// Wait for 10s.
				time.Sleep(10 * time.Second)
				if (FailNum + SuccNum) > 0 {
					log.Printf("\n\tIF: %s SuccNum:%d FailNum:%d FailureRate:%f\n\n", grabIF, SuccNum, FailNum, float32((FailNum))/float32((FailNum+SuccNum)))
				}
			}
		}()
	}
	newDB, err := db.NewReplayDB()
	if err != nil {
		log.Println("Open ReplayDB Error ", err)
	}
	return &ProxyRoute{
		client: &http.Client{Transport: tr},
		url:    newurl,
		GrabIF: grabIF,
		db:     newDB}
}

func (proxy *ProxyRoute) doReq(NeedLog bool, path, method, requestBody string, newRq *http.Request) (resp *http.Response, res []byte) {
	now := time.Now()
	resp, err := proxy.client.Do(newRq)
	if resp != nil {
		defer resp.Body.Close()
	}

	LogOutPut(NeedLog, "Time used: ", time.Since(now))
	if err != nil {
		log.Println("get error ", err)
	} else {
		if resp.Header.Get("Content-Encoding") == "gzip" {
			resp.Body, err = gzip.NewReader(resp.Body)
			if err != nil {
				panic(err)
			}
		}

		res, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("ioutil read err ", err)
		}

		if NeedLog {
			if resp.StatusCode != 200 {
				FailNum++
			} else {
				SuccNum++
			}
		}

		LogOutPut(NeedLog, "Get response : ")
		ResponseFormat(NeedLog, resp, string(res))

		err = proxy.db.StoreRequest(path, method, requestBody, string(res), resp.StatusCode)
		if err != nil {
			log.Println("Store data failed ", err)
		}
		//proxy.db.SerilizeToFile()
	}
	return
}

func (proxy *ProxyRoute) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	newbody := make([]byte, req.ContentLength)
	req.Body.Read(newbody)
	NeedLog := strings.Contains(req.RequestURI, proxy.GrabIF)

	newRq, err := http.NewRequest(req.Method, proxy.url+req.RequestURI, ioutil.NopCloser(bytes.NewReader(newbody)))
	if err != nil {
		log.Println("new request error ", err)
	}

	newRq.Header = req.Header
	path := strings.Split(req.RequestURI, "?")

	LogOutPut(NeedLog, "New Request: ")
	RequstFormat(NeedLog, newRq, string(newbody))
	resphttp, res := proxy.doReq(NeedLog, path[0], req.Method, string(newbody), newRq)
	for key, _ := range resphttp.Header {
		w.Header().Set(key, strings.Join(resphttp.Header[key], ";"))
	}

	w.Write(res)
	// res, err := proxy.db.GetResponse(path[0], req.Method, string(newbody))
	// if err != nil || len(res) == 0 {
	// 	log.Println("Cannot get response from replaydb ", err)
	//
	// } else {
	// 	log.Println("Get response from replaydb ", string(res))
	// 	w.Write(res)
	// }

}
