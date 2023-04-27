package securities

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	models "github.com/oriolus/moex/securities/models"
)

const DateFormant = "2006-01-02"

type fileDescription struct {
	Metadata map[string]any
	Columns  []string
	Data     [][]any
}

type fileBoard struct {
	Metadata map[string]map[string]any
	Columns  []string
	Data     [][]any
}

type fileModel struct {
	Description fileDescription
	Boards      fileBoard
}

func parseFloat(value string) float32 {

	val, err := strconv.ParseFloat(value, 32)
	if err != nil {
		panic(value)
	}
	return float32(val)
}

func parseInt(value string) int32 {

	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		panic(value)
	}
	return int32(val)
}

func parseLong(value string) int64 {

	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(value)
	}
	return val
}

func getSecurity(columns []string, values [][]any) models.Security {

	nameIndex, valueIndex := 0, 0
	for ind, column := range columns {
		if column == "name" {
			nameIndex = ind
		} else if column == "value" {
			valueIndex = ind
		}
	}

	sec := models.Security{}

	for _, rawRow := range values {

		name := rawRow[nameIndex].(string)
		value := rawRow[valueIndex]

		if name == "SECID" {
			sec.SecId = value.(string)
		} else if name == "NAME" {
			sec.Name = value.(string)
		} else if name == "SHORTNAME" {
			sec.Shortname = value.(string)
		} else if name == "ISIN" {
			sec.Isin = value.(string)
		} else if name == "ISSUEDATE" {

			date, parseErr := time.Parse(DateFormant, value.(string))
			if parseErr != nil {
				fmt.Println(parseErr)
			} else {
				sec.IssueDate = date
			}

		} else if name == "MATDATE" {

			date, parseErr := time.Parse(DateFormant, value.(string))
			if parseErr != nil {
				fmt.Println(parseErr)
			} else {
				sec.MatDate = date
			}

		} else if name == "INITIALFACEVALUE" {
			sec.InitialfaceValue = parseFloat(value.(string))
		} else if name == "FACEUNIT" {
			sec.FaceUnit = value.(string)
		} else if name == "LATNAME" {
			sec.Latname = value.(string)
		} else if name == "STARTDATEMOEX" {

			date, parseErr := time.Parse(DateFormant, value.(string))
			if parseErr != nil {
				fmt.Println(parseErr)
			} else {
				sec.StartDateMoex = date
			}

		} else if name == "LISTLEVEL" {
			sec.ListLevel = parseInt(value.(string))
		} else if name == "DAYSTOREDEMPTION" {
			sec.DaysToRedemption = parseInt(value.(string))
		} else if name == "ISSUESIZE" {
			sec.IssueSize = parseLong(value.(string))
		} else if name == "FACEVALUE" {
			sec.FaceValue = parseFloat(value.(string))
		} else if name == "ISQUALIFIEDINVESTORS" {
			sec.IsQualifiedInvestors = value.(string) != "0"
		} else if name == "COUPONFREQUENCY" {
			sec.CouponFrequency = parseInt(value.(string))
		} else if name == "COUPONDATE" {

			date, parseErr := time.Parse(DateFormant, value.(string))
			if parseErr != nil {
				fmt.Println(parseErr)
			} else {
				sec.CouponDate = date
			}

		} else if name == "COUPONPERCENT" {
			sec.CouponPercent = parseFloat(value.(string))
		} else if name == "COUPONVALUE" {
			sec.CouponValue = parseFloat(value.(string))
		} else if name == "EVENINGSESSION" {
			sec.EveningSession = value.(string) != "0"
		} else if name == "TYPENAME" {
			sec.Typename = value.(string)
		} else if name == "GROUP" {
			sec.Group = value.(string)
		} else if name == "TYPE" {
			sec.Type = name
		} else if name == "GROUPNAME" {
			sec.Groupname = value.(string)
		} else if name == "EMITTER_ID" {
			sec.EmitterId = parseInt(value.(string))
		}

	}

	return sec
}

func getBoard(columns []string, columnTypes map[string]string, values []any) models.Board {

	var board models.Board
	boardVal := reflect.ValueOf(&board).Elem()

	for ind, column := range columns {

		if values[ind] != nil {

			columnType := columnTypes[column]
			columnSetter := boardVal.FieldByName(strings.Title(column))

			if columnType == "string" {

				columnSetter.SetString(values[ind].(string))

			} else if columnType == "int32" {

				val := int32(values[ind].(float64))
				columnSetter.SetInt(int64(val))

			} else if columnType == "date" {

				date, parseErr := time.Parse(DateFormant, values[ind].(string))
				if parseErr != nil {
					fmt.Println("bad time at ", column, values[ind].(string))
				}

				columnSetter.Set(reflect.ValueOf(date))

			} else {
				panic("unknown type: " + columnType + " for column " + column)
			}

		}

	}

	return board
}

func getBoards(boards fileBoard) []models.Board {

	parsedBoards := make([]models.Board, len(boards.Data))
	metadata := make(map[string]string, len(boards.Data))

	for _, column := range boards.Columns {
		metadata[column] = boards.Metadata[column]["type"].(string)
	}

	for ind, boardData := range boards.Data {
		parsedBoards[ind] = getBoard(boards.Columns, metadata, boardData)
	}

	return parsedBoards
}

func Transform(filename string) (*models.SecurityWrapper, error) {

	buf, readFileErr := ioutil.ReadFile(filename)

	if readFileErr != nil {
		return nil, readFileErr
	}

	var model fileModel
	json.Unmarshal(buf, &model)

	sec := getSecurity(model.Description.Columns, model.Description.Data)
	boards := getBoards(model.Boards)

	result := models.SecurityWrapper{sec, boards}

	return &result, nil
}

func TransformAll(directory string) ([]models.SecurityWrapper, error) {

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	wraps := make([]models.SecurityWrapper, len(files))

	for ind, file := range files {

		wrap, err := Transform(path.Join(directory, file.Name()))

		if err != nil {
			panic(file.Name())
		} else {
			wraps[ind] = *wrap
		}

	}

	return wraps, nil
}
