import json
from pb import m4t_pb2

class Speakers:
    def __init__(self,root_path):
        self.speakers = None
        self.root_path = root_path


    def load(self):
        with open(self.root_path + 'speakers.json', 'r', encoding='utf-8') as file:
            data = json.load(file)
        speakers = [m4t_pb2.Speaker(**d) for d in data]
        print(speakers)
        self.speakers = {speaker.audio_id: speaker for speaker in speakers}


    def save(self):
        serialized_list = [speaker.__dict__ for speaker in self.speakers]
        data = json.dumps(serialized_list, indent=2)
        with open(self.root_path + 'speakers.json', 'w', encoding='utf-8') as file:
            file.write(data)

    def all(self):
        if self.speakers == None:
            self.load()
        return self.speakers.values()
    
    def upsert(self,speaker, audio_data):
        if self.speakers == None:
            self.load()
        return self.speakers