from flask import Flask
from flask import jsonify
app = Flask(__name__)
import grpc
import blob_pb2
import blob_pb2_grpc
import time

@app.route('/test')
def move():
    print("Start service")
    try:
      channel = grpc.insecure_channel('localhost:3000')
      stub = blob_pb2_grpc.BlobStub(channel)
      data = blob_pb2.Position()
      data.x = 0
      data.y = 0
      blobRequest = blob_pb2.BlobRequest()
      before = time.time()
      blobResponse = stub.Move(blobRequest)
      after = time.time()
      print(after - before)
      print(blobResponse)
      print("players: ", blobResponse.players.decode("utf-8"))
      print("food: ", blobResponse.food.decode("utf-8"))

      return jsonify({
          "x": blobResponse.position.x,
          "y": blobResponse.position.y,
          "alive": blobResponse.alive,
          "mass": blobResponse.mass,
          "players": blobResponse.players.decode("utf-8"),
          "food": blobResponse.food.decode("utf-8")
      })
    except Exception as e:
      print(e)
      return e
if __name__ == '__main__':
    channel = grpc.insecure_channel('localhost:3000')
    stub = blob_pb2_grpc.BlobStub(channel)
    data = blob_pb2.Position()
    data.x = 0
    data.y = 0
    blobRequest = blob_pb2.BlobRequest()
    before = time.time()
    blobResponse = stub.Move(blobRequest)
    after = time.time()
    print(after - before)
    print(blobResponse)
    print("players: ", blobResponse.players.decode("utf-8"))
    print("food: ", blobResponse.food.decode("utf-8"))
    app.run()
