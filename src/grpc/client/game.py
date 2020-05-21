import pygame,random,math

pygame.init()
PLAYER_COLORS = [(37,7,255),(35,183,253),(48,254,241),(19,79,251),(255,7,230),(255,7,23),(6,254,13)]
CELL_COLORS = [(80,252,54),(36,244,255),(243,31,46),(4,39,243),(254,6,178),(255,211,7),(216,6,254),(145,255,7),(7,255,182),(255,6,86),(147,7,255)]

SCREEN_WIDTH = 800
SCREEN_HEIGHT = 500
surface = pygame.display.set_mode((SCREEN_WIDTH, SCREEN_HEIGHT))

t_surface = pygame.Surface((95,25),pygame.SRCALPHA) #transparent rect for score
t_lb_surface = pygame.Surface((155,278),pygame.SRCALPHA) #transparent rect for leaderboard
t_surface.fill((50,50,50,80))
t_lb_surface.fill((50,50,50,80))

font = pygame.font.SysFont('Ubuntu',20,True)
big_font = pygame.font.SysFont('Ubuntu',24,True)

pygame.display.set_caption("Agar.io")

cell_list = list()
clock = pygame.time.Clock()

def drawText(message,pos,color=(255,255,255)):
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

    def centre(self,blobOrPos):
        if(isinstance(blobOrPos,Player)):
            p = blobOrPos
            self.x = (p.startX-(p.x*self.zoom))-p.startX+((SCREEN_WIDTH/2))
            self.y = (p.startY-(p.y*self.zoom))-p.startY+((SCREEN_HEIGHT/2))
        elif(type(blobOrPos) == tuple):
            self.x,self.y = blobOrPos

class Player:
    def __init__(self,surface,name = ""):
        self.startX = self.x = random.randint(100,400)
        self.startY = self.y = random.randint(100,400)
        self.mass = 20
        self.surface = surface
        self.color = PLAYER_COLORS[random.randint(0,len(PLAYER_COLORS)-1)]
        self.name = name
        self.pieces = list()
        piece = Piece(surface,(self.x,self.y),self.color,self.mass,self.name)

    def update(self):
        self.move()
        self.collisionDetection()

    def collisionDetection(self):
        for cell in cell_list:
            if(getDistance((cell.x,cell.y),(self.x,self.y)) <= self.mass/2):
                self.mass+=0.5
                cell_list.remove(cell)

    def move(self):
        dX,dY = pygame.mouse.get_pos()
        rotation = math.atan2(dY-(float(SCREEN_HEIGHT)/2),dX-(float(SCREEN_WIDTH)/2))*180/math.pi
        speed = 4
        vx = speed * (90-math.fabs(rotation))/90
        vy = 0
        if(rotation < 0):
            vy = -speed + math.fabs(vx)
        else:
            vy = speed - math.fabs(vx)
        self.x += vx
        self.y += vy

    def feed(self):
        pass

    def draw(self,cam):
        col = self.color
        zoom = cam.zoom
        x = cam.x
        y = cam.y
        pygame.draw.circle(self.surface,(col[0]-int(col[0]/3),int(col[1]-col[1]/3),int(col[2]-col[2]/3)),(int(self.x*zoom+x),int(self.y*zoom+y)),int((self.mass/2+3)*zoom))
        pygame.draw.circle(self.surface,col,(int(self.x*cam.zoom+cam.x),int(self.y*cam.zoom+cam.y)),int(self.mass/2*zoom))
        if(len(self.name) > 0):
            fw, fh = font.size(self.name)
            drawText(self.name, (self.x*cam.zoom+cam.x-int(fw/2),self.y*cam.zoom+cam.y-int(fh/2)),(50,50,50))

class Piece:
    def __init__(self,surface,pos,color,mass,name,transition=False):
        self.x,self.y = pos
        self.mass = mass
        self.splitting = transition
        self.surface = surface
        self.name = name

    def draw(self):
        pass

class Cell:
    def __init__(self,surface):
        self.x = random.randint(20,1980)
        self.y = random.randint(20,1980)
        self.mass = 7
        self.surface = surface
        self.color = CELL_COLORS[random.randint(0,len(CELL_COLORS)-1)]

    def draw(self,cam):
        pygame.draw.circle(self.surface,self.color,(int((self.x*cam.zoom+cam.x)),int(self.y*cam.zoom+cam.y)),int(self.mass*cam.zoom))

def spawn_cells(numOfCells):
    for i in range(numOfCells):
        cell = Cell(surface)
        cell_list.append(cell)

def draw_grid():
    for i in range(0,2001,25):
        pygame.draw.line(surface,(230,240,240),(0+camera.x,i*camera.zoom+camera.y),(2001*camera.zoom+camera.x,i*camera.zoom+camera.y),3)
        pygame.draw.line(surface,(230,240,240),(i*camera.zoom+camera.x,0+camera.y),(i*camera.zoom+camera.x,2001*camera.zoom+camera.y),3)

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
blob = Player(surface,"Viliami")
spawn_cells(2000)

while(True):
    clock.tick(70)
    for e in pygame.event.get():
        if(e.type == pygame.QUIT):
            pygame.quit()
            quit()
    blob.update()
    camera.zoom = 100/(blob.mass)+0.3
    camera.centre(blob)
    surface.fill((242,251,255))
    draw_grid()

    for c in cell_list:
        c.draw(camera)
    blob.draw(camera)

    draw_HUD()
    draw_leaderboard(['testing'])
    pygame.display.flip()
