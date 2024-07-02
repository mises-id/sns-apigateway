# ssh config mises_alpha
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sns-apigateway ./cmd/main.go
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o sns-apigateway-darwin ./cmd/main.go
upload:
	scp ./main mises_alpha:/apps/sns-apigateway/
replace:
	ssh mises_alpha "mv /apps/sns-apigateway/main /apps/sns-apigateway/sns-apigateway"
restart:
	ssh mises_alpha "sudo supervisorctl restart apigateway"
deploy: build \
	upload \
	replace \
	restart
upload-backup:
	scp ./main mises_backup:/apps/sns-apigateway/
replace-backup:
	ssh mises_backup "mv /apps/sns-apigateway/main /apps/sns-apigateway/sns-apigateway"
restart-backup:
	ssh mises_backup "sudo supervisorctl restart apigateway"
deploy-backup: build \
	upload-backup \
	replace-backup \
	restart-backup
#mises-master
upload-master:
	scp ./main mises_master:/apps/sns-apigateway/
replace-master:
	ssh mises_master "mv /apps/sns-apigateway/main /apps/sns-apigateway/sns-apigateway"
restart-master:
	ssh mises_master "sudo supervisorctl restart apigateway"
deploy-master: build \
	upload-master \
	replace-master \
	restart-master
