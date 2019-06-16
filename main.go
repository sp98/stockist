package main

import (
	"log"

	"github.com/stockist/pkg/instrument"
)

func main() {

	inst := instrument.New()
	if !validateInstrument(inst) {
		return
	}

	inst.StartProcessing()

}

func validateInstrument(inst *instrument.Instrument) bool {
	if len(inst.Exchange) == 0 || len(inst.Interval) == 0 || len(inst.Name) == 0 || len(inst.Symbol) == 0 || len(inst.Token) == 0 || len(inst.AccessToken) == 0 || len(inst.APIKey) == 0 || len(inst.APISecret) == 0 {
		log.Printf("One or more instrument fields are empty - %+v", inst)
		return false
	}

	return true
}
