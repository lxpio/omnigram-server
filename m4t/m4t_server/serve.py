from concurrent import futures
import grpc
import audio_pb2
import audio_pb2_grpc
import numpy as np
from scipy.io.wavfile import write
from io import BytesIO

class TextToAudioServicer(audio_pb2_grpc.TextToAudioServicer):
    def ConvertTextToAudio(self, request, context):
        # Convert text to audio data (simple example)
        text = request.text
        sample_rate = 44100
        duration = 5
        frequency = 440
        t = np.linspace(0, duration, int(sample_rate * duration), endpoint=False)
        audio_data = 0.5 * np.sin(2 * np.pi * frequency * t)
        
        # Save audio data to a BytesIO object
        audio_stream = BytesIO()
        write(audio_stream, sample_rate, audio_data)
        
        # Get audio data bytes
        audio_bytes = audio_stream.getvalue()
        
        return audio_pb2.AudioResponse(audio_data=audio_bytes)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    audio_pb2_grpc.add_TextToAudioServicer_to_server(TextToAudioServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Server is running...")
    server.wait_for_termination()

if __name__ == '__main__':
    serve()