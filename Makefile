create-forum-binary:
	go get -u -v github.com/mailcourses/technopark-dbms-forum@master
	go build github.com/mailcourses/technopark-dbms-forum

build:
	docker build -t forum .

run-docker:
	docker run  --memory 2G --log-opt max-size=5M --log-opt max-file=3 -p 5000:5000 -t forum

run-tests-func:
	curl -vvv -X POST http://localhost:5000/api/service/clear
	./technopark-dbms-forum func -u http://localhost:5000/api/ -r report.html

run-tests-perf:
	curl -X 'POST' http://localhost:5000/api/service/clear
	./technopark-dbms-forum fill -u http://localhost:5000/api/ --timeout=900
	./technopark-dbms-forum perf -u http://localhost:5000/api/  --duration=600 --step=60
