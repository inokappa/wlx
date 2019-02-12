default: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

depend: ## 依存パッケージの導入
	@dep ensure

test: ## test テストの実行
	@go test -v

build: ## バイナリをビルドする
	@./build.sh

image: ## Docker イメージをビルドする
	@docker build -t wlx .

run: ## Docker コンテナを起動する
	@docker run -d --rm --name=wlx -p 20200:20200 wlx

release: release ## バイナリをリリースする. 引数に `_VER=バージョン番号` を指定する.
	@ghr -u inokappa -r wlx v${_VER} ./pkg/
