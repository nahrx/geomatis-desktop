import json
import cv2

def rasterFeaturePoints_tes(tes):
    return tes
def rasterFeaturePoints(filePath,resize=False,rotation=0,view=False):
    imgRaw = cv2.imread(filePath)
    scale_percent = 40
    if resize:
         # percent of original size
        width = int(imgRaw.shape[1] * scale_percent / 100)
        height = int(imgRaw.shape[0] * scale_percent / 100)
        dim = (width, height)
        imgRaw = cv2.resize(imgRaw, dim, interpolation = cv2.INTER_AREA)
    
    if rotation == 0:
        img = imgRaw
    elif rotation == 90:
        img = cv2.rotate(imgRaw, cv2.ROTATE_90_CLOCKWISE)
    elif rotation == 180:
        img = cv2.rotate(imgRaw, cv2.ROTATE_180)
    elif rotation == -90:
        img = cv2.rotate(imgRaw, cv2.ROTATE_90_COUNTERCLOCKWISE)
    
    #blur = cv2.pyrMeanShiftFiltering(img, 11, 21)
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    thresh = cv2.threshold(gray, 0, 255, cv2.THRESH_BINARY_INV + cv2.THRESH_OTSU)[1]

    cnts = cv2.findContours(thresh, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)
    cnts = cnts[0] if len(cnts) == 2 else cnts[1]
    container_peri = 0

    for c in cnts:
        peri = cv2.arcLength(c, True)
        approx = cv2.approxPolyDP(c, 0.015 * peri, True)
        if len(approx) == 4:
            if container_peri < peri :
                container_peri = peri
                container_approx = approx
    #         x,y,w,h = cv2.boundingRect(approx)
    #         cv2.rectangle(image,(x,y),(x+w,y+h),(36,255,12),2)
    
    #print(container_peri)
    #print(container_approx)
    points = []
    for approx in container_approx.tolist():
        if resize:
            for i,_ in enumerate(approx[0]):
                approx[0][i] = approx[0][i]*100/scale_percent
        points.append(approx[0])
    list = {'points': points}
    jsonString = json.dumps(list)
    if view:
        cv2.drawContours(img, container_approx, -1, (0, 0, 255), 5)
        cv2.imshow('image', img)
        #cv2.imshow('thresh', thresh)
        #cv2.imshow('blur', blur)
        cv2.waitKey()
    

    return jsonString
#print(rasterFeaturePoints('D:/01 CODE/GO/georeferensi-otomatis-desktop/example/raster-scanned-maps/64710110040011.jpg',True,0,True))
# print(rasterFeaturePoints("64710100060058 - rotateRight.jpg",True,0,True))