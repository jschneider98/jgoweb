package jgoweb

import (
	"database/sql"
	"encoding/json"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
	"net/url"
	"time"
)

// QueueJob
type QueueJob struct {
	Id          sql.NullString   `json:"Id" validate:"omitempty,uuid"`
	AccountId   sql.NullString   `json:"AccountId" validate:"required,uuid"`
	Name        sql.NullString   `json:"Name" validate:"required,min=1,max=255"`
	Description sql.NullString   `json:"Description" validate:"required,min=1,max=255"`
	Priority    sql.NullString   `json:"Priority" validate:"omitempty,int"`
	Data        sql.NullString   `json:"Data" validate:"omitempty"`
	Status      sql.NullString   `json:"Status" validate:"omitempty,min=1,max=255"`
	QueuedAt    sql.NullString   `json:"QueuedAt" validate:"omitempty,rfc3339"`
	StartedAt   sql.NullString   `json:"StartedAt" validate:"omitempty,rfc3339"`
	CheckinAt   sql.NullString   `json:"CheckinAt" validate:"omitempty,rfc3339"`
	EndedAt     sql.NullString   `json:"EndedAt" validate:"omitempty,rfc3339"`
	Error       sql.NullString   `json:"Error" validate:"omitempty"`
	Ctx         ContextInterface `json:"-" validate:"-"`
}

// Empty new model
func NewQueueJob(ctx ContextInterface) (*QueueJob, error) {
	qj := &QueueJob{Ctx: ctx}
	qj.SetDefaults()

	return qj, nil
}

// Set defaults
func (qj *QueueJob) SetDefaults() {
	qj.SetPriority("90")
	qj.SetQueuedAt(time.Now().Format(time.RFC3339))
}

// New model with data
func NewQueueJobWithData(ctx ContextInterface, req *web.Request) (*QueueJob, error) {
	qj, err := NewQueueJob(ctx)

	if err != nil {
		return nil, err
	}

	err = qj.Hydrate(req)

	if err != nil {
		return nil, err
	}

	return qj, nil
}

// Factory Method
func FetchQueueJobById(ctx ContextInterface, id string) (*QueueJob, error) {
	var qj []QueueJob

	stmt := ctx.Select("*").
		From("queue.jobs").
		Where("id = ?", id).
		Limit(1)

	_, err := stmt.Load(&qj)

	if err != nil {
		return nil, err
	}

	if len(qj) == 0 {
		return nil, nil
	}

	qj[0].Ctx = ctx

	return &qj[0], nil
}

//
func (qj *QueueJob) ProcessSubmit(req *web.Request) (string, bool, error) {
	err := qj.Hydrate(req)

	if err != nil {
		return "", false, err
	}

	err = qj.Ctx.GetValidator().Struct(qj)

	if err != nil {
		return util.GetNiceErrorMessage(err, "</br>"), false, nil
	}

	err = qj.Save()

	if err != nil {
		return "", false, err
	}

	return "Queue Job saved.", true, nil
}

// Hydrate the model with data
func (qj *QueueJob) Hydrate(req *web.Request) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}

	qj.SetId(req.PostFormValue("Id"))
	qj.SetAccountId(req.PostFormValue("AccountId"))
	qj.SetName(req.PostFormValue("Name"))
	qj.SetDescription(req.PostFormValue("Description"))
	qj.SetPriority(req.PostFormValue("Priority"))
	qj.SetData(req.PostFormValue("Data"))
	qj.SetStatus(req.PostFormValue("Status"))
	qj.SetQueuedAt(req.PostFormValue("QueuedAt"))
	qj.SetStartedAt(req.PostFormValue("StartedAt"))
	qj.SetCheckinAt(req.PostFormValue("CheckinAt"))
	qj.SetEndedAt(req.PostFormValue("EndedAt"))
	qj.SetError(req.PostFormValue("Error"))

	return nil
}

