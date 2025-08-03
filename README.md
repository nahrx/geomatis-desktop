
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
1. Go to Georeference Page\n
	![Georeference page](/example/images/img1.png "This is a georeference page.")
2. Example of georeferencing WS maps, using master map GeoJSON from a local directory. [example here](https://github.com/nahrx/geomatis-desktop/example)\n
	![Georeferencing WS maps](/example/images/img2.png "Process of georeferencing WS maps" )
3. If successful, a log like the following will appear:\n
	![Georeference log](/example/images/img3.png "Georeference log result")
4. And a world file (e.g. .jgw) will be created in the same directory where the raster map file is stored.\n
	![World file](/example/images/img4.png "world file")
5. The results can be checked in QGIS – tested and verified on QGIS version 3.34.14.\n
	![QGis result](/example/images/img5.png "Result in QGis")

