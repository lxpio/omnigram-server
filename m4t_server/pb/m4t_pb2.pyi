from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class TextRequest(_message.Message):
    __slots__ = ["text", "lang", "audio_id"]
    TEXT_FIELD_NUMBER: _ClassVar[int]
    LANG_FIELD_NUMBER: _ClassVar[int]
    AUDIO_ID_FIELD_NUMBER: _ClassVar[int]
    text: str
    lang: str
    audio_id: str
    def __init__(self, text: _Optional[str] = ..., lang: _Optional[str] = ..., audio_id: _Optional[str] = ...) -> None: ...

class AudioResponse(_message.Message):
    __slots__ = ["audio_data"]
    AUDIO_DATA_FIELD_NUMBER: _ClassVar[int]
    audio_data: bytes
    def __init__(self, audio_data: _Optional[bytes] = ...) -> None: ...
