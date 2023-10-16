# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import pb.m4t_pb2 as m4t__pb2


class TextToAudioStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.ConvertTextToAudio = channel.unary_unary(
                '/m4t.TextToAudio/ConvertTextToAudio',
                request_serializer=m4t__pb2.TextRequest.SerializeToString,
                response_deserializer=m4t__pb2.AudioResponse.FromString,
                )


class TextToAudioServicer(object):
    """Missing associated documentation comment in .proto file."""

    def ConvertTextToAudio(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_TextToAudioServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'ConvertTextToAudio': grpc.unary_unary_rpc_method_handler(
                    servicer.ConvertTextToAudio,
                    request_deserializer=m4t__pb2.TextRequest.FromString,
                    response_serializer=m4t__pb2.AudioResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'm4t.TextToAudio', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class TextToAudio(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def ConvertTextToAudio(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/m4t.TextToAudio/ConvertTextToAudio',
            m4t__pb2.TextRequest.SerializeToString,
            m4t__pb2.AudioResponse.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)