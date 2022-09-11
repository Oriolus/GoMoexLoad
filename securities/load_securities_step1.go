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

type GetResult struct {
	SecId  string
	Status string
	Error  string
}

const destinationDir = "securities_info"
const OK_STATUS = "OK"
const FAIL_STATUS = "FAIL"

type StatusManager struct {
	Statuses map[string]GetResult
}

func (s *StatusManager) GetResults() []GetResult {
	result := make([]GetResult, len(s.Statuses))
	ind := 0
	for _, val := range s.Statuses {
		result[ind] = val
		ind++
	}
	return result
}

func (s *StatusManager) HandleOk(secId string, result string) {
	s.Statuses[secId] = GetResult{secId, OK_STATUS, ""}
	s.writeOkFile(secId, []byte(result))
}

func (s *StatusManager) HandleFail(secId string, errorText string) {
	s.Statuses[secId] = GetResult{secId, FAIL_STATUS, errorText}
	s.writeErrorFile(secId, errorText)
}

func (s *StatusManager) writeErrorFile(secId string, errorTest string) {
	ioutil.WriteFile(destinationDir+"\\"+FAIL_STATUS+"_"+secId+".json", []byte(errorTest), fs.ModeCharDevice)
}

func (s *StatusManager) writeOkFile(secId string, result []byte) {
	ioutil.WriteFile(destinationDir+"\\"+OK_STATUS+"_"+secId+".json", result, fs.ModeCharDevice)
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
	baseUrl := "http://iss.moex.com/iss/securities/"
	url := baseUrl + secId + ".json"

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

	manager := StatusManager{make(map[string]GetResult, 0)}

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
