
.PHONY: run
run:
	./build/crack-manager -c ./static_config.yaml

.PHONY: generate-coverage
generate-coverage:
	cd ./build
	lcov --capture --gcov-tool /usr/bin/gcov-12 --directory  --output-file coverage.info
	lcov --remove coverage.info '/usr/*' '*/test/*' --output-file coverage_filtered.info 
	genhtml coverage_filtered.info --output-directory coverage_report

# .PHONY: generate-swagger
# generate-swagger:
# 	swag init -d "./" -g "cmd/api/main.go"