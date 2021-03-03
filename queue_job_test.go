// +build integration

package jgoweb

import (
	"github.com/gocraft/web"
	"net/http"
	"strings"
	"testing"
)

//
func TestFetchSystemJobById(t *testing.T) {
	InitMockCtx()
	InitMockSystemJob()

	// force not found
	id := "00000000-0000-0000-0000-000000000000"
	sj, err := FetchSystemJobById(MockCtx, id)

	if err != nil {
		t.Errorf("\nERROR: Failed to fetch SystemJob by id: %v\n", err)
		return
	}

	if sj != nil {
		t.Errorf("\nERROR: Should have failed to find SystemJob: %v\n", id)
		return
	}

	sj, err = FetchSystemJobById(MockCtx, MockSystemJob.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if sj == nil {
		t.Errorf("\nERROR: Should have found SystemJob with Id: %v\n", MockSystemJob.GetId())
		return
	}

	if sj.GetId() != MockSystemJob.GetId() {
		t.Errorf("\nERROR: Fetch mismatch. Expected: %v Got: %v\n", MockSystemJob.GetId(), sj.GetId())
		return
	}
}

//
func TestSystemJobId(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetId()
	testVal := "test"

	MockSystemJob.SetId("")

	if MockSystemJob.Id.Valid {
		t.Errorf("ERROR: Id should be invalid.\n")
	}

	if MockSystemJob.GetId() != "" {
		t.Errorf("ERROR: Set Id failed. Should have a blank value. Got: %s", MockSystemJob.GetId())
	}

	MockSystemJob.SetId(testVal)

	if !MockSystemJob.Id.Valid {
		t.Errorf("ERROR: Id should be valid.\n")
	}

	if MockSystemJob.GetId() != testVal {
		t.Errorf("ERROR: Set Id failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetId())
	}

	MockSystemJob.SetId(origVal)
}

//
func TestSystemJobName(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetName()
	testVal := "test"

	MockSystemJob.SetName("")

	if MockSystemJob.Name.Valid {
		t.Errorf("ERROR: Name should be invalid.\n")
	}

	if MockSystemJob.GetName() != "" {
		t.Errorf("ERROR: Set Name failed. Should have a blank value. Got: %s", MockSystemJob.GetName())
	}

	MockSystemJob.SetName(testVal)

	if !MockSystemJob.Name.Valid {
		t.Errorf("ERROR: Name should be valid.\n")
	}

	if MockSystemJob.GetName() != testVal {
		t.Errorf("ERROR: Set Name failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetName())
	}

	MockSystemJob.SetName(origVal)
}

//
func TestSystemJobDescription(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetDescription()
	testVal := "test"

	MockSystemJob.SetDescription("")

	if MockSystemJob.Description.Valid {
		t.Errorf("ERROR: Description should be invalid.\n")
	}

	if MockSystemJob.GetDescription() != "" {
		t.Errorf("ERROR: Set Description failed. Should have a blank value. Got: %s", MockSystemJob.GetDescription())
	}

	MockSystemJob.SetDescription(testVal)

	if !MockSystemJob.Description.Valid {
		t.Errorf("ERROR: Description should be valid.\n")
	}

	if MockSystemJob.GetDescription() != testVal {
		t.Errorf("ERROR: Set Description failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetDescription())
	}

	MockSystemJob.SetDescription(origVal)
}

//
func TestSystemJobPriority(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetPriority()
	testVal := "test"

	MockSystemJob.SetPriority("")

	if MockSystemJob.Priority.Valid {
		t.Errorf("ERROR: Priority should be invalid.\n")
	}

	if MockSystemJob.GetPriority() != "" {
		t.Errorf("ERROR: Set Priority failed. Should have a blank value. Got: %s", MockSystemJob.GetPriority())
	}

	MockSystemJob.SetPriority(testVal)

	if !MockSystemJob.Priority.Valid {
		t.Errorf("ERROR: Priority should be valid.\n")
	}

	if MockSystemJob.GetPriority() != testVal {
		t.Errorf("ERROR: Set Priority failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetPriority())
	}

	MockSystemJob.SetPriority(origVal)
}

//
func TestSystemJobData(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetData()
	testVal := "test"

	MockSystemJob.SetData("")

	if MockSystemJob.Data.Valid {
		t.Errorf("ERROR: Data should be invalid.\n")
	}

	if MockSystemJob.GetData() != "" {
		t.Errorf("ERROR: Set Data failed. Should have a blank value. Got: %s", MockSystemJob.GetData())
	}

	MockSystemJob.SetData(testVal)

	if !MockSystemJob.Data.Valid {
		t.Errorf("ERROR: Data should be valid.\n")
	}

	if MockSystemJob.GetData() != testVal {
		t.Errorf("ERROR: Set Data failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetData())
	}

	MockSystemJob.SetData(origVal)
}

//
func TestSystemJobStatus(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetStatus()
	testVal := "test"

	MockSystemJob.SetStatus("")

	if MockSystemJob.Status.Valid {
		t.Errorf("ERROR: Status should be invalid.\n")
	}

	if MockSystemJob.GetStatus() != "" {
		t.Errorf("ERROR: Set Status failed. Should have a blank value. Got: %s", MockSystemJob.GetStatus())
	}

	MockSystemJob.SetStatus(testVal)

	if !MockSystemJob.Status.Valid {
		t.Errorf("ERROR: Status should be valid.\n")
	}

	if MockSystemJob.GetStatus() != testVal {
		t.Errorf("ERROR: Set Status failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetStatus())
	}

	MockSystemJob.SetStatus(origVal)
}

//
func TestSystemJobQueuedAt(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetQueuedAt()
	testVal := "test"

	MockSystemJob.SetQueuedAt("")

	if MockSystemJob.QueuedAt.Valid {
		t.Errorf("ERROR: QueuedAt should be invalid.\n")
	}

	if MockSystemJob.GetQueuedAt() != "" {
		t.Errorf("ERROR: Set QueuedAt failed. Should have a blank value. Got: %s", MockSystemJob.GetQueuedAt())
	}

	MockSystemJob.SetQueuedAt(testVal)

	if !MockSystemJob.QueuedAt.Valid {
		t.Errorf("ERROR: QueuedAt should be valid.\n")
	}

	if MockSystemJob.GetQueuedAt() != testVal {
		t.Errorf("ERROR: Set QueuedAt failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetQueuedAt())
	}

	MockSystemJob.SetQueuedAt(origVal)
}

//
func TestSystemJobStartedAt(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetStartedAt()
	testVal := "test"

	MockSystemJob.SetStartedAt("")

	if MockSystemJob.StartedAt.Valid {
		t.Errorf("ERROR: StartedAt should be invalid.\n")
	}

	if MockSystemJob.GetStartedAt() != "" {
		t.Errorf("ERROR: Set StartedAt failed. Should have a blank value. Got: %s", MockSystemJob.GetStartedAt())
	}

	MockSystemJob.SetStartedAt(testVal)

	if !MockSystemJob.StartedAt.Valid {
		t.Errorf("ERROR: StartedAt should be valid.\n")
	}

	if MockSystemJob.GetStartedAt() != testVal {
		t.Errorf("ERROR: Set StartedAt failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetStartedAt())
	}

	MockSystemJob.SetStartedAt(origVal)
}

//
func TestSystemJobCheckinAt(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetCheckinAt()
	testVal := "test"

	MockSystemJob.SetCheckinAt("")

	if MockSystemJob.CheckinAt.Valid {
		t.Errorf("ERROR: CheckinAt should be invalid.\n")
	}

	if MockSystemJob.GetCheckinAt() != "" {
		t.Errorf("ERROR: Set CheckinAt failed. Should have a blank value. Got: %s", MockSystemJob.GetCheckinAt())
	}

	MockSystemJob.SetCheckinAt(testVal)

	if !MockSystemJob.CheckinAt.Valid {
		t.Errorf("ERROR: CheckinAt should be valid.\n")
	}

	if MockSystemJob.GetCheckinAt() != testVal {
		t.Errorf("ERROR: Set CheckinAt failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetCheckinAt())
	}

	MockSystemJob.SetCheckinAt(origVal)
}

//
func TestSystemJobEndedAt(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetEndedAt()
	testVal := "test"

	MockSystemJob.SetEndedAt("")

	if MockSystemJob.EndedAt.Valid {
		t.Errorf("ERROR: EndedAt should be invalid.\n")
	}

	if MockSystemJob.GetEndedAt() != "" {
		t.Errorf("ERROR: Set EndedAt failed. Should have a blank value. Got: %s", MockSystemJob.GetEndedAt())
	}

	MockSystemJob.SetEndedAt(testVal)

	if !MockSystemJob.EndedAt.Valid {
		t.Errorf("ERROR: EndedAt should be valid.\n")
	}

	if MockSystemJob.GetEndedAt() != testVal {
		t.Errorf("ERROR: Set EndedAt failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetEndedAt())
	}

	MockSystemJob.SetEndedAt(origVal)
}

//
func TestSystemJobError(t *testing.T) {
	InitMockSystemJob()
	origVal := MockSystemJob.GetError()
	testVal := "test"

	MockSystemJob.SetError("")

	if MockSystemJob.Error.Valid {
		t.Errorf("ERROR: Error should be invalid.\n")
	}

	if MockSystemJob.GetError() != "" {
		t.Errorf("ERROR: Set Error failed. Should have a blank value. Got: %s", MockSystemJob.GetError())
	}

	MockSystemJob.SetError(testVal)

	if !MockSystemJob.Error.Valid {
		t.Errorf("ERROR: Error should be valid.\n")
	}

	if MockSystemJob.GetError() != testVal {
		t.Errorf("ERROR: Set Error failed. Expected: %s, Got: %s", testVal, MockSystemJob.GetError())
	}

	MockSystemJob.SetError(origVal)
}

//
func TestSystemJobInsert(t *testing.T) {
	InitMockSystemJob()
	Name := "Name Insert"
	Description := "Description Insert"
	Priority := "Priority Insert"
	Data := "Data Insert"
	Status := "Status Insert"
	QueuedAt := "QueuedAt Insert"
	StartedAt := "StartedAt Insert"
	CheckinAt := "CheckinAt Insert"
	EndedAt := "EndedAt Insert"
	Error := "Error Insert"

	sj, err := NewSystemJob(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: NewSystemJob() failed. %v\n", err)
	}

	sj.SetName(Name)
	sj.SetDescription(Description)
	sj.SetPriority(Priority)
	sj.SetData(Data)
	sj.SetStatus(Status)
	sj.SetQueuedAt(QueuedAt)
	sj.SetStartedAt(StartedAt)
	sj.SetCheckinAt(CheckinAt)
	sj.SetEndedAt(EndedAt)
	sj.SetError(Error)

	err = sj.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !sj.Id.Valid {
		t.Errorf("\nERROR: SystemJob.Id should be set.\n")
	}

	// verify write
	sj, err = FetchSystemJobById(MockCtx, sj.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sj == nil || sj.GetName() != Name || sj.GetDescription() != Description || sj.GetPriority() != Priority || sj.GetData() != Data || sj.GetStatus() != Status || sj.GetQueuedAt() != QueuedAt || sj.GetStartedAt() != StartedAt || sj.GetCheckinAt() != CheckinAt || sj.GetEndedAt() != EndedAt || sj.GetError() != Error {
		t.Errorf("\nERROR: SystemJob does not match save values. Insert failed.\n")
	}
}

//
func TestSystemJobUpdate(t *testing.T) {
	InitMockSystemJob()
	Name := "Name Update"
	Description := "Description Update"
	Priority := "Priority Update"
	Data := "Data Update"
	Status := "Status Update"
	QueuedAt := "QueuedAt Update"
	StartedAt := "StartedAt Update"
	CheckinAt := "CheckinAt Update"
	EndedAt := "EndedAt Update"
	Error := "Error Update"

	MockSystemJob.SetName(Name)
	MockSystemJob.SetDescription(Description)
	MockSystemJob.SetPriority(Priority)
	MockSystemJob.SetData(Data)
	MockSystemJob.SetStatus(Status)
	MockSystemJob.SetQueuedAt(QueuedAt)
	MockSystemJob.SetStartedAt(StartedAt)
	MockSystemJob.SetCheckinAt(CheckinAt)
	MockSystemJob.SetEndedAt(EndedAt)
	MockSystemJob.SetError(Error)

	err := MockSystemJob.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	sj, err := FetchSystemJobById(MockCtx, MockSystemJob.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sj == nil || sj.GetName() != Name || sj.GetDescription() != Description || sj.GetPriority() != Priority || sj.GetData() != Data || sj.GetStatus() != Status || sj.GetQueuedAt() != QueuedAt || sj.GetStartedAt() != StartedAt || sj.GetCheckinAt() != CheckinAt || sj.GetEndedAt() != EndedAt || sj.GetError() != Error {
		t.Errorf("\nERROR: SystemJob does not match save values. Update failed.\n")
	}
}

//
func TestSystemJobDelete(t *testing.T) {
	InitMockSystemJob()
	err := MockSystemJob.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	sj, err := FetchSystemJobById(MockCtx, MockSystemJob.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if sj != nil {
		t.Errorf("\nERROR: Delete failed. Fetch should return nil.\n")
		return
	}

	MockSystemJob = nil
}

//
func TestNewSystemJobWithData(t *testing.T) {
	httpReq, err := http.NewRequest("GET", "http://example.com", nil)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	req := &web.Request{}
	req.Request = httpReq

	_, err = NewSystemJobWithData(MockCtx, req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestSystemJobProcessSubmit(t *testing.T) {
	sj, err := NewSystemJob(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq, err := http.NewRequest("POST", "http://example.com", strings.NewReader("z=post&Name=Name&Description=Description&Priority=Priority&Data=Data&Status=Status&QueuedAt=QueuedAt&StartedAt=StartedAt&CheckinAt=CheckinAt&EndedAt=EndedAt&Error=Error"))

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req := &web.Request{}
	req.Request = httpReq

	msg, saved, err := sj.ProcessSubmit(req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if !saved {
		t.Errorf("\nERROR: %v", msg)
	}
}
