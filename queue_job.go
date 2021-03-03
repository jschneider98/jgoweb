package jgoweb

import (
	"database/sql"
	"encoding/json"
	"github.com/gocraft/web"
	"github.com/jschneider98/jgoweb/util"
	"net/url"
	"time"
)

// SystemJob
type SystemJob struct {
	Id          sql.NullString   `json:"Id" validate:"omitempty,uuid"`
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

// DataValues
type DataValues struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Empty new model
func NewSystemJob(ctx ContextInterface) (*SystemJob, error) {
	sj := &SystemJob{Ctx: ctx}
	sj.SetDefaults()

	return sj, nil
}

// Set defaults
func (sj *SystemJob) SetDefaults() {
	sj.SetPriority("90")
	sj.SetQueuedAt(time.Now().Format(time.RFC3339))
}

// New model with data
func NewSystemJobWithData(ctx ContextInterface, req *web.Request) (*SystemJob, error) {
	sj, err := NewSystemJob(ctx)

	if err != nil {
		return nil, err
	}

	err = sj.Hydrate(req)

	if err != nil {
		return nil, err
	}

	return sj, nil
}

// Factory Method
func FetchSystemJobById(ctx ContextInterface, id string) (*SystemJob, error) {
	var sj []SystemJob

	stmt := ctx.Select("*").
		From("system.jobs").
		Where("id = ?", id).
		Limit(1)

	_, err := stmt.Load(&sj)

	if err != nil {
		return nil, err
	}

	if len(sj) == 0 {
		return nil, nil
	}

	sj[0].Ctx = ctx

	return &sj[0], nil
}

//
func (sj *SystemJob) ProcessSubmit(req *web.Request) (string, bool, error) {
	err := sj.Hydrate(req)

	if err != nil {
		return "", false, err
	}

	err = sj.Ctx.GetValidator().Struct(sj)

	if err != nil {
		return util.GetNiceErrorMessage(err, "</br>"), false, nil
	}

	err = sj.Save()

	if err != nil {
		return "", false, err
	}

	return "System Job saved.", true, nil
}

// Hydrate the model with data
func (sj *SystemJob) Hydrate(req *web.Request) error {
	err := req.ParseForm()

	if err != nil {
		return err
	}

	sj.SetId(req.PostFormValue("Id"))
	sj.SetName(req.PostFormValue("Name"))
	sj.SetDescription(req.PostFormValue("Description"))
	sj.SetPriority(req.PostFormValue("Priority"))
	sj.SetData(req.PostFormValue("Data"))
	sj.SetStatus(req.PostFormValue("Status"))
	sj.SetQueuedAt(req.PostFormValue("QueuedAt"))
	sj.SetStartedAt(req.PostFormValue("StartedAt"))
	sj.SetCheckinAt(req.PostFormValue("CheckinAt"))
	sj.SetEndedAt(req.PostFormValue("EndedAt"))
	sj.SetError(req.PostFormValue("Error"))

	return nil
}

// Validate the model
func (sj *SystemJob) IsValid() error {
	return sj.Ctx.GetValidator().Struct(sj)
}

// Insert/Update based on pkey value
func (sj *SystemJob) Save() error {
	err := sj.IsValid()

	if err != nil {
		return err
	}

	if !sj.Id.Valid {
		return sj.Insert()
	} else {
		return sj.Update()
	}
}

// Insert a new record
func (sj *SystemJob) Insert() error {
	query := `
INSERT INTO
system.jobs (name,
	description,
	priority,
	data,
	status,
	queued_at,
	started_at,
	checkin_at,
	ended_at,
	error)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
RETURNING id

`

	stmt, err := sj.Ctx.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	err = stmt.QueryRow(sj.Name,
		sj.Description,
		sj.Priority,
		sj.Data,
		sj.Status,
		sj.QueuedAt,
		sj.StartedAt,
		sj.CheckinAt,
		sj.EndedAt,
		sj.Error).Scan(&sj.Id)

	if err != nil {
		return err
	}

	return nil
}

// Update a record
func (sj *SystemJob) Update() error {
	if !sj.Id.Valid {
		return nil
	}

	_, err := sj.Ctx.Update("system.jobs").
		Set("id", sj.Id).
		Set("name", sj.Name).
		Set("description", sj.Description).
		Set("priority", sj.Priority).
		Set("data", sj.Data).
		Set("status", sj.Status).
		Set("queued_at", sj.QueuedAt).
		Set("started_at", sj.StartedAt).
		Set("checkin_at", sj.CheckinAt).
		Set("ended_at", sj.EndedAt).
		Set("error", sj.Error).
		Where("id = ?", sj.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

// Hard delete a record
func (sj *SystemJob) Delete() error {

	if !sj.Id.Valid {
		return nil
	}

	_, err := sj.Ctx.DeleteFrom("system.jobs").
		Where("id = ?", sj.Id).
		Exec()

	if err != nil {
		return err
	}

	return nil
}

//
func (sj *SystemJob) GetId() string {

	if sj.Id.Valid {
		return sj.Id.String
	}

	return ""
}

//
func (sj *SystemJob) SetId(val string) {

	if val == "" {
		sj.Id.Valid = false
		sj.Id.String = ""

		return
	}

	sj.Id.Valid = true
	sj.Id.String = val
}

//
func (sj *SystemJob) GetName() string {

	if sj.Name.Valid {
		return sj.Name.String
	}

	return ""
}

//
func (sj *SystemJob) SetName(val string) {

	if val == "" {
		sj.Name.Valid = false
		sj.Name.String = ""

		return
	}

	sj.Name.Valid = true
	sj.Name.String = val
}

//
func (sj *SystemJob) GetDescription() string {

	if sj.Description.Valid {
		return sj.Description.String
	}

	return ""
}

//
func (sj *SystemJob) SetDescription(val string) {

	if val == "" {
		sj.Description.Valid = false
		sj.Description.String = ""

		return
	}

	sj.Description.Valid = true
	sj.Description.String = val
}

//
func (sj *SystemJob) GetPriority() string {

	if sj.Priority.Valid {
		return sj.Priority.String
	}

	return ""
}

//
func (sj *SystemJob) SetPriority(val string) {

	if val == "" {
		sj.Priority.Valid = false
		sj.Priority.String = ""

		return
	}

	sj.Priority.Valid = true
	sj.Priority.String = val
}

//
func (sj *SystemJob) GetData() string {

	if sj.Data.Valid {
		return sj.Data.String
	}

	return ""
}

//
func (sj *SystemJob) SetData(val string) {

	if val == "" {
		sj.Data.Valid = false
		sj.Data.String = ""

		return
	}

	sj.Data.Valid = true
	sj.Data.String = val
}

//
func (sj *SystemJob) GetStatus() string {

	if sj.Status.Valid {
		return sj.Status.String
	}

	return ""
}

//
func (sj *SystemJob) SetStatus(val string) {

	if val == "" {
		sj.Status.Valid = false
		sj.Status.String = ""

		return
	}

	sj.Status.Valid = true
	sj.Status.String = val
}

//
func (sj *SystemJob) GetQueuedAt() string {

	if sj.QueuedAt.Valid {
		return sj.QueuedAt.String
	}

	return ""
}

//
func (sj *SystemJob) SetQueuedAt(val string) {

	if val == "" {
		sj.QueuedAt.Valid = false
		sj.QueuedAt.String = ""

		return
	}

	sj.QueuedAt.Valid = true
	sj.QueuedAt.String = val
}

//
func (sj *SystemJob) GetStartedAt() string {

	if sj.StartedAt.Valid {
		return sj.StartedAt.String
	}

	return ""
}

//
func (sj *SystemJob) SetStartedAt(val string) {

	if val == "" {
		sj.StartedAt.Valid = false
		sj.StartedAt.String = ""

		return
	}

	sj.StartedAt.Valid = true
	sj.StartedAt.String = val
}

//
func (sj *SystemJob) GetCheckinAt() string {

	if sj.CheckinAt.Valid {
		return sj.CheckinAt.String
	}

	return ""
}

//
func (sj *SystemJob) SetCheckinAt(val string) {

	if val == "" {
		sj.CheckinAt.Valid = false
		sj.CheckinAt.String = ""

		return
	}

	sj.CheckinAt.Valid = true
	sj.CheckinAt.String = val
}

//
func (sj *SystemJob) GetEndedAt() string {

	if sj.EndedAt.Valid {
		return sj.EndedAt.String
	}

	return ""
}

//
func (sj *SystemJob) SetEndedAt(val string) {

	if val == "" {
		sj.EndedAt.Valid = false
		sj.EndedAt.String = ""

		return
	}

	sj.EndedAt.Valid = true
	sj.EndedAt.String = val
}

//
func (sj *SystemJob) GetError() string {

	if sj.Error.Valid {
		return sj.Error.String
	}

	return ""
}

//
func (sj *SystemJob) SetError(val string) {

	if val == "" {
		sj.Error.Valid = false
		sj.Error.String = ""

		return
	}

	sj.Error.Valid = true
	sj.Error.String = val
}

// ************

//
func (sj *SystemJob) GetDataValues() (url.Values, error) {
	values := url.Values{}
	var data []DataValues

	jsonStr := sj.GetData()

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
func (sj *SystemJob) Fail(err error) error {
	sj.SetError(err.Error())
	sj.SetEndedAt((time.Now()).Format(time.RFC3339))

	return sj.Save()
}

//
func (sj *SystemJob) Start() error {
	sj.SetStartedAt((time.Now()).Format(time.RFC3339))

	return sj.Save()
}

//
func (sj *SystemJob) End() error {
	sj.SetEndedAt((time.Now()).Format(time.RFC3339))

	return sj.Save()
}

//
func (sj *SystemJob) Checkin(status string) error {
	sj.SetStatus(status)
	sj.SetCheckinAt((time.Now()).Format(time.RFC3339))

	return sj.Save()
}
