package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Endpoint struct {
	Name        string `json:"name"`
	GroupName   string `json:"sub_category_name"`
	FullPath    string `json:"full_path"`
	Description string `json:"description"`
	Method      string `json:"method"`
	Handler     http.HandlerFunc
}

type ApiCategoryGroup struct {
	Name       string
	Categories []ApiCategory
}

type ApiCategory struct {
	Name      string
	Endpoints []Endpoint
}

type API struct {
	CategoriesGroup []ApiCategoryGroup
	Categories      []ApiCategory
}

// API

func NewApi() API {
	return API{
		CategoriesGroup: []ApiCategoryGroup{},
		Categories:      []ApiCategory{},
	}
}

func (api *API) AddCategory(category ApiCategory) {
	api.Categories = append(api.Categories, category)
}

func (api *API) AddCategoryGroup(category ApiCategoryGroup) {
	api.CategoriesGroup = append(api.CategoriesGroup, category)
}

// API Categories

func NewApiCategory(name string) ApiCategory {
	return ApiCategory{
		Name: name,
	}
}

func (category *ApiCategory) AddEndpoint(data Endpoint) *ApiCategory {
	if len(data.GroupName) != 0 {
		data.FullPath += fmt.Sprintf("/%s", strings.ToLower(data.GroupName))
	}

	data.FullPath += fmt.Sprintf("/%s", strings.ToLower(category.Name))
	data.FullPath += fmt.Sprintf("/%s", strings.ToLower(data.Name))

	category.Endpoints = append(category.Endpoints, data)

	return category
}

func NewApiCategoryGroup(name string) ApiCategoryGroup {
	return ApiCategoryGroup{
		Name:       name,
		Categories: []ApiCategory{},
	}
}

func (subcategory *ApiCategoryGroup) AddCategory(category ApiCategory) *ApiCategoryGroup {
	subcategory.Categories = append(subcategory.Categories, category)

	return subcategory
}
