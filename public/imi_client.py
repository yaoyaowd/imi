import grpc

from public import imi_pb2
from public import imi_pb2_grpc


def create_simple_clients(host):
    channel = grpc.insecure_channel(host)
    stub = imi_pb2_grpc.IMIStub(channel)
    return stub


def search(client, cid):
    response = client.Search(imi_pb2.SearchRequest(cid=cid))
    print response


# Only use it for internal test now.
if __name__ == '__main__':
    client = create_simple_clients('alice:8999')
    search(client, '57de2c3c0f70d0106e2e9846')
