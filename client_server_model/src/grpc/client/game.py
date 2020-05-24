#imports for grpc
import grpc
import blob_pb2
import blob_pb2_grpc

import pygame,random,math
import asyncio

pygame.init()
PLAYER_COLORS = [(37,7,255),(35,183,253),(48,254,241),(19,79,251),(255,7,230),(255,7,23),(6,254,13)]
FOOD_COLORS = [(80,252,54),(36,244,255),(243,31,46),(4,39,243),(254,6,178),(255,211,7),(216,6,254),(145,255,7),(7,255,182),(255,6,86),(147,7,255)]

SCREEN_WIDTH = 800
SCREEN_HEIGHT = 500
surface = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))

FOOD_MASS = 7
ZOOM_CONSTANT = 1
MAP_LENGTH = 10000

t_surface = pygame.Surface((95,25),pygame.SRCALPHA) #transparent rect for score
t_lb_surface = pygame.Surface((155,278),pygame.SRCALPHA) #transparent rect for leaderboard
t_surface.fill((50,50,50,80))
t_lb_surface.fill((50,50,50,80))

font = pygame.font.SysFont('Ubuntu',20,True)
big_font = pygame.font.SysFont('Ubuntu',24,True)

pygame.display.set_caption("Agar.io")

food_list = list()
clock = pygame.time.Clock()

#grpc constants
channel = grpc.insecure_channel('localhost:3000')
stub = blob_pb2_grpc.BlobStub(channel)

def drawText(message,pos,color=(255,255,255)):
    pos = (int(pos[0]), int(pos[1]))
    surface.blit(font.render(message,1,color),pos)

# Gets euclidean distance between two positions
def getDistance(pos1,pos2):
    px,py = pos1
    p2x,p2y = pos2
    diffX = math.fabs(px-p2x)
    diffY = math.fabs(py-p2y)

    return ((diffX**2)+(diffY**2))**(0.5)

class Camera:
    def __init__(self):
        self.x = 0
        self.y = 0
        self.width = SCREEN_WIDTH
        self.height = SCREEN_HEIGHT
        self.zoom = 0.5

    def center(self,blobOrPos):
        p = blobOrPos
        self.x = (p.startX-(p.x*self.zoom))-p.startX+((SCREEN_WIDTH/2))
        self.y = (p.startY-(p.y*self.zoom))-p.startY+((SCREEN_HEIGHT/2))

class Blob:
    def __init__(self,surface,name = ""):
        initRequest = blob_pb2.InitRequest()
        initResponse = stub.Init(initRequest)
        self.startX = self.x = initResponse.x
        self.startY = self.y = initResponse.y
        self.mass = initResponse.mass
        self.surface = surface
        self.color = PLAYER_COLORS[random.randint(0,len(PLAYER_COLORS)-1)]
        self.name = initResponse.id
        self.pieces = list()
        piece = Piece(surface,(self.x,self.y),self.color,self.mass,self.name)

    def update(self):
        self.move()

    def move(self):
        dX,dY = pygame.mouse.get_pos()
        rotation = math.atan2(dY-(float(SCREEN_HEIGHT)/2),dX-(float(SCREEN_WIDTH)/2))*180/math.pi
        speed = 5-1
        vx = speed * (90-math.fabs(rotation))/90
        vy = 0
        if(rotation < 0):
            vy = -speed + math.fabs(vx)
        else:
            vy = speed - math.fabs(vx)

        # print("start pos: ", dX, dY)
        moveRequest = blob_pb2.MoveRequest()
        moveRequest.id = self.name
        moveRequest.x = vx
        moveRequest.y = vy
        moveResponse = stub.Move(moveRequest)

        # print("end pos: ", moveResponse.x, moveResponse.y)
        self.x = moveResponse.x
        self.y = moveResponse.y

    def draw(self,cam):
        regionRequest = blob_pb2.RegionRequest()
        regionResponse = stub.Region(regionRequest)

        players = regionResponse.players
        # print(players)
        for player in players:
            if player.id == self.name:
                #update player mass
                self.x = player.x 
                self.y = player.y 
                self.mass = player.mass
            col = self.color
            zoom = cam.zoom
            x = cam.x
            y = cam.y
            pygame.draw.circle(self.surface,(col[0]-int(col[0]/3),int(col[1]-col[1]/3),int(col[2]-col[2]/3)),(int(player.x*zoom+x),int(player.y*zoom+y)),int((player.mass/2+3)*zoom))
            pygame.draw.circle(self.surface,col,(int(player.x*cam.zoom+cam.x),int(player.y*cam.zoom+cam.y)),int(player.mass/2*zoom))
            if(len(player.id) > 0):
                fw, fh = font.size(player.id)
                drawText(player.id, (player.x*cam.zoom+cam.x-int(fw/2),player.y*cam.zoom+cam.y-int(fh/2)),(50,50,50))

        foods = regionResponse.foods    
        for food in foods:
            #only draw food if food is on screen

            # color = FOOD_COLORS[random.randint(0,len(FOOD_COLORS)-1)]
            color = FOOD_COLORS[0]
            pygame.draw.circle(self.surface, color, (int((food.x*cam.zoom+cam.x)),int(food.y*cam.zoom+cam.y)),int(FOOD_MASS*cam.zoom))


