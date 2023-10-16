



from tts_server import TTSModel,ClonerManager



model = TTSModel("./HHD1/XTTS-v1/")

manager = ClonerManager(model)


cloner = manager.get_cloner('female_001')



print("Inference...")

bytes = cloner.text_to_speech("Omnigram 是Flutter编写的支持多平台文件阅读和听书客户端。","zh-cn")

# Specify the file path
file_path = 'output.wav'

# Open the file in binary write mode ('wb')
with open(file_path, 'wb') as file:
    # Write the bytes to the file
    file.write(bytes)
