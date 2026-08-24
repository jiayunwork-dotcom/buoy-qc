读海洋浮标 CSV 的 Go 命令行质控工具：按范围、突刺、平坦、梯度打标，再输出 Hs、平均周期、八扇区风玫瑰和 Gumbel 50 年重现波高；serve 把同一套分析接到 JSON API。

# buoy-qc 是海洋浮标观测数据质控与海况分析命令行工具

对浮标传感器 CSV 数据执行范围/突刺/平坦/梯度质控，并计算有效波高、平均波周期、风玫瑰图和极值重现水位。

## 构建 / 运行 / 测试

```text
go build ./...
go run . -readings example/readings.csv
go test ./...
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
