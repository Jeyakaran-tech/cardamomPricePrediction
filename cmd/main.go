package main

import (
	logger "log"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/Jeyakaran-tech/cardamomPricePrediction/cardamom"
	"google.golang.org/appengine"
)

func main() {

	logger.Print("starting server...")
	functions.HTTP("cardamomData", cardamom.CardamomDataExtract)
	appengine.Main()

}
