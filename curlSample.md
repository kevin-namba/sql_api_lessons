user/getへのリクエスト
curl -X GET  -H "x-token: sampletoken1"   http://localhost:8080/user/get

user/create
curl -X POST  -H "Content-Type: application/json" -d '{"name":"kevin6"}'
http://192.168.33.10:8080/user/create

user/update
curl -X POST  -H "xtoken: rXYsplSo" -H "Content-Type: application/json" -d '{"name":"kevin5"}' http://192.168.33.10:8080/user/update
