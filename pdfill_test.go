package pdfill

import "testing"

func TestFill(t *testing.T) {
	form := Form{
		"address": "Mars 11234455555",
		"total":   300.00,
	}

	err := Fill(form, "./invoice_test.pdf", "./result.pdf")
	if err != nil {
		t.Error(err)
	}
}

func TestFdfContent(t *testing.T) {
	form := Form{
		"address": "Mars 112344",
		"total":   300.45,
	}

	expect := "<</T(address)/V(Mars 112344)>><</T(total)/V(300.45)>>"
	result := form.FdfContent()
	if result != expect {
		t.Error("Want FdfContent =", expect)
		t.Error("got", result)
	}
}
