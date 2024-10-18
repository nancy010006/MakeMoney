 
# 使用官方 Go 鏡像作為構建環境
FROM golang:1.22.1 AS builder

# 設置工作目錄
WORKDIR /app

# 複製 go mod 和 sum 文件
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製源代碼
COPY . .

# 構建應用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# 使用輕量級的 alpine 鏡像
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 從構建器階段複製二進制文件
COPY --from=builder /app/main .

# 暴露應用端口
EXPOSE 8080

# 運行應用
CMD ["./main"]