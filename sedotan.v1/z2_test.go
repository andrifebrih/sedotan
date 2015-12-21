package sedotan

import (
	"fmt"
	"github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/csv"
	_ "github.com/eaciit/dbox/dbc/json"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/toolkit"
	"testing"
	"time"
)

func TestGrab(t *testing.T) {
	t.Skip()

	url := "http://www.ariefdarmawan.com"
	g := NewGrabber(url, "GET", &Config{})
	if e := g.Grab(nil); e != nil {
		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
		return
	}

	fmt.Printf("Result:\n%s\n", g.ResultString()[:200])
}

func TestPost(t *testing.T) {
	url := "http://www.dce.com.cn/PublicWeb/MainServlet"
	GrabConfig := Config{}

	dataurl := toolkit.M{}
	dataurl["Pu00231_Input.trade_date"] = "20151214"
	dataurl["Pu00231_Input.variety"] = "i"
	dataurl["Pu00231_Input.trade_type"] = "0"
	dataurl["Submit"] = "Go"
	dataurl["action"] = "Pu00231_result"

	GrabConfig.SetFormValues(dataurl)
	g := NewGrabber(url, "POST", &GrabConfig)

	g.Config.DataSettings = make(map[string]*DataSetting)

	tempDataSetting := DataSetting{}
	tempDataSetting.RowSelector = "table .table tbody tr"
	tempDataSetting.Column(0, &GrabColumn{Alias: "Contract", Selector: "td:nth-child(1)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "Open", Selector: "td:nth-child(2)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "High", Selector: "td:nth-child(3)"})

	g.Config.DataSettings["SELECT01"] = &tempDataSetting

	if e := g.Grab(nil); e != nil {
		t.Errorf("Unable to grab %s. Error: %s", url, e.Error())
		return
	}

	docs := []toolkit.M{}

	e := g.ResultFromHtml("SELECT01", &docs)
	if e != nil {
		t.Errorf("Unable to read: %s", e.Error())
	}

	for _, doc := range docs {
		fmt.Println(doc)
	}

}

func TestServiceGrabGet(t *testing.T) {

	xGrabService := newGrabService()
	xGrabService.Name = "goldshfecom"
	xGrabService.Url = "http://www.shfe.com.cn/en/products/Gold/"

	xGrabService.SourceType = SourceType_Http

	xGrabService.GrabInterval = 5 * time.Minute
	xGrabService.TimeOutInterval = 1 * time.Minute //time.Hour, time.Minute, time.Second

	xGrabService.TimeOutIntervalInfo = fmt.Sprintf("%v %s", 1, "seconds")

	//==For Data Grab Config/Data Grabber          ===========================================
	// tempGrab.ServGrabber = sedotan.NewGrabber(tempGrab.Url, mapVal.Get("calltype", "").(string), &GrabConfig)

	grabConfig := Config{}
	// if has grabconfig

	// CallType     string
	// FormValues   toolkit.M
	// AuthType     string
	// AuthUserId   string
	// AuthPassword string
	//==================

	xGrabService.ServGrabber = NewGrabber(xGrabService.Url, "GET", &grabConfig)

	//===================================================================

	//==For Data Log          ===========================================

	// 	logpath = tempLogConf.Get("logpath", "").(string)
	// 	filename = tempLogConf.Get("filename", "").(string)
	// 	filepattern = tempLogConf.Get("filepattern", "").(string)

	logpath := "E:\\data\\vale\\log"
	filename := "LOG-GRABSHFETEST"
	filepattern := "20060102"

	logconf, e := toolkit.NewLog(false, true, logpath, filename, filepattern)
	if e != nil {
		t.Errorf("Error Found : ", e)
	}

	xGrabService.Log = logconf
	//===================================================================

	//===================================================================
	//==Data Setting and Destination Save =====================

	xGrabService.ServGrabber.DataSettings = make(map[string]*DataSetting)
	xGrabService.DestDbox = make(map[string]*DestInfo)

	// ==For Every Data Setting ===============================
	tempDataSetting := DataSetting{}
	tempDestInfo := DestInfo{}

	tempDataSetting.RowSelector = "#tab_conbox li:nth-child(2) .sjtable .listshuju tbody tr"
	tempDataSetting.Column(0, &GrabColumn{Alias: "Code", Selector: "td:nth-child(1)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "LongSpeculation", Selector: "td:nth-child(2)"})
	tempDataSetting.Column(0, &GrabColumn{Alias: "ShortSpeculation", Selector: "td:nth-child(3)"})

	xGrabService.ServGrabber.DataSettings["DATA01"] = &tempDataSetting //DATA01 use name in datasettings

	ci := dbox.ConnectionInfo{}
	ci.Host = "E:\\data\\vale\\Data_GrabTest.csv"
	ci.Database = ""
	ci.UserName = ""
	ci.Password = ""
	ci.Settings = toolkit.M{}.Set("useheader", true).Set("delimiter", ",")

	tempDestInfo.collection = ""
	tempDestInfo.desttype = "csv"

	tempDestInfo.IConnection, e = dbox.NewConnection(tempDestInfo.desttype, &ci)
	if e != nil {
		t.Errorf("Error Found : ", e)
	}

	xGrabService.DestDbox["DATA01"] = &tempDestInfo

	//===================================================================

	e = xGrabService.StartService()
	if e != nil {
		t.Errorf("Error Found : ", e)
	} else {

		for i := 0; i < 100; i++ {
			fmt.Printf(".")
			time.Sleep(3000 * time.Millisecond)
		}

		e = xGrabService.StopService()
		if e != nil {
			t.Errorf("Error Found : ", e)
		}
	}
}
