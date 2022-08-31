package tests

import (
	"github.com/revel/revel/testing"
)

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

//func (t *AppTest) TestSaveHotels() {
//	arr := []models.Hotel{
//		{ID: 23, Name: "TestName", Price: 124},
//		{ID: 23, Name: "", Price: 122},
//		{Name: ";;s;;", Price: 122},
//	}
//	for i := range arr {
//		data, _ := json.Marshal(i)
//		reader := bytes.NewReader(data)
//		t.Post("/save", "application/json", reader)
//		t.
//	}
//}

func (t *AppTest) After() {
	println("Tear down")
}
