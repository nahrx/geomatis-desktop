
# geomatis-desktop
This application is used to georeference WSS, WS, and WB maps resulting from survey/census activities conducted by BPS-Statistics Indonesia. The application generates world files (i.e. generate .jgw files for .jpg raster maps) containing georeferencing information for the WSS/WS/WB raster maps.

Up to now, georeferencing at BPS has been done manually, one map at a time, using QGIS. This method is time-consuming, as georeferencing a single map takes around 2–5 minutes, while the number of maps to be georeferenced can reach thousands, corresponding to the number of local environment unit or census blocks in a regency/city. Therefore, this application was developed to make the georeferencing process faster and more efficient. In just 1–2 minutes, thousands of maps can be automatically georeferenced.

## Distribution
check out our ready to use geomatis-desktop [installer here](https://github.com/nahrx/geomatis-desktop/releases)

## Installation for Development
The app shells out to Python for raster feature-point detection (`util.GetRasterFeaturePoints`). Which Python it uses is decided automatically at compile time via a Go build tag, so nothing needs to be hardcoded or toggled by hand:
- `wails dev` automatically compiles with the `dev` tag, so the app looks up `python`/`python3` on the system `PATH` (see `util/python_dev.go`).
- `wails build` (and `go build` without `-tags dev`) do **not** set the `dev` tag, so the app instead uses the Python runtime bundled next to the executable, at `<app directory>/python-embed/python.exe` (see `util/python_prod.go`). This is what gets shipped to end users, who don't need Python installed at all.

Setup steps:
1. Install Go
2. rename ```config.development.json``` into ```config.json```
3. For `wails dev`: install Python 3 on your machine, and make sure it's on `PATH` as `python` or `python3`, then install the libraries `pypy.py` needs:
	```
	pip install opencv-python numpy
	```
4. For `wails build` / `wails build -nsis` (producing a build to distribute): set up the embedded Python runtime once per machine, since `build/bin` is gitignored:
	1. Download the Windows "embeddable package" for Python 3.13 (64-bit or 32-bit, matching the target build) from [python.org](https://www.python.org/downloads/windows/) and extract it into `build/bin/python-embed`.
	2. In `build/bin/python-embed/python313._pth`, uncomment `import site` (this is required for `pip`/`site-packages` to work).
	3. Bootstrap `pip`:
		```
		cd build/bin/python-embed
		curl -o get-pip.py https://bootstrap.pypa.io/get-pip.py
		./python.exe get-pip.py
		```
	4. Install the required libraries:
		```
		./python.exe -m pip install opencv-python numpy
		```
	5. Copy `pypy.py` (project root) into `build/bin/python-embed/pypy.py`. Note: the `.` entry in `python313._pth` resolves to the folder containing `python.exe`, not the process's working directory, so `pypy.py` must live inside `python-embed/`, not just next to the app's `.exe`.
5. Run the application using wails command below
	```wails dev```
	or
	```wails build``` to compile the code into the `build/bin` directory (`python-embed`, set up above, must already be present there).
	or
	```wails build -nsis``` to additionally package everything (the executable + `python-embed`) into a single Windows installer at `build/bin/Geomatis-amd64-installer.exe`. This requires [NSIS](https://nsis.sourceforge.io/) to be installed (`winget install NSIS.NSIS`).

## Quick Guide
1. Go to Georeference Page\
	<img src="/example/images/img1.png" alt="This is a georeference page." style="width:400px;"/>
2. Example of georeferencing maps, using master map GeoJSON from a local directory. [example here](https://github.com/nahrx/geomatis-desktop/example)\
	<img src="/example/images/img2.png" alt="Process of georeferencing maps" style="width:400px;"/>
3. Select the map type, then select raster files. The raster file name must begin with the key attribute matching the selected map type: **IDSUBSLS** (16 digits) for WSS maps, or **IDSLS** (14 digits) for WS maps, or **IDBS** (14 digits) for WB maps -- for example: `6403110001000100.jpg` (WSS), `64710500010001_WS.jpg` (WS), `64710500010001_WB.jpg` (WB). The program will take the first N digits of the file name (matching the key's digit count above) to match it with the master map. Furthermore, the scanned raster file must be in good quality, with no folded paper, especially in the map container area, as this is the part read by the computer vision program. The **Margin** field controls the blank space percentage between the map container's corner markers and the actual map content (default `0.05` = 5%); adjust it if the georeferenced result comes out slightly off-scale.
4. If successful, a log like the following will appear:\
	<img src="/example/images/img3.png" alt="Georeference log result" style="width:400px;"/>
5. And a world file (e.g. .jgw) will be created in the same directory where the raster map file is stored.\
	<img src="/example/images/img4.png" alt="world files" style="width:500px;"/>
6. The results can be checked in QGIS – tested and verified on QGIS version 3.34.14.\
	<img src="/example/images/img5.png" alt="result in QGis" style="width:500px;"/>
