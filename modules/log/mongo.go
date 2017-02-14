package log

import (
	"github.com/innovandalism/shodan/util"
	"gopkg.in/mgo.v2"
)

func MongoLogger(uri *string, dataChan chan *LogMessage) {
	session, err := mgo.Dial(*uri)
	if err != nil {
		util.ReportThreadError(true, err)
	}

	defer session.Close()

	for {
		msg := <-dataChan
		err := session.DB("").C("log").Insert(msg)
		if err != nil {
			util.ReportThreadError(false, err)
		}
	}

}
