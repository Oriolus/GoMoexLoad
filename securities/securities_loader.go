package securities

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"time"

	models "github.com/oriolus/moex/securities/models"
)

const DestinationDir 	= "securities_info"
const MoexBaseUrl 		= "http://iss.moex.com/iss/securities/"

type LoadStatus string

const (
	LoadStatusOk		= "OK"
	LoadStatusFailed	= "FAILED"
)

type getResult struct {
	SecId  string
	Status string
	Error  string
}

func NewGetResultSuccess(secId string) getResult {
	return getResult{
		SecId:	secId,
		Status:	LoadStatusOk,
	}
}

func NewGetResultFailed(secId string, fialReason string) getResult {
	return getResult {
		SecId:		secId,
		Status:		LoadStatusFailed,
		Error:		fialReason,
	}
}

type StatusManager struct {
	statuses map[string]getResult
}

func NewStatusManager(length int) StatusManager {
	return StatusManager{
		statuses:	make(map[string]getResult, length),
	}
}

func (s *StatusManager) GetResults() []getResult {
	
	result := make([]getResult, len(s.statuses))
	
	ind := 0
	for _, val := range s.statuses {
		result[ind] = val
		ind++
	}
	
	return result
}

func (s *StatusManager) HandleOk(secId string, result string) {
	s.statuses[secId] = NewGetResultSuccess(secId)
	s.writeOkFile(secId, []byte(result))
}

func (s *StatusManager) HandleFail(secId string, errorText string) {
	s.statuses[secId] = NewGetResultFailed(secId, errorText)
	s.writeErrorFile(secId, errorText)
}

func (s *StatusManager) writeErrorFile(secId string, errorTest string) {
	ioutil.WriteFile(getFileName(secId, LoadStatusFailed), []byte(errorTest), fs.ModeCharDevice)
}

func (s *StatusManager) writeOkFile(secId string, result []byte) {
	ioutil.WriteFile(getFileName(secId, LoadStatusOk), result, fs.ModeCharDevice)
}

func getFileName(secId string, status LoadStatus) string {
	return DestinationDir + "\\" + string(status) + "_" + secId + ".json"
}

func GetSecurites(filename string) ([]models.SecurityDim, error) {

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var emptySecs []models.SecurityDim
	json.Unmarshal(buf, &emptySecs)

	return emptySecs, nil
}

func ReadHttpGet(secId string) (result *string, err error) {

	url := MoexBaseUrl + secId + ".json"

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	res, resultErr := ioutil.ReadAll(resp.Body)

	if resultErr != nil {
		return nil, resultErr
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(string(res))
	}

	stringRes := string(res)
	return &stringRes, nil
}

func Load(dimensionFile string) {

	securites, err := GetSecurites(dimensionFile)

	if err != nil {
		fmt.Println(err)
		return
	}

	manager :=  NewStatusManager(len(securites))

	for i, sec := range securites {

		str, err := ReadHttpGet(sec.Secid)

		if err != nil {
			manager.HandleFail(sec.Secid, err.Error())
			fmt.Println(sec.Secid, "....failed", i)
		} else {
			manager.HandleOk(sec.Secid, *str)
			fmt.Println(sec.Secid, "....success", i)
		}

		time.Sleep(2 * time.Second)
	}

	buf, e := json.Marshal(manager.GetResults())

	if e != nil {
		fmt.Println(e.Error())
	}

	ioutil.WriteFile("statuses.json", buf, fs.ModeCharDevice)

}
