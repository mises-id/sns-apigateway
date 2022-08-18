# ssh config mises_alpha
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/main.go
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
