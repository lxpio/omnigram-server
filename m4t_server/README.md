# m4t server

当前 `m4t server`` 提供xTTS对外的grpc服务。

```protobuf
service TextToAudio {
  rpc ConvertTextToAudio(TextRequest) returns (AudioResponse);
}
```


### 1. 部署方式

1. 安装 python 环境（略）

2. 启动grpc 服务

```bash
# conda create -n m4t python=3.10
cd ${PROJECT_DIR}/m4t_server
pip install -r ./requirements.txt

python serve.py

```

3. systemd 服务

将如下两个变量 `MY_PYTHON_PATH`, `MY_MODEL_PATH` 替换为自己实际的目录：

```bash
cd ${PROJECT_DIR}/m4t_server
sudo MY_PYTHON_PATH='/opt/anaconda3/envs/m4t/bin/python' MY_MODEL_PATH='./model/xtts_v1' ./install.sh

```


### 开发选项

```
conda create -n m4t python=3.10
pip install pipreqs
python3 -m  pipreqs.pipreqs . --force
```


<!-- # systemctl start llmchain.service -->