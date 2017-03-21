
pgo:
	cd client && go build -o $(GOBIN)/pgo pgo.go
clean:
	rm $(GOBIN)/pgo $(GOBIN)/postgres-operator
operatorimage:
	cd operator && go install postgres-operator.go
	cp $(GOBIN)/postgres-operator bin/postgres-operator
	docker build -t postgres-operator -f $(CCP_BASEOS)/Dockerfile.postgres-operator.$(CCP_BASEOS) .
	docker tag postgres-operator crunchydata/postgres-operator:$(CCP_BASEOS)-$(CO_VERSION)
lsimage:
	docker build -t lspvc -f $(CO_BASEOS)/Dockerfile.lspvc.$(CO_BASEOS) .
	docker tag lspvc crunchydata/lspvc:$(CO_BASEOS)-$(CO_VERSION)
all:
	make operatorimage
default:
	all