class Piece:
    def __init__(self,surface,pos,color,mass,name,transition=False):
        self.x,self.y = pos
        self.mass = mass
        self.splitting = transition
        self.surface = surface
        self.name = name

    def draw(self):
        pass

class Food:
    def __init__(self,surface):
        self.x = random.randint(20,1980)
        self.y = random.randint(20,1980)
        self.mass = 7
        self.surface = surface
        self.color = FOOD_COLORS[random.randint(0,len(FOOD_COLORS)-1)]

    # def draw(self,cam):
    #     pygame.draw.circle(self.surface,self.color,(int((self.x*cam.zoom+cam.x)),int(self.y*cam.zoom+cam.y)),int(self.mass*cam.zoom))

def spawn_foods(numOfFoods):
    pass
    # for i in range(numOfFoods):
    #     food = Food(surface)
    #     food_list.append(food)

def draw_grid():
    for i in range(0,MAP_LENGTH,25):
        pygame.draw.line(surface,(230,240,240),(int(0+camera.x),int(i*camera.zoom+camera.y)),(int(MAP_LENGTH*camera.zoom+camera.x),int(i*camera.zoom+camera.y)),3)
        pygame.draw.line(surface,(230,240,240),(int(i*camera.zoom+camera.x),int(+camera.y)),(int(i*camera.zoom+camera.x),int(MAP_LENGTH*camera.zoom+camera.y)),3)

def draw_leaderboard(leaders):
    LEADERBOARD_X_INSET = 157
    LEADERBOARD_Y_INSET = 20
    ROW_Y_OFFSET = 45

    surface.blit(big_font.render("Leaderboard",0,(255,255,255)),(SCREEN_WIDTH-LEADERBOARD_X_INSET,
        LEADERBOARD_Y_INSET))
    for idx, player in enumerate(leaders):
        drawText(player, (SCREEN_WIDTH-LEADERBOARD_X_INSET, (idx+1) * ROW_Y_OFFSET))

def draw_HUD():
    w,h = font.size("Score: "+str(int(blob.mass*2))+" ")
    surface.blit(pygame.transform.scale(t_surface,(w,h)),(8,SCREEN_HEIGHT-30))
    surface.blit(t_lb_surface,(SCREEN_WIDTH-160,15))
    drawText("Score: " + str(int(blob.mass*2)),(10,SCREEN_HEIGHT-30))


camera = Camera()
blob = Blob(surface,"Viliami")
# spawn_foods(2000)

while(True):
    clock.tick(70)
    for e in pygame.event.get():
        if(e.type == pygame.QUIT):
            pygame.quit()
            quit()
    blob.update()
    camera.zoom = ZOOM_CONSTANT/(blob.mass)+0.3
    
    camera.center(blob)
    # print(blob.x, blob.y)
    surface.fill((242,251,255))
    draw_grid()

    # for c in food_list:
    #     c.draw(camera)
    blob.draw(camera)

    draw_HUD()
    draw_leaderboard(['testing'])
    pygame.display.flip()
