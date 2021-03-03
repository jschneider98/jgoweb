// +build integration

package jgoweb

import (
	"github.com/gocraft/web"
	"net/http"
	"strings"
	"testing"
)

//
func TestFetchQueueJobById(t *testing.T) {
	InitMockCtx()
	InitMockQueueJob()

	// force not found
	id := "00000000-0000-0000-0000-000000000000"
	qj, err := FetchQueueJobById(MockCtx, id)

	if err != nil {
		t.Errorf("\nERROR: Failed to fetch QueueJob by id: %v\n", err)
		return
	}

	if qj != nil {
		t.Errorf("\nERROR: Should have failed to find QueueJob: %v\n", id)
		return
	}

	qj, err = FetchQueueJobById(MockCtx, MockQueueJob.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if qj == nil {
		t.Errorf("\nERROR: Should have found QueueJob with Id: %v\n", MockQueueJob.GetId())
		return
	}

	if qj.GetId() != MockQueueJob.GetId() {
		t.Errorf("\nERROR: Fetch mismatch. Expected: %v Got: %v\n", MockQueueJob.GetId(), qj.GetId())
		return
	}
}

//
func TestQueueJobId(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetId()
	testVal := "test"

	MockQueueJob.SetId("")

	if MockQueueJob.Id.Valid {
		t.Errorf("ERROR: Id should be invalid.\n")
	}

	if MockQueueJob.GetId() != "" {
		t.Errorf("ERROR: Set Id failed. Should have a blank value. Got: %s", MockQueueJob.GetId())
	}

	MockQueueJob.SetId(testVal)

	if !MockQueueJob.Id.Valid {
		t.Errorf("ERROR: Id should be valid.\n")
	}

	if MockQueueJob.GetId() != testVal {
		t.Errorf("ERROR: Set Id failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetId())
	}

	MockQueueJob.SetId(origVal)
}

//
func TestQueueJobAccountId(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetAccountId()
	testVal := "test"

	MockQueueJob.SetAccountId("")

	if MockQueueJob.AccountId.Valid {
		t.Errorf("ERROR: AccountId should be invalid.\n")
	}

	if MockQueueJob.GetAccountId() != "" {
		t.Errorf("ERROR: Set AccountId failed. Should have a blank value. Got: %s", MockQueueJob.GetAccountId())
	}

	MockQueueJob.SetAccountId(testVal)

	if !MockQueueJob.AccountId.Valid {
		t.Errorf("ERROR: AccountId should be valid.\n")
	}

	if MockQueueJob.GetAccountId() != testVal {
		t.Errorf("ERROR: Set AccountId failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetAccountId())
	}

	MockQueueJob.SetAccountId(origVal)
}

//
func TestQueueJobName(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetName()
	testVal := "test"

	MockQueueJob.SetName("")

	if MockQueueJob.Name.Valid {
		t.Errorf("ERROR: Name should be invalid.\n")
	}

	if MockQueueJob.GetName() != "" {
		t.Errorf("ERROR: Set Name failed. Should have a blank value. Got: %s", MockQueueJob.GetName())
	}

	MockQueueJob.SetName(testVal)

	if !MockQueueJob.Name.Valid {
		t.Errorf("ERROR: Name should be valid.\n")
	}

	if MockQueueJob.GetName() != testVal {
		t.Errorf("ERROR: Set Name failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetName())
	}

	MockQueueJob.SetName(origVal)
}

//
func TestQueueJobDescription(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetDescription()
	testVal := "test"

	MockQueueJob.SetDescription("")

	if MockQueueJob.Description.Valid {
		t.Errorf("ERROR: Description should be invalid.\n")
	}

	if MockQueueJob.GetDescription() != "" {
		t.Errorf("ERROR: Set Description failed. Should have a blank value. Got: %s", MockQueueJob.GetDescription())
	}

	MockQueueJob.SetDescription(testVal)

	if !MockQueueJob.Description.Valid {
		t.Errorf("ERROR: Description should be valid.\n")
	}

	if MockQueueJob.GetDescription() != testVal {
		t.Errorf("ERROR: Set Description failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetDescription())
	}

	MockQueueJob.SetDescription(origVal)
}

//
func TestQueueJobPriority(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetPriority()
	testVal := "test"

	MockQueueJob.SetPriority("")

	if MockQueueJob.Priority.Valid {
		t.Errorf("ERROR: Priority should be invalid.\n")
	}

	if MockQueueJob.GetPriority() != "" {
		t.Errorf("ERROR: Set Priority failed. Should have a blank value. Got: %s", MockQueueJob.GetPriority())
	}

	MockQueueJob.SetPriority(testVal)

	if !MockQueueJob.Priority.Valid {
		t.Errorf("ERROR: Priority should be valid.\n")
	}

	if MockQueueJob.GetPriority() != testVal {
		t.Errorf("ERROR: Set Priority failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetPriority())
	}

	MockQueueJob.SetPriority(origVal)
}

//
func TestQueueJobData(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetData()
	testVal := "test"

	MockQueueJob.SetData("")

	if MockQueueJob.Data.Valid {
		t.Errorf("ERROR: Data should be invalid.\n")
	}

	if MockQueueJob.GetData() != "" {
		t.Errorf("ERROR: Set Data failed. Should have a blank value. Got: %s", MockQueueJob.GetData())
	}

	MockQueueJob.SetData(testVal)

	if !MockQueueJob.Data.Valid {
		t.Errorf("ERROR: Data should be valid.\n")
	}

	if MockQueueJob.GetData() != testVal {
		t.Errorf("ERROR: Set Data failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetData())
	}

	MockQueueJob.SetData(origVal)
}

//
func TestQueueJobStatus(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetStatus()
	testVal := "test"

	MockQueueJob.SetStatus("")

	if MockQueueJob.Status.Valid {
		t.Errorf("ERROR: Status should be invalid.\n")
	}

	if MockQueueJob.GetStatus() != "" {
		t.Errorf("ERROR: Set Status failed. Should have a blank value. Got: %s", MockQueueJob.GetStatus())
	}

	MockQueueJob.SetStatus(testVal)

	if !MockQueueJob.Status.Valid {
		t.Errorf("ERROR: Status should be valid.\n")
	}

	if MockQueueJob.GetStatus() != testVal {
		t.Errorf("ERROR: Set Status failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetStatus())
	}

	MockQueueJob.SetStatus(origVal)
}

//
func TestQueueJobQueuedAt(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetQueuedAt()
	testVal := "test"

	MockQueueJob.SetQueuedAt("")

	if MockQueueJob.QueuedAt.Valid {
		t.Errorf("ERROR: QueuedAt should be invalid.\n")
	}

	if MockQueueJob.GetQueuedAt() != "" {
		t.Errorf("ERROR: Set QueuedAt failed. Should have a blank value. Got: %s", MockQueueJob.GetQueuedAt())
	}

	MockQueueJob.SetQueuedAt(testVal)

	if !MockQueueJob.QueuedAt.Valid {
		t.Errorf("ERROR: QueuedAt should be valid.\n")
	}

	if MockQueueJob.GetQueuedAt() != testVal {
		t.Errorf("ERROR: Set QueuedAt failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetQueuedAt())
	}

	MockQueueJob.SetQueuedAt(origVal)
}

//
func TestQueueJobStartedAt(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetStartedAt()
	testVal := "test"

	MockQueueJob.SetStartedAt("")

	if MockQueueJob.StartedAt.Valid {
		t.Errorf("ERROR: StartedAt should be invalid.\n")
	}

	if MockQueueJob.GetStartedAt() != "" {
		t.Errorf("ERROR: Set StartedAt failed. Should have a blank value. Got: %s", MockQueueJob.GetStartedAt())
	}

	MockQueueJob.SetStartedAt(testVal)

	if !MockQueueJob.StartedAt.Valid {
		t.Errorf("ERROR: StartedAt should be valid.\n")
	}

	if MockQueueJob.GetStartedAt() != testVal {
		t.Errorf("ERROR: Set StartedAt failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetStartedAt())
	}

	MockQueueJob.SetStartedAt(origVal)
}

//
func TestQueueJobCheckinAt(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetCheckinAt()
	testVal := "test"

	MockQueueJob.SetCheckinAt("")

	if MockQueueJob.CheckinAt.Valid {
		t.Errorf("ERROR: CheckinAt should be invalid.\n")
	}

	if MockQueueJob.GetCheckinAt() != "" {
		t.Errorf("ERROR: Set CheckinAt failed. Should have a blank value. Got: %s", MockQueueJob.GetCheckinAt())
	}

	MockQueueJob.SetCheckinAt(testVal)

	if !MockQueueJob.CheckinAt.Valid {
		t.Errorf("ERROR: CheckinAt should be valid.\n")
	}

	if MockQueueJob.GetCheckinAt() != testVal {
		t.Errorf("ERROR: Set CheckinAt failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetCheckinAt())
	}

	MockQueueJob.SetCheckinAt(origVal)
}

//
func TestQueueJobEndedAt(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetEndedAt()
	testVal := "test"

	MockQueueJob.SetEndedAt("")

	if MockQueueJob.EndedAt.Valid {
		t.Errorf("ERROR: EndedAt should be invalid.\n")
	}

	if MockQueueJob.GetEndedAt() != "" {
		t.Errorf("ERROR: Set EndedAt failed. Should have a blank value. Got: %s", MockQueueJob.GetEndedAt())
	}

	MockQueueJob.SetEndedAt(testVal)

	if !MockQueueJob.EndedAt.Valid {
		t.Errorf("ERROR: EndedAt should be valid.\n")
	}

	if MockQueueJob.GetEndedAt() != testVal {
		t.Errorf("ERROR: Set EndedAt failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetEndedAt())
	}

	MockQueueJob.SetEndedAt(origVal)
}

//
func TestQueueJobError(t *testing.T) {
	InitMockQueueJob()
	origVal := MockQueueJob.GetError()
	testVal := "test"

	MockQueueJob.SetError("")

	if MockQueueJob.Error.Valid {
		t.Errorf("ERROR: Error should be invalid.\n")
	}

	if MockQueueJob.GetError() != "" {
		t.Errorf("ERROR: Set Error failed. Should have a blank value. Got: %s", MockQueueJob.GetError())
	}

	MockQueueJob.SetError(testVal)

	if !MockQueueJob.Error.Valid {
		t.Errorf("ERROR: Error should be valid.\n")
	}

	if MockQueueJob.GetError() != testVal {
		t.Errorf("ERROR: Set Error failed. Expected: %s, Got: %s", testVal, MockQueueJob.GetError())
	}

	MockQueueJob.SetError(origVal)
}

//
func TestQueueJobInsert(t *testing.T) {
	InitMockQueueJob()
	AccountId := "AccountId Insert"
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

	qj, err := NewQueueJob(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: NewQueueJob() failed. %v\n", err)
	}

	qj.SetAccountId(AccountId)
	qj.SetName(Name)
	qj.SetDescription(Description)
	qj.SetPriority(Priority)
	qj.SetData(Data)
	qj.SetStatus(Status)
	qj.SetQueuedAt(QueuedAt)
	qj.SetStartedAt(StartedAt)
	qj.SetCheckinAt(CheckinAt)
	qj.SetEndedAt(EndedAt)
	qj.SetError(Error)

	err = qj.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if !qj.Id.Valid {
		t.Errorf("\nERROR: QueueJob.Id should be set.\n")
	}

	// verify write
	qj, err = FetchQueueJobById(MockCtx, qj.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if qj == nil || qj.GetAccountId() != AccountId || qj.GetName() != Name || qj.GetDescription() != Description || qj.GetPriority() != Priority || qj.GetData() != Data || qj.GetStatus() != Status || qj.GetQueuedAt() != QueuedAt || qj.GetStartedAt() != StartedAt || qj.GetCheckinAt() != CheckinAt || qj.GetEndedAt() != EndedAt || qj.GetError() != Error {
		t.Errorf("\nERROR: QueueJob does not match save values. Insert failed.\n")
	}
}

//
func TestQueueJobUpdate(t *testing.T) {
	InitMockQueueJob()
	AccountId := "AccountId Update"
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

	MockQueueJob.SetAccountId(AccountId)
	MockQueueJob.SetName(Name)
	MockQueueJob.SetDescription(Description)
	MockQueueJob.SetPriority(Priority)
	MockQueueJob.SetData(Data)
	MockQueueJob.SetStatus(Status)
	MockQueueJob.SetQueuedAt(QueuedAt)
	MockQueueJob.SetStartedAt(StartedAt)
	MockQueueJob.SetCheckinAt(CheckinAt)
	MockQueueJob.SetEndedAt(EndedAt)
	MockQueueJob.SetError(Error)

	err := MockQueueJob.Save()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	qj, err := FetchQueueJobById(MockCtx, MockQueueJob.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if qj == nil || qj.GetAccountId() != AccountId || qj.GetName() != Name || qj.GetDescription() != Description || qj.GetPriority() != Priority || qj.GetData() != Data || qj.GetStatus() != Status || qj.GetQueuedAt() != QueuedAt || qj.GetStartedAt() != StartedAt || qj.GetCheckinAt() != CheckinAt || qj.GetEndedAt() != EndedAt || qj.GetError() != Error {
		t.Errorf("\nERROR: QueueJob does not match save values. Update failed.\n")
	}
}

//
func TestQueueJobDelete(t *testing.T) {
	InitMockQueueJob()
	err := MockQueueJob.Delete()

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	// verify write
	qj, err := FetchQueueJobById(MockCtx, MockQueueJob.GetId())

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	if qj != nil {
		t.Errorf("\nERROR: Delete failed. Fetch should return nil.\n")
		return
	}

	MockQueueJob = nil
}

//
func TestNewQueueJobWithData(t *testing.T) {
	httpReq, err := http.NewRequest("GET", "http://example.com", nil)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}

	req := &web.Request{}
	req.Request = httpReq

	_, err = NewQueueJobWithData(MockCtx, req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
	}
}

//
func TestQueueJobProcessSubmit(t *testing.T) {
	qj, err := NewQueueJob(MockCtx)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq, err := http.NewRequest("POST", "http://example.com", strings.NewReader("z=post&AccountId=AccountId&Name=Name&Description=Description&Priority=Priority&Data=Data&Status=Status&QueuedAt=QueuedAt&StartedAt=StartedAt&CheckinAt=CheckinAt&EndedAt=EndedAt&Error=Error"))

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	req := &web.Request{}
	req.Request = httpReq

	msg, saved, err := qj.ProcessSubmit(req)

	if err != nil {
		t.Errorf("\nERROR: %v\n", err)
		return
	}

	if !saved {
		t.Errorf("\nERROR: %v", msg)
	}
}
