PROJECT_NAME := "OpenCoursePlatform-Go"
PKG := "github.com/OpenCoursePlatform/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
 
test-coverage: ## Run tests with coverage
	@go test -p 1 -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt