package capital

import (
	"fmt"
)

/**
 * This file is used solely for testing during development of this vendor, so that we do not break other vendors by modifying central main func
 */

func main() {
	//most important thing is authentication
	if err := authenticate(); err != nil {
		fmt.Println(err)
		panic("failed capital auth")
	}

	fmt.Println("Authenticated with capital")

	//start websocket to start listening to price updates to update our local price cache
	fmt.Println("About to call stream quotes")

	streamer := &QuoteStreamer{}
	if err := streamer.StreamQuotes(); err != nil {
		fmt.Println(err)
		panic("failed streaming capital quotes")
	}

	fmt.Println("Listening for price updates from capital")
	select {}

}
