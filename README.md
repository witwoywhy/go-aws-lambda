# AWS Lambda Go + Gin + SAM

ตัวอย่าง AWS Lambda ด้วย Go ที่ใช้ Gin เป็น HTTP router และใช้
[`aws-lambda-go-api-proxy/gin`](https://github.com/awslabs/aws-lambda-go-api-proxy)
แปลง API Gateway event ให้เป็น Gin request

ข้อมูล user เป็น mock แบบ in-memory มีข้อมูลเริ่มต้น 2 รายการ:

```json
[
  {"id": "1", "name": "Alice", "email": "alice@example.com"},
  {"id": "2", "name": "Bob", "email": "bob@example.com"}
]
```

## โครงสร้าง

```text
main.go
├── Gin adapter -> Lambda
└── handler.NewRouter()
    ├── router.go
    ├── create_user.go
    ├── select_user.go
    ├── update_user.go
    ├── delete_user.go
    └── user_store.go
```

`handler/router.go` ลงทะเบียน Gin routes และใช้ `gins.Log(nil)` สำหรับ request/response logging ส่วน handler ของแต่ละ API แยกเป็นคนละไฟล์

## API

| Method | Path | Response |
|---|---|---|
| GET | `/users` | object `{ "users": [...] }`, `200` |
| GET | `/users/:id` | user ตาม id, `200` หรือ `404` |
| POST | `/users` | สร้าง user, `201` |
| PUT | `/users/:id` | แก้ไข user, `200` หรือ `404` |
| DELETE | `/users/:id` | ลบ user, `200` หรือ `404` |

Body สำหรับ create/update:

```json
{"name":"Carol","email":"carol@example.com"}
```

## สิ่งที่ต้องติดตั้ง

- Go
- AWS SAM CLI
- Docker และต้องเปิด Docker service อยู่

ตรวจสอบเครื่องมือ:

```bash
go version
sam --version
docker info
```

## รัน unit test

```bash
cd /home/home-01/projects/aws-lambda
go test ./...
go vet ./...
```

## Build และ validate ด้วย SAM

```bash
export SAM_CLI_TELEMETRY=0
export AWS_DEFAULT_REGION=ap-southeast-1

sam validate --template-file template.yaml --lint
sam build --template-file template.yaml --no-cached
```

ไฟล์ที่ build แล้วอยู่ที่ `.aws-sam/build/LambdaFunction/bootstrap`

## ทดสอบ Lambda ด้วย SAM

`events/local.json` จำลอง `GET /users`

```bash
sam local invoke LambdaFunction \
  --template-file .aws-sam/build/template.yaml \
  --event events/local.json
```

ผลลัพธ์สำคัญที่คาดหวัง:

```json
{
  "statusCode": 200,
  "body": "{\"users\":[{\"id\":\"1\",\"name\":\"Alice\",\"email\":\"alice@example.com\"},{\"id\":\"2\",\"name\":\"Bob\",\"email\":\"bob@example.com\"}]}"
}
```

## ทดสอบผ่าน HTTP ด้วย SAM

เปิด terminal แรก:

```bash
sam local start-api \
  --template-file .aws-sam/build/template.yaml \
  --warm-containers LAZY
```

เปิด terminal ที่สอง แล้วรัน:

```bash
# select all
curl -i http://127.0.0.1:3000/users

# create: ได้ id 3
curl -i -X POST http://127.0.0.1:3000/users \
  -H 'content-type: application/json' \
  -d '{"name":"Carol","email":"carol@example.com"}'

# select one
curl -i http://127.0.0.1:3000/users/3

# update
curl -i -X PUT http://127.0.0.1:3000/users/3 \
  -H 'content-type: application/json' \
  -d '{"name":"Carol Updated","email":"carol.updated@example.com"}'

# delete
curl -i -X DELETE http://127.0.0.1:3000/users/3

# ตรวจว่าถูกลบแล้ว: ได้ 404
curl -i http://127.0.0.1:3000/users/3
```

`--warm-containers LAZY` จำเป็นสำหรับการทดสอบลำดับ create → update → delete เพราะ mock store อยู่ใน memory ของ container เดียว ข้อมูลจะ reset เมื่อเกิด cold start หรือมี Lambda หลาย instance

## ผลการทดสอบล่าสุด

```text
go test ./...                         PASS
go vet ./...                          PASS
sam validate --lint                   PASS
sam build                             PASS
sam local invoke                     PASS (200)
sam local start-api CRUD              PASS (200, 201, 200, 200, 200, 404)
response_body object logging           PASS (`response_body.users` มี user list)
```

## หมายเหตุ

นี่เป็น mock user สำหรับ local test เท่านั้น ไม่ได้บันทึกข้อมูลถาวร หากใช้จริงควรเปลี่ยน `user_store.go` เป็น database หรือ service ภายนอก
