package cmd

import (
	"fmt"

	"github.com/ttacon/chalk"
)

var top = chalk.Cyan.Color(`
                     -+*    ++-.              
                  :=**+.     =**+:            
                .+***:        .+**+:          
               -***+            =***+         
              =***=              -***+        
             -+++=                -+++=       
            :++++.     .::::.      ++++-      
           .=+++-    -=++++++=-.   :=+++:`)
var bottom = chalk.Blue.Color(`
      .:-==:...    .============:    ...:-=-:.
      -===-        ==============.       :===-
      :----       .--------------:       ----:
       :---:       --------------.      :---- 
        :::::      .:::::::::::::      :::::. 
         :::::.      ::::::::::.      :::::.  
          .::::..      ......       .::::.    
            .:::::..            ..::::::.     
              .::::::::......::::::::.        
                 .::::::::::::::::.. 
      
`)

func Banner() {
	fmt.Printf("%s%s", top, bottom)
	fmt.Println(chalk.Bold.TextStyle(chalk.White.Color("\t\t   Go Conduit CLI\n")))
}
