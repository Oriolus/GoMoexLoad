package securities

import "time"

type SecurityDim struct {
	Id                  int32
	Secid               string
	Shortname           string
	Regnumber           string
	Name                string
	Isin                string
	Is_traded           int32
	Emitent_id          int32
	Emitent_title       string
	Emitent_inn         string
	Emitent_okpo        string
	Gosreg              string
	Type                string
	Group               string
	Primary_boardid     string
	Marketprice_boardid string
}

type Security struct {
	SecId                string    // Код ценной бумаги
	Name                 string    // Полное наименование
	Shortname            string    // Краткое наименование
	Isin                 string    // ISIN код
	IssueDate            time.Time // Дата начала торгов
	MatDate              time.Time // Дата погашения
	InitialfaceValue     float32   // Первоначальная номинальная стоимость
	FaceUnit             string    // Валюта номинала
	Latname              string    // Английское наименование
	StartDateMoex        time.Time // Дата начала торгов на Московской Бирже
	ListLevel            int32     // Уровень листинга
	DaysToRedemption     int32     // Дней до погашения
	IssueSize            int64     // Объем выпуска
	FaceValue            float32   // Номинальная стоимость
	IsQualifiedInvestors bool      // Бумаги для квалифицированных инвесторов
	CouponFrequency      int32     // Периодичность выплаты купона в год
	CouponDate           time.Time // Дата выплаты купона
	CouponPercent        float32   // Ставка купона, %
	CouponValue          float32   // Сумма купона, в валюте номинала
	EveningSession       bool      // Допуск к вечерней дополнительной торговой сессии
	Typename             string    // Вид\/категория ценной бумаги
	Group                string    // Код типа инструмента
	Type                 string    // Тип бумаги
	Groupname            string    // Типа инструмента
	EmitterId            int32     // Код эмитента
}

type Board struct {
	Secid          string
	Boardid        string
	Title          string
	Board_group_id int32
	Market_id      int32
	Market         string
	Engine_id      int32
	Engine         string
	Is_traded      int32
	Decimals       int32
	History_from   time.Time
	History_till   time.Time
	Listed_from    time.Time
	Listed_till    time.Time
	Is_primary     int32
	Currencyid     string
}

type SecurityWrapper struct {
	Value  Security
	Boards []Board
}
