package main

import (
	"log"
	"net/http"
	"strings"
)

var (
	BaseApi API
)

func GetEndpointFromPath(path string) (bool, *Endpoint) {
	var isCategoryGroup bool
	splittedPath := strings.Split(path, "/")

	if len(splittedPath) == 4 {
		isCategoryGroup = true
	}

	// Find category

	var foundCategory *ApiCategory

	if !isCategoryGroup {
		for _, category := range BaseApi.Categories {
			if category.Name == splittedPath[1] {
				foundCategory = &category
				break
			}
		}
	} else {
		var foundCategoryGroup *ApiCategoryGroup

		for _, categoryGroup := range BaseApi.CategoriesGroup {
			if categoryGroup.Name == splittedPath[1] {
				foundCategoryGroup = &categoryGroup
				break
			}
		}

		if foundCategoryGroup == nil {
			return false, nil
		}

		for _, category := range foundCategoryGroup.Categories {
			if category.Name == splittedPath[2] {
				foundCategory = &category
				break
			}
		}
	}

	if foundCategory == nil {
		return false, nil
	}

	// Find Endpoint
	for _, endpoint := range foundCategory.Endpoints {
		if (isCategoryGroup && splittedPath[3] == endpoint.Name) || (!isCategoryGroup && splittedPath[2] == endpoint.Name) {
			return true, &endpoint
		}
	}

	return false, nil
}

func HandleApiRequest(w http.ResponseWriter, r *http.Request) {
	foundEndpoint, endpointData := GetEndpointFromPath(r.URL.Path)

	if !foundEndpoint {
		ip := getClientIP(r)
		log.Printf("%s Tried to access %s, RETURNED: %d (ENDPOINT NOT FOUND)", ip, r.URL.Path, http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if endpointData.Method == r.Method || r.Method == http.MethodOptions {
		endpointData.Handler.ServeHTTP(w, r)
	} else {
		ip := getClientIP(r)
		log.Printf("%s Tried to access %s, RETURNED: %d (METHOD NOT ALLOWED)", ip, r.URL.Path, http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func BuildApi() {
	BaseApi = NewApi()

	FilesCategory := NewApiCategory("f")
	FilesCategory.
		AddEndpoint(Endpoint{ // UPLOAD FILE ENDPOINT
			Name:    "u",
			Method:  http.MethodPost,
			Handler: FilesUpload,
		}).
		AddEndpoint(Endpoint{ // SHOW FILE ENDPOINT
			Name:    "s",
			Method:  http.MethodGet,
			Handler: FilesShow,
		})

	BaseApi.AddCategory(FilesCategory)
}
