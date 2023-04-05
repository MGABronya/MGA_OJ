package selfInspection

import (
	"MGA_OJ/util"
	"log"
	"strconv"
)

func ImgInspection() {
	log.Println("The profile is being detected...")
	if !util.FileExit("config/application.yml") {
		log.Println("ERROR!!!" + "The configuration file does not exist! Please make sure that the configuration file application.yml exists in the config directory under the current directory, the contents of the configuration file can be found at " +
			"https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E9%83%A8%E7%BD%B2%E6%96%87%E6%A1%A3.md")
	}
	log.Println("The default image is being checked for presence...")
	log.Println("If this deployment is not intended to support MGA_OJ, please ignore the results of this self-test.")
	for i := 1; i <= 9; i++ {
		num := strconv.Itoa(i)
		if !util.FileExit("MGA" + num + ".jpg") {
			log.Println("Warning!" + "MGA" + num + ".jpg" + " doesn't exist!")
		}
	}
}
