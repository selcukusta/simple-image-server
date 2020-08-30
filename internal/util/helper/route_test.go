package helper

import (
	"testing"

	"github.com/selcukusta/simple-image-server/internal/util/constant"
)

func TestIsRouteFitForPatternOne(t *testing.T) {

	input := "/i/gdrive/100/100x200/gtc/1cqCxdmf5YRK4KoneQykr8ic_sbWYHQqg"
	val, params := IsRouteFit(constant.Patterns, input)
	if !val {
		t.Error("Route should be matched")
	}

	if len(params) != 6 {
		t.Error("Paramete count should be 6")
	}

	if params["slug"] == "" || params["slug"] != "gdrive" {
		t.Error("Slug should be matched")
	}

	if params["quality"] == "" || params["quality"] != "100" {
		t.Error("Quality should be matched")
	}

	if params["width"] == "" || params["width"] != "100" {
		t.Error("Width should be matched")
	}

	if params["height"] == "" || params["height"] != "200" {
		t.Error("Height should be matched")
	}

	if params["option"] == "" || params["option"] != "gtc" {
		t.Error("Option should be matched")
	}

	if params["path"] == "" || params["path"] != "1cqCxdmf5YRK4KoneQykr8ic_sbWYHQqg" {
		t.Error("Path should be matched")
	}
}

func TestIsRouteFitForPatternTwo(t *testing.T) {

	input := "/i/gdrive/100/100x200/1cqCxdmf5YRK4KoneQykr8ic_sbWYHQqg"
	val, params := IsRouteFit(constant.Patterns, input)
	if !val {
		t.Error("Route should be matched")
	}

	if len(params) != 5 {
		t.Error("Paramete count should be 5")
	}

	if params["slug"] == "" || params["slug"] != "gdrive" {
		t.Error("Slug should be matched")
	}

	if params["quality"] == "" || params["quality"] != "100" {
		t.Error("Quality should be matched")
	}

	if params["width"] == "" || params["width"] != "100" {
		t.Error("Width should be matched")
	}

	if params["height"] == "" || params["height"] != "200" {
		t.Error("Height should be matched")
	}

	if params["path"] == "" || params["path"] != "1cqCxdmf5YRK4KoneQykr8ic_sbWYHQqg" {
		t.Error("Path should be matched")
	}
}

func TestIsRouteFitForPatternThree(t *testing.T) {

	input := "/i/gdrive/webp/100/100x200/ctg/1cqCxdmf5YRK4KoneQykr8ic_sbWYHQqg"
	val, params := IsRouteFit(constant.Patterns, input)
	if !val {
		t.Error("Route should be matched")
	}

	if len(params) != 7 {
		t.Error("Paramete count should be 7")
	}

	if params["slug"] == "" || params["slug"] != "gdrive" {
		t.Error("Slug should be matched")
	}

	if params["webp"] == "" || params["webp"] != "webp" {
		t.Error("Webp should be matched")
	}

	if params["quality"] == "" || params["quality"] != "100" {
		t.Error("Quality should be matched")
	}

	if params["width"] == "" || params["width"] != "100" {
		t.Error("Width should be matched")
	}

	if params["height"] == "" || params["height"] != "200" {
		t.Error("Height should be matched")
	}

	if params["option"] == "" || params["option"] != "ctg" {
		t.Error("Option should be matched")
	}

	if params["path"] == "" || params["path"] != "1cqCxdmf5YRK4KoneQykr8ic_sbWYHQqg" {
		t.Error("Path should be matched")
	}
}

func TestIsRouteFitForPatternFour(t *testing.T) {

	input := "/i/gdrive/webp/100/100x200/1cqCxdmf5YRK4KoneQykr8ic_sbWYHQqg"
	val, params := IsRouteFit(constant.Patterns, input)
	if !val {
		t.Error("Route should be matched")
	}

	if len(params) != 6 {
		t.Error("Paramete count should be 6")
	}

	if params["slug"] == "" || params["slug"] != "gdrive" {
		t.Error("Slug should be matched")
	}

	if params["webp"] == "" || params["webp"] != "webp" {
		t.Error("Webp should be matched")
	}

	if params["quality"] == "" || params["quality"] != "100" {
		t.Error("Quality should be matched")
	}

	if params["width"] == "" || params["width"] != "100" {
		t.Error("Width should be matched")
	}

	if params["height"] == "" || params["height"] != "200" {
		t.Error("Height should be matched")
	}

	if params["path"] == "" || params["path"] != "1cqCxdmf5YRK4KoneQykr8ic_sbWYHQqg" {
		t.Error("Path should be matched")
	}
}
