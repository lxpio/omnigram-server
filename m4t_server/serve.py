#


from concurrent import futures
import grpc
from pb import m4t_pb2_grpc,m4t_pb2
from tts_server import  ClonerManager, TTSModel
import argparse

class TextToAudioServicer(m4t_pb2_grpc.TextToAudioServicer):

    def __init__(self,manager):
        self.manager = manager

    def ConvertTextToAudio(self, request, context):
        # Convert text to audio data (simple example)
        print(request.audio_id, request.text,request.lang)

        clone = self.manager.get_cloner(request.audio_id)
     
        
        # Get audio data bytes
        audio_bytes = clone.text_to_speech(request.text,request.lang)
        
        return m4t_pb2.AudioResponse(audio_data=audio_bytes)



if __name__ == '__main__':

    parser = argparse.ArgumentParser()
    parser.add_argument("--host", type=str, default="localhost")
    parser.add_argument("--port", type=int, default=50051)
    parser.add_argument("--model-path", type=str, default="./HHD1/XTTS-v1/")

    args = parser.parse_args()

    manager = ClonerManager(TTSModel(args.model_path))

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    m4t_pb2_grpc.add_TextToAudioServicer_to_server(TextToAudioServicer(manager), server)
    # server.add_insecure_port('[::]:50051')
    server.add_insecure_port( args.host + ':'+ args.port)

    server.start()
    print("Server is running...")
    server.wait_for_termination()