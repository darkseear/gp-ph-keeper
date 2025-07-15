# Makefile для кроссплатформенной сборки

# Версия приложения (можно передавать при вызове make)
BUILD_VERSION = 0.0.1

# Имя бинарного файла
APP_NAME = gophkeeper
SERVER_APP_NAME = serverGophkeeper

# Директория для бинарных файлов
BIN_DIR_SERVER = server/bin
BIN_DIR_CLIENT = client/bin

# Платформы для сборки
PLATFORMS = windows linux darwin
ARCHS = amd64 arm64

# Флаги для сборки
LDFLAGS = -X main.buildVersion=v$(BUILD_VERSION) -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'
  
.PHONY: build-server build-all build-client-win clean test

# Сборка для всех платформ и архитектур
build-all:
	@for platform in $(PLATFORMS); do \
		for arch in $(ARCHS); do \
			if [ "$$platform" = "windows" ]; then \
				EXT=".exe"; \
			else \
				EXT=""; \
			fi; \
			GOOS=$$platform GOARCH=$$arch go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR_CLIENT)/$(APP_NAME)-$$platform-$$arch$$EXT ./client/cmd/main.go; \
			echo "Built client/$(BIN_DIR)/$(APP_NAME)-$(BUILD_VERSION)-$$platform-$$arch$$EXT"; \
		done \
	done

# Сборка только для Windows 
build-client-win:
	@mkdir -p $(BIN_DIR_CLIENT)
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR_CLIENT)/$(APP_NAME)-$(BUILD_VERSION).exe ./client/cmd/main.go
	@echo "Built Windows executable: $(BIN_DIR_CLIENT)/$(APP_NAME)-$(BUILD_VERSION).exe"

#Сборка только для Windows 
build-server-win:
	@mkdir -p $(BIN_DIR_SERVER)
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(BIN_DIR_SERVER)/$(SERVER_APP_NAME)-$(BUILD_VERSION).exe ./server/cmd/main.go
	@echo "Built Windows executable: $(BIN_DIR_SERVER)/$(SERVER_APP_NAME)-$(BUILD_VERSION).exe"

# Очистка бинарных файлов
clean-client:
	@rm -rf $(BIN_DIR_CLIENT)
	@echo "Cleaned bin directory"
clean-server:
	@rm -rf $(BIN_DIR_SERVER)
	@echo "Cleaned bin directory"

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out	