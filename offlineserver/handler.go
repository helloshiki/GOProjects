package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"strings"
)

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	r.ParseForm() //解析参数，默认是不会解析的

	log.Println(net.ParseIP(strings.Split(r.RemoteAddr, ":")[0]))
	dec := json.NewDecoder(r.Body)
	var result interface{}
	dec.Decode(&result)
	log.Println("Req:", result)
}

func ATS(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var checkInfo ATSReq
	err := dec.Decode(&checkInfo)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("ATS Req %+v\n", checkInfo)
	rsp := RepoCreateATSRsp(&checkInfo)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(rsp); err != nil {
		panic(err)
	}
}

func RecommandationProducts(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)
	var id RecommandInfo
	err := dec.Decode(&id)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Printf("RecommandProducts Req%+v\n", id)

	RecommandIds := RepoCreateRecommandationProducts(id.ProductId)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(RecommandIds); err != nil {
		panic(err)
	}
	log.Printf("RecommandProducts Rsp%+v\n", RecommandIds)
}

func CreateCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	var customer CustomerCreate
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(&customer)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Customer Create Req%+v\n", customer)

	Account := RepoCreateAccount(customer)

	if err := json.NewEncoder(w).Encode(Account); err != nil {
		panic(err)
	}
	log.Printf("Customer Create Rsp%+v\n", Account)
}

func CustomerAddressNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	dec := json.NewDecoder(r.Body)

	var addInfo CustomerAddress
	err := dec.Decode(&addInfo)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	Rs := RepoCreateAddress(&addInfo)

	if err = json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}

	log.Printf("Create address info %+v\n", addInfo)

}

func CustomerAddressUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	addressId := GetIdFromStr(ps.ByName("id"))
	
	dec := json.NewDecoder(r.Body)

	var addInfo CustomerAddress
	err := dec.Decode(&addInfo)

	if err != nil {
		HandleError(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	Rs := RepoUpdateAddress(addressId, &addInfo)

	if err = json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}
}

func GetCustomerAddress(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	
	customerId := GetIdFromStr(r.Form["$filter"][0])
	Rs := RepoGetCustomerAddress(customerId)
	
	if err := json.NewEncoder(w).Encode(Rs); err != nil {
		panic(err)
	}
}