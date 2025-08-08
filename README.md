
# geomatis-desktop
This application is used to georeference  WS and WB maps resulting from survey/census activities conducted by BPS (Statistics Indonesia). The application generates world files (i.e., .jgw files for .jpg raster maps) containing georeferencing information for the WS/WB raster maps.

Up to now, georeferencing at BPS has been done manually, one map at a time, using QGIS. This method is time-consuming, as georeferencing a single map takes around 2–5 minutes, while the number of maps to be georeferenced can reach thousands, corresponding to the number of local environment unit or census blocks in a regency/city. Therefore, this application was developed to make the georeferencing process faster and more efficient. In just 1–2 minutes, thousands of maps can be automatically georeferenced.

## Distribution
check out our ready to use geomatis-desktop [installer here](https://github.com/nahrx/geomatis-desktop/releases)

## Installation
1. Install python 3
2. Install PIP
3. Install opencv library using this command
	``` pip install opencv-python```
4. Install Go 
6. rename ```config.development.json``` into ```config.json```
7. Run the application using wails
	```wails dev```
	or
	```wails build``` for production

## Quick Guide
1. Go to Georeference Page\
	<img src="/example/images/img1.png" alt="This is a georeference page." style="width:400px;"/>
2. Example of georeferencing WS maps, using master map GeoJSON from a local directory. [example here](https://github.com/nahrx/geomatis-desktop/example)\
	<img src="/example/images/img2.png" alt="Process of georeferencing WS maps" style="width:400px;"/>
3. Select raster files. The raster files name must begin with IDSLS attribute from local environment unit master map, or begin with IDBS from block cencus master map. for example: 64710500010001.jpg, 64710500010001_WS.jpg. The program will take the first 14 digits of the file name to match it with the IDSLS/IDBS in the digital master map. Furthermore, The scanned raster file must be in good quality, with no folded paper, especially in the map container area, as this is the part read by the computer vision program.
3. If successful, a log like the following will appear:\
	<img src="/example/images/img3.png" alt="Georeference log result" style="width:400px;"/>
4. And a world file (e.g. .jgw) will be created in the same directory where the raster map file is stored.\
	<img src="/example/images/img4.png" alt="world files" style="width:500px;"/>
5. The results can be checked in QGIS – tested and verified on QGIS version 3.34.14.\
	<img src="/example/images/img5.png" alt="result in QGis" style="width:500px;"/>