// Validate the model
func (qj *QueueJob) IsValid() error {
	return qj.Ctx.GetValidator().Struct(qj)
}

// Insert/Update based on pkey value
func (qj *QueueJob) Save() error {
	err := qj.IsValid()

	if err != nil {
		return err
	}

	if !qj.Id.Valid {
		return qj.Insert()
	} else {
		return qj.Update()
	}
}

// Insert a new record
func (qj *QueueJob) Insert() error {
	query := `
INSERT INTO
queue.jobs (account_id,
	name,
	description,
	priority,
	data,
	status,
	queued_at,
	started_at,
	checkin_at,
	ended_at,
	error)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
RETURNING id

`

	stmt, err := qj.Ctx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(qj.AccountId,
		qj.Name,
		qj.Description,
		qj.Priority,
		qj.Data,
		qj.Status,
		qj.QueuedAt,
		qj.StartedAt,
		qj.CheckinAt,
		qj.EndedAt,
		qj.Error).Scan(&qj.Id)

	if err != nil {
		return err
	}

	return nil
}

// Update a record
func (qj *QueueJob) Update() error {
	if !qj.Id.Valid {
		return nil
	}

	_, err := qj.Ctx.Update("queue.jobs").
		Set("id", qj.Id).
		Set("account_id", qj.AccountId).
		Set("name", qj.Name).
		Set("description", qj.Description).
		Set("priority", qj.Priority).
		Set("data", qj.Data).
		Set("status", qj.Status).
		Set("queued_at", qj.QueuedAt).
		Set("started_at", qj.StartedAt).
		Set("checkin_at", qj.CheckinAt).
		Set("ended_at", qj.EndedAt).
		Set("error", qj.Error).
		Where("id = ?", qj.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

// Hard delete a record
func (qj *QueueJob) Delete() error {

	if !qj.Id.Valid {
		return nil
	}

	_, err := qj.Ctx.DeleteFrom("queue.jobs").
		Where("id = ?", qj.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

//
func (qj *QueueJob) GetId() string {

	if qj.Id.Valid {
		return qj.Id.String
	}

	return ""
}

//
func (qj *QueueJob) SetId(val string) {

	if val == "" {
		qj.Id.Valid = false
		qj.Id.String = ""

		return
	}

	qj.Id.Valid = true
	qj.Id.String = val
}

//
func (qj *QueueJob) GetAccountId() string {

	if qj.AccountId.Valid {
		return qj.AccountId.String
	}

	return ""
}

//
func (qj *QueueJob) SetAccountId(val string) {

	if val == "" {
		qj.AccountId.Valid = false
		qj.AccountId.String = ""

		return
	}

	qj.AccountId.Valid = true
	qj.AccountId.String = val
}

//
func (qj *QueueJob) GetName() string {

	if qj.Name.Valid {
		return qj.Name.String
	}

	return ""
}

//
func (qj *QueueJob) SetName(val string) {

	if val == "" {
		qj.Name.Valid = false
		qj.Name.String = ""

		return
	}

	qj.Name.Valid = true
	qj.Name.String = val
}

//
func (qj *QueueJob) GetDescription() string {

	if qj.Description.Valid {
		return qj.Description.String
	}

	return ""
}

//
func (qj *QueueJob) SetDescription(val string) {

	if val == "" {
		qj.Description.Valid = false
		qj.Description.String = ""

		return
	}

	qj.Description.Valid = true
	qj.Description.String = val
}

//
func (qj *QueueJob) GetPriority() string {

	if qj.Priority.Valid {
		return qj.Priority.String
	}

	return ""
}

//
func (qj *QueueJob) SetPriority(val string) {

	if val == "" {
		qj.Priority.Valid = false
		qj.Priority.String = ""

		return
	}

	qj.Priority.Valid = true
	qj.Priority.String = val
}

//
func (qj *QueueJob) GetData() string {

	if qj.Data.Valid {
		return qj.Data.String
	}

	return ""
}

//
func (qj *QueueJob) SetData(val string) {

	if val == "" {
		qj.Data.Valid = false
		qj.Data.String = ""

		return
	}

	qj.Data.Valid = true
	qj.Data.String = val
}

//
func (qj *QueueJob) GetStatus() string {

	if qj.Status.Valid {
		return qj.Status.String
	}

	return ""
}

//
func (qj *QueueJob) SetStatus(val string) {

	if val == "" {
		qj.Status.Valid = false
		qj.Status.String = ""

		return
	}

	qj.Status.Valid = true
	qj.Status.String = val
}

//
func (qj *QueueJob) GetQueuedAt() string {

	if qj.QueuedAt.Valid {
		return qj.QueuedAt.String
	}

	return ""
}

//
func (qj *QueueJob) SetQueuedAt(val string) {

	if val == "" {
		qj.QueuedAt.Valid = false
		qj.QueuedAt.String = ""

		return
	}

	qj.QueuedAt.Valid = true
	qj.QueuedAt.String = val
}

//
func (qj *QueueJob) GetStartedAt() string {

	if qj.StartedAt.Valid {
		return qj.StartedAt.String
	}

	return ""
}

//
func (qj *QueueJob) SetStartedAt(val string) {

	if val == "" {
		qj.StartedAt.Valid = false
		qj.StartedAt.String = ""

		return
	}

	qj.StartedAt.Valid = true
	qj.StartedAt.String = val
}

//
func (qj *QueueJob) GetCheckinAt() string {

	if qj.CheckinAt.Valid {
		return qj.CheckinAt.String
	}

	return ""
}

//
func (qj *QueueJob) SetCheckinAt(val string) {

	if val == "" {
		qj.CheckinAt.Valid = false
		qj.CheckinAt.String = ""

		return
	}

	qj.CheckinAt.Valid = true
	qj.CheckinAt.String = val
}

//
func (qj *QueueJob) GetEndedAt() string {

	if qj.EndedAt.Valid {
		return qj.EndedAt.String
	}

	return ""
}

//
func (qj *QueueJob) SetEndedAt(val string) {

	if val == "" {
		qj.EndedAt.Valid = false
		qj.EndedAt.String = ""

		return
	}

	qj.EndedAt.Valid = true
	qj.EndedAt.String = val
}

//
func (qj *QueueJob) GetError() string {

	if qj.Error.Valid {
		return qj.Error.String
	}

	return ""
}

//
func (qj *QueueJob) SetError(val string) {

	if val == "" {
		qj.Error.Valid = false
		qj.Error.String = ""

		return
	}

	qj.Error.Valid = true
	qj.Error.String = val
}

// ************

// DataValues
type DataValues struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

//
func (qj *QueueJob) GetDataValues() (url.Values, error) {
	values := url.Values{}
	var data []DataValues

	jsonStr := qj.GetData()

	if jsonStr == "" {
		return values, nil
	}

	err := json.Unmarshal([]byte(jsonStr), &data)

	if err != nil {
		return nil, err
	}

	for _, item := range data {
		values.Set(item.Key, item.Value)
	}

	return values, nil
}

//
func (qj *QueueJob) Fail(err error) error {
	qj.SetError(err.Error())
	qj.SetEndedAt((time.Now()).Format(time.RFC3339))

	return qj.Save()
}

//
func (qj *QueueJob) Start() error {
	qj.SetStartedAt((time.Now()).Format(time.RFC3339))

	return qj.Save()
}

//
func (qj *QueueJob) End() error {
	qj.SetEndedAt((time.Now()).Format(time.RFC3339))

	return qj.Save()
}

//
func (qj *QueueJob) Checkin(status string) error {
	qj.SetStatus(status)
	qj.SetCheckinAt((time.Now()).Format(time.RFC3339))

	return qj.Save()
}
