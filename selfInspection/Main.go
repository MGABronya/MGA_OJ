package selfInspection

import (
	"MGA_OJ/util"
	"log"
)

func MainInspection() {
	log.Println("The profile is being detected...")
	if !util.FileExit("./config/application.yml") {
		log.Println("ERROR!!!" + "The configuration file does not exist! Please make sure that the configuration file application.yml exists in the config directory under the current directory, the contents of the configuration file can be found at " +
			"https://github.com/MGABronya/MGA_OJ/blob/main/document/mgaOJ%E7%9A%84%E9%83%A8%E7%BD%B2%E6%96%87%E6%A1%A3.md")
	}
}
