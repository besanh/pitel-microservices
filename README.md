<div align="center">
<h1>Inbox Marketing API Microservices</h1>
</div>

# Git Workflow

Mỗi thành viên khi nhận hoặc đảm nhiệm một tính năng để phát triển, vui lòng clone từ commit từ branch dev gần nhất và đặt tên branch của mình theo format như sau

`dev*{function}*{username hoặc name của người làm}`

Git commit vào mỗi cuối ngày để follow tiến độ.
Git commit tuân theo pattern sau:

`{ACTION}-{nội dung tính năng vắng tất (tối đa 50 ký tự)}`

Ví dụ: UPDATE - add component table campaign.

Sau khi hoàn tất thì tạo một pull request vào branch dev và báo lead review code.

# Microservices:

<div align="center">
<img src="docs/assets/images/grpc-example.webp" />
</div>


# Architecture:

- DB: PostgreSQL
- Programing Language: Go (v1.20 or above)
- Design pattern: Repository Pattern
- Protocol: HTTP (API REST) + gRPC
- Cache: Redis, Mem Cache
- Queue Messages: RabbitMQ, Redis