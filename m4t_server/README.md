# m4t server

当前 `m4t server`` 提供xTTS对外的grpc服务。

```protobuf
service TextToAudio {
  rpc ConvertTextToAudio(TextRequest) returns (AudioResponse);
  rpc TTSStream(TextRequest) returns (stream WAVListResponse);
}
```


### 1. 部署方式

#### 1.1 下载模型

当前TTS基于 https://github.com/coqui-ai/TTS 实现，默认模型使用的是 [Hugging Face Hub](https://huggingface.co/coqui/XTTS-v2)，
运行前前往下载。

#### 1.2 docker 方式

```bash
# 这里假定你下载的目录为 /opt/MY_TTS/XTTS-v2
# git clone https://huggingface.co/coqui/XTTS-v2 /opt/MY_TTS/XTTS-v2
# 直接启动
docker run --rm -v /opt/MY_TTS/XTTS-v2:/models/XTTS m4t-server:0.0.2
```


#### 1.3 linux 方式

1. 安装 python 环境（略）

2. 前台启动服务

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


### 2. 开发选项

```
conda create -n m4t python=3.10
pip install pipreqs
python3 -m  pipreqs.pipreqs . --force
```