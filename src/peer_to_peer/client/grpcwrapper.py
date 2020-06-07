import grpc
import player_pb2_grpc
import player_pb2

from collections import Counter
import time
import json
import sys

class GRPCWrapper:
    def __init__(self, IP, run):
        self.ip = IP
        self.port = IP.split(":")[1]
        self.channel = grpc.insecure_channel(IP)
        self.stub = player_pb2_grpc.PlayerStub(self.channel)
        self.log_file_name = "./logs/" + str(run) + "|" + self.port + ".json"
        self.region_rtts = []
        self.move_rtts = []

    # flushes grpc rtts to file
    def flush(self, sig, frame):
        print("Flushing data to logfile...")
        rtt_data = {}
        rtt_data["region"] = self.region_rtts
        rtt_data["move"] = self.move_rtts
        with open(self.log_file_name, "w+") as f:
            json.dump(rtt_data, f, indent=2)
        sys.exit(0)
    
    def respawn(self):
        self.channel.close()
        time.sleep(10)
        self.channel = grpc.insecure_channel(self.ip)
        self.stub = player_pb2_grpc.PlayerStub(self.channel)

    def init(self):
        initRequest = player_pb2.InitRequest()
        # print('Making init req')
        initResponse = self.stub.Init(initRequest)
        return initResponse

    def region(self):
        regionRequest = player_pb2.RegionRequest()
        start = time.time()
        regionResponse = self.stub.Region(regionRequest)
        end = time.time()
        runtime = end - start
        self.region_rtts.append(runtime)
        return regionResponse

    def move(self, name, dx, dy):
        moveRequest = player_pb2.MoveRequest()
        moveRequest.id = name
        moveRequest.x = dx
        moveRequest.y = dy
        start = time.time()
        moveResponse = self.stub.Move(moveRequest)
        end = time.time()
        runtime = end - start
        self.move_rtts.append(runtime)
        return moveResponse
  
    