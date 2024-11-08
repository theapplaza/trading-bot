package capital

import (
	"fmt"
)


func Start(){
	//most important thing is authentication
	if err := authenticate(); err != nil {
		fmt.Println(err)
		panic("failed capital auth")
	}

	fmt.Println("Authenticated with capital")

	//start websocket to start listening to price updates to update our local price cache
		fmt.Println("About to call stream quotes")

	if err := streamQuotes(); err != nil {
		fmt.Println(err)
		panic("failed streaming capital quotes")
	}

	fmt.Println("Listening for price updates from capital")


	//run scheduled jobs to check for signals and issue actions to be performed if needed


}