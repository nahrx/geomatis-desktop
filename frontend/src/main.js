import './style.css';
import logo from './assets/images/logo-universal.png';
import { SelectRasterFiles, 
  SelectGeojsonFileForGeoreference, 
  ProcessGeoreference, 
  ConnectToDB, 
  DisconnectDB, 
  SaveDbConfig, 
  LoadDbConfig, 
  GetMasterMaps, 
  SelectGeojsonFile, 
  CreateMasterMaps, 
  DeleteMasterMap } 
  from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
<main class="container" role="main">
  <img id="logo" class="logo">
  <h1>Geomatis</h1>
  <p class="text-center" style="margin-bottom:20px">Automatically georeference WS maps from BPS survey/census activities.</p>
  <nav class="tabs" role="tablist" aria-label="Main Sections">
    <div class="tab public active" role="tab" tabindex="-1" aria-selected="false" aria-controls="georeferensi" id="tab-georeferensi" data-tab="georeferensi">Georeference</div>
    <div class="tab public" role="tab" tabindex="0" aria-selected="true" aria-controls="database" id="tab-database" data-tab="database">DB Connection</div>
    <div class="tab private disabled" role="tab" tabindex="-1" aria-selected="false" aria-controls="master" id="tab-master" data-tab="master">DB Master Maps</div>
    <div class="tab public" role="tab" tabindex="-1" aria-selected="false" aria-controls="about" id="tab-about" data-tab="about">About</div>
  </nav>
  <!-- Georeferensi Tab -->
  <section class="input-box active" role="tabpanel" aria-labelledby="tab-georeferensi" id="georeferensi" tabindex="0" hidden>
    <label style="display: block;margin-top:5px"> Choose Map Type </label>
    <div class="map-option">
      <input type="radio" id="ws-type" name="masterMapType" value="ws" checked="checked">
      <div>
        <label for="ws-type">WS Map</label>
      </div>
      <input type="radio" id="wb-type" name="masterMapType" value="wb">
      <div>
        <label for="wb-type">WB Map</label>
      </div>
    </div>
    <label style="display: block;margin-top:5px"> Choose Master Map Source </label>
    <div class="source-option">
      <input type="radio" id="file-source" name="masterMapSource" value="file" checked="checked">
      <div>
        <label for="file-source">From File</label>
      </div>
      <input type="radio" id="database-source" name="masterMapSource" value="database">
      <div>
        <label for="database-source">From Database</label>
      </div>
    </div>
    <div class="file-source-input">
      <label style="display: block; margin-bottom: 15px; margin-top:5px"> Select Master Map </label>
      <button class="btn upload-btn w-full" id="selectGeojsonButton" onclick="selectGeojson()">
        <span class="material-icons">upload</span> Select a geojson file </button>
      <div class="select-files" id="geojsonFile" style="padding:0"> No geojson selected </div>
    </div>
    <div class="database-source-input">
      <label style="display: block; margin-bottom: 15px; margin-top:5px"> Select Master Map </label>
      <div style="display: flex; align-items: center; gap: 8px;">
        <select id="masterMapSelect" class="input" aria-describedby="mapSelectDesc" aria-required="true" style="flex-grow: 1;">
          <option id="mapSelectDesc">Database not connected</option>
        </select>
        <button class="btn btn-refresh" id="refreshMasterMapsSelect" onclick="refreshMasterMapsSelect()" title="Refresh Master Maps" style="margin:2px">
          <span class="material-icons" aria-hidden="true">refresh</span>
        </button>
      </div>
    </div>
    <div>
      <label style="display: block; margin-bottom: 15px; margin-top:5px"> Select Raster Files </label>
      <button class="btn upload-btn w-full" id="selectFileButton" onclick="selectRasterFiles()">
        <span class="material-icons">upload</span> Select raster files (jpg/png) </button>
      <div class="select-files" id="rasterFiles"> No raster selected </div>
    </div>
    <button class="btn w-full" id="submitFilesButton" onclick="processGeoreference()">
      <span class="material-icons">play_arrow</span> Start Georeference </button>
  </section>
  <!-- Database Setting Tab -->
  <section class="input-box " role="tabpanel" aria-labelledby="tab-database" id="database" tabindex="0">
    <p>Setup and connect to postgresql database. instead of uploading the master map from file, we can use the master map from database.</p>
    <label for="databaseHost">Host</label>
    <input id="databaseHost" type="text" value="localhost" />
    <label for="databasePort">Database Port</label>
    <input id="databasePort" type="text" value="5432" />
    <label for="databaseName">Database Name</label>
    <input id="databaseName" type="text" />
    <label for="databaseUsername">User Name</label>
    <input id="databaseUsername" type="text" />
    <label for="databasePassword">Password</label>
    <input id="databasePassword" type="password" />
    <div style="display:flex; gap: 12px; margin-top:12px;">
      <button type="button" class="btn" id="connectDbBtn" onclick="connectDatabase()">
        <span class="material-icons" aria-hidden="true">link</span>Connect to Database </button>
      <button type="button" class="btn warning-btn" id="disconnectDbBtn" style="display: none;" onclick="disconnectDatabase()">
        <span class="material-icons" aria-hidden="true">link_off</span>Disconnect Database </button>
      <button type="button" class="save-btn btn" id="saveDbBtn" onclick="saveDatabaseConfig()">
        <span class="material-icons" aria-hidden="true">save</span>Save Settings </button>
    </div>
  </section>
  <!-- Master Maps Tab -->
  <section class="input-box" role="tabpanel" aria-labelledby="tab-master" id="master" tabindex="0" hidden>
    <p>Manage master maps from database</p>
    <div class="button-bar" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem;">
      <button class="warning-btn  btn" type="button" id="uploadMasterMapBtn" onclick="uploadMasterMap()">
        <span class="material-icons" aria-hidden="true">upload_file</span>Upload Master Maps </button>
      <button class="refresh-btn btn" type="button" id="refreshMasterMapBtn" onclick="refreshMasterMapsTable()">
        <span class="material-icons" aria-hidden="true">refresh</span>Refresh </button>
    </div>
    <div class="table-container" role="region" aria-label="Master Maps List">
      <table id="masterMapsTable">
        <thead>
          <tr>
            <th scope="col">Name</th>
            <th scope="col">Dimension</th>
            <th scope="col">SRID</th>
            <th scope="col">Type</th>
            <th scope="col" aria-label="Delete action"></th>
          </tr>
        </thead>
        <tbody>
          <!-- Rows will be dynamically added -->
        </tbody>
      </table>
    </div>
  </section>
  <!-- About Tab -->
  <section class="input-box" role="tabpanel" aria-labelledby="tab-about" id="about" tabindex="0" hidden>
    <h3>Geomatis-Desktop Version 1.0</h3>
    <h4>About</h4>
    <p>Geomatis is used for automatic georeferencing of WS maps resulting from BPS survey activities. This app will create world files (.jwg for .jpg image / .pgw for .png image) to store georeferencing information for raster images.</p>
    <h4>Requirements</h4>
    <ul>
      <li> The scanned map file name must begin with <strong>IDSLS</strong>, for example: <code>64710500010001.jpg</code>, <code>64710500010001_WS.jpg</code>. The program will take the first 14 digits of the file name to match it with the IDSLS in the digital SLS master. </li>
      <li> The WS map scan result must be of good quality, with no folded paper, especially in the map container area, as this is the part read by the computer vision program. </li>
      <li> The WS map scan must not be upside down. </li>
    </ul>
    <h4>Quick Guide</h4>
    <h5>From file</h5>
    <ol>
      <li>Go to the <strong>Georeference</strong> tab. </li>
      <li>Select map type <strong>(WS/WB map)</strong>
      </li>
      <li>Click <strong>From File</strong> button </li>
      <li>Upload the master map geojson file from computer</li>
      <li>Upload the map rasters to be georeferenced from computer</li>
      <li>Start the georeferencing process.</li>
    </ol>
    <h5>From database</h5>
    <ol>
      <li>Connect to the PostgreSQL database that has been previously created ( <a href="https://www.w3schools.com/postgresql/postgresql_install.php" target="_blank">tutorial guide</a>) </li>
      <li>Add and manage SLS master data in the <strong>Master Maps</strong> tab </li>
      <li>Go to the <strong>Georeference</strong> tab </li>
      <li>Select map type <strong>(WS/WB map)</strong>
      </li>
      <li>Click <strong>From Database</strong> button </li>
      <li>Upload the master map available in the database</li>
      <li>Upload the map rasters to be georeferenced from computer</li>
      <li>Start the georeferencing process</li>
    </ol>
    <h4>Contact Programmer</h4>
    <p>
      <b>nahar.nasrullah@bps.go.id</b>
    </p>
    <p>BPS Kota Balikpapan</p>
  </section>
  <div>
    <div class="log-result" id="log" aria-live="polite" aria-atomic="true"></div>
  </div>
</main>
`;
document.getElementById('logo').src = logo;
const masterMapSelect = document.getElementById("masterMapSelect");
const masterMapsTableBody = document.querySelector('#masterMapsTable tbody');

const tabs = document.querySelectorAll('.tab');
const tabsPrivate = document.querySelectorAll('.tab.private');
const tabsPublic = document.querySelectorAll('.tab.public');
const panels = document.querySelectorAll('.input-box');

var rasterFiles = document.getElementById('rasterFiles');
var geojsonFile = document.getElementById('geojsonFile');
var log = document.getElementById('log');
let tabHandlers = new Map();

var selectedRasterFiles = [];
var selectedGeojsonFile = "";


// Database ///////////////////////////////

document.addEventListener("DOMContentLoaded", () => {
  LoadDbConfig().then((config) => {
    document.getElementById("databaseHost").value = config.DB_HOST;
    document.getElementById("databasePort").value = config.DB_PORT;
    document.getElementById("databaseName").value = config.DB_DATABASE;
    document.getElementById("databaseUsername").value = config.DB_USERNAME;
    document.getElementById("databasePassword").value = config.DB_PASSWORD;
  }).catch((err) => {
    log.innerHTML += "\nFailed to load DB config:" + err;
  });
});

window.saveDatabaseConfig = function () {
  log.innerHTML = 'saving database config...';
  try {
    const config = getDbConfig();

    SaveDbConfig(config).then((map) => {
      log.innerHTML = `Configuration saved! ("${config.DB_DATABASE}" at ${config.DB_HOST} as ${config.DB_USERNAME}.)`;
    }).catch((error) => {
      log.innerHTML = "Error saving config: " + error;
    })
  } catch (error) {
    log.innerHTML = "Error saving config: " + error;
  }
};
window.connectDatabase = function () {
  log.innerHTML = 'Connecting to database...';
  const config = getDbConfig();
  ConnectToDB(config).then((maps) => {
    log.innerHTML = `Connected to database "${config.DB_DATABASE}" at ${config.DB_HOST} as ${config.DB_USERNAME}.`;
    enableMenu();
  }).catch((err) => {
    log.innerHTML = "Error connect to database: " + error;
  })
}
window.disconnectDatabase = function () {

  log.innerHTML = 'Disconnecting from database...';
  const config = getDbConfig();
  DisconnectDB(config).then((maps) => {
    log.innerHTML = `Disconnect database "${config.DB_DATABASE}" at ${config.DB_HOST} as ${config.DB_USERNAME}.`;
    disableMenu();
  }).catch((err) => {
    log.innerHTML = "Error disconnect database: " + error;
  })
}

function getDbConfig() {
  const dbHost = document.getElementById("databaseHost").value.trim();
  const dbPort = document.getElementById("databasePort").value.trim();
  const dbName = document.getElementById("databaseName").value.trim();
  const dbUser = document.getElementById("databaseUsername").value.trim();
  const dbPass = document.getElementById("databasePassword").value;

  const config = {
    DB_HOST: dbHost,
    DB_PORT: parseInt(dbPort),
    DB_DATABASE: dbName,
    DB_USERNAME: dbUser,
    DB_PASSWORD: dbPass
  };
  return config
}

function enableMenu() {
  document.getElementById('connectDbBtn').style.display = 'none';
  document.getElementById('disconnectDbBtn').style.display = 'inline-flex';
  refreshMasterMapsTable();
  disableDatabaseForm();
  enableTabListeners(tabs);
}

function disableMenu() {
  document.getElementById('connectDbBtn').style.display = 'inline-flex';
  document.getElementById('disconnectDbBtn').style.display = 'none';
  enableDatabaseForm();
  disableTabListeners(tabsPrivate);
}
function enableDatabaseForm() {
  const inputs = document.getElementById("database").querySelectorAll("input");
  inputs.forEach(input => {
    input.disabled = false;
  });
}
function disableDatabaseForm() {
  const inputs = document.getElementById("database").querySelectorAll("input");
  inputs.forEach(input => {
    input.disabled = true;
  });

}

window.refreshMasterMapsSelect = function () {
  masterMapSelect.innerHTML = '<option disabled>Select Database...</option>'; // Clear existing "Loading..." option
  GetMasterMaps().then((masterMaps) => {
    if (masterMaps == null) {
      const option = document.createElement('option');
      option.textContent = "there is no data";
      masterMapSelect.appendChild(option);
    }
    masterMaps.forEach((map) => {
      const option = document.createElement('option');
      option.value = map.name;
      option.textContent = map.name;
      masterMapSelect.appendChild(option);
    });
  }).catch((err) => {
    log.innerHTML += "\nerror during load data : " + err;
  })
}

window.refreshMasterMapsTable = function () {
  GetMasterMaps().then((masterMaps) => {
    masterMapsTableBody.innerHTML = '';
    masterMapSelect.innerHTML = '<option disabled>Select Database...</option>';
    if (masterMaps == null) {
      const tr = document.createElement('tr');
      const td = document.createElement('td');
      td.colSpan = 5;
      td.style.textAlign = 'center';
      td.style.color = '#94a3b8';
      td.textContent = 'No master maps available.';
      tr.appendChild(td);
      masterMapsTableBody.appendChild(tr);
      log.innerHTML += "\nThere is no master maps yet";
      return;
    }
    masterMaps.forEach((map, i) => {
      // select option in georeferensi menu
      const option = document.createElement('option');
      option.value = map.name;
      option.textContent = map.name;
      masterMapSelect.appendChild(option);

      // table content
      const tr = document.createElement('tr');
      tr.innerHTML = `
            <td>${map.name}</td>
            <td>${map.dimension}</td>
            <td>${map.srid}</td>
            <td>${map.type}</td>
            <td><button class="delete-btn btn" aria-label="Delete ${map.name}" data-name="${map.name}">Delete</button></td>
          `;
      masterMapsTableBody.appendChild(tr);
    });
    masterMapsTableBody.querySelectorAll('.delete-btn').forEach(btn => {
      btn.addEventListener('click', (e) => {
        const name = e.currentTarget.getAttribute('data-name');
        if (confirm(`Are you sure you want to delete "${name}"?`)) {
          deleteMasterMap(name)
        }
      });
    });
  }).catch((err) => {
    log.innerHTML += "\nFailed to load master maps:" + err;
    masterMapSelect.innerHTML = '<option disabled>Error loading maps</option>';
  });
}
window.uploadMasterMap = function () {
  SelectGeojsonFile().then((filePath) => {
    CreateMasterMaps(filePath).then((result) => {
      refreshMasterMapsTable();
      log.innerHTML += "\ngeojson file uploaded succesfully ";
    }).catch((err) => {
      log.innerHTML += "\nError upload geojson file : " + err;
    })
  }).catch((err) => {
    log.innerHTML += "\nError select geojson file : " + err;
  })
}
window.deleteMasterMap = function (name) {
  DeleteMasterMap(name).then((result) => {
    refreshMasterMapsTable();
    log.innerHTML += "\nMastermap " + name + "deleted succesfully ";
  }).catch((err) => {
    log.innerHTML += "\nError deleteting master map : " + name + err;
  })
}
function addGeojsonFile(gFile){
  geojsonFile.innerHTML = ""
  selectedGeojsonFile = gFile;
  const listItem = document.createElement('div');
  listItem.setAttribute('title', gFile);
  var filename = gFile.replace(/^.*[\\/]/, '');
  listItem.innerHTML = '<span>' + filename + '</span>' + '<button class="x"data-path=\"' + gFile + '\"onclick="removeGeojsonFile()">x</button>';
  geojsonFile.appendChild(listItem);
  geojsonFile.style.display = "block";
}

window.removeGeojsonFile = function () {
  //e.parentNode.parentNode.removeChild(e.parentNode);
  selectedGeojsonFile = "";
  geojsonFile.innerHTML = "";
  geojsonFile.style.display = "none";
}

window.selectGeojson = function () {
  const masterMapType = document.querySelector('input[name="masterMapType"]:checked').value;
  SelectGeojsonFileForGeoreference(masterMapType)
    .then((result) => {
      addGeojsonFile(result)
    })
    .catch((err) => {
      log.innerHTML += "Error select file :" + err;
    });
};

function addRasterFiles(rasters){
  if (selectedRasterFiles.length == 0) {
        rasterFiles.innerHTML = ""
      }
      rasters.forEach(function (item, index) {
        var index = selectedRasterFiles.indexOf(item);
        if (index > -1) { // only when item is found
          return;
        }
        selectedRasterFiles.push(item);
        const listItem = document.createElement('div');
        listItem.setAttribute('title', item);
        var filename = item.replace(/^.*[\\/]/, '');
        listItem.innerHTML = '<span>' + filename + '</span>' + '<button class="x"data-path=\"' + item + '\"onclick="removeRasterFiles(this)">x</button>';
        rasterFiles.appendChild(listItem);
        rasterFiles.style.display = "grid";
      });
}

window.removeRasterFiles = function (e) {
  var fpath = e.getAttribute("data-path");
  e.parentNode.parentNode.removeChild(e.parentNode);
  var index = selectedRasterFiles.indexOf(fpath);
  if (index > -1) { // only splice array when item is found
    selectedRasterFiles.splice(index, 1); // 2nd parameter means remove one item only
  }
  if (selectedRasterFiles.length == 0) {
    rasterFiles.style.display = "none";
    rasterFiles.innerHTML = "";
  }
};

window.selectRasterFiles = function () {
  SelectRasterFiles()
    .then((result) => {
      addRasterFiles(result)
    })
    .catch((err) => {
      log.innerHTML += "Error select file :" + err;
    });
};

window.processGeoreference = function () {
  var selectedMap = "";
  const masterMapType = document.querySelector('input[name="masterMapType"]:checked').value;
  const masterMapSource = document.querySelector('input[name="masterMapSource"]:checked').value;

  if (masterMapSource == "database") {
    selectedMap = document.getElementById('masterMapSelect').value;
  } else if (masterMapSource == "file") {
    selectedMap = selectedGeojsonFile;
  }

  if (selectedMap == "" || selectedMap == null) {
    log.innerHTML = "no master map selected"
    return;
  }

  if (selectedRasterFiles.length === 0) {
    log.innerHTML = "no raster file selected"
    return;
  }
  // Call App.Select(name)
  try {
    log.innerHTML = "<span class='spinner'></span>georeference loading...";
    ProcessGeoreference(selectedRasterFiles, selectedMap, masterMapType, masterMapSource)
      .then((result) => {
        resetRasterFiles();
        log.innerHTML = "";
        result.forEach(function (item, index) {
          const listItem = document.createElement('p');
          if (item.indexOf('success :') != 0) {
            listItem.setAttribute('class', 'red');
          }
          listItem.textContent = item;
          log.appendChild(listItem);
        });
      })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
    alert(error);
  }
};

function resetRasterFiles(){
  selectedRasterFiles.length = 0;
  rasterFiles.style.display = "none";
  rasterFiles.innerHTML = "";
}

function resetForm() {
  resetRasterFiles();
  removeGeojsonFile();
}

function disableTabListeners(tabs) {
  tabs.forEach(tab => {
    const handlers = tabHandlers.get(tab);
    if (handlers) {
      tab.removeEventListener('click', handlers.clickHandler);
      tab.removeEventListener('keydown', handlers.keyHandler);
      tabHandlers.delete(tab);
    }
    tab.setAttribute('aria-disabled', 'true');
    tab.classList.add('disabled');
  });
}
enableTabListeners(tabsPublic);
function enableTabListeners(tabs) {
  tabs.forEach(tab => {
    const clickHandler = () => {
      const target = tab.getAttribute('data-tab');
      switchTab(target, tab);
    };

    const keyHandler = (e) => {
      let index = Array.from(tabs).indexOf(e.target);
      if (e.key === 'ArrowRight') {
        e.preventDefault();
        let nextIndex = (index + 1) % tabs.length;
        tabs[nextIndex].focus();
        tabs[nextIndex].click();
      }
      if (e.key === 'ArrowLeft') {
        e.preventDefault();
        let prevIndex = (index - 1 + tabs.length) % tabs.length;
        tabs[prevIndex].focus();
        tabs[prevIndex].click();
      }
    };

    tab.addEventListener('click', clickHandler);
    tab.addEventListener('keydown', keyHandler);

    tabHandlers.set(tab, { clickHandler, keyHandler });
    // Re-enable visually
    tab.removeAttribute('aria-disabled');
    tab.classList.remove('disabled');
  });
}

function switchTab(tabName, tabElement) {
  tabs.forEach(t => {
    t.classList.remove('active');
    t.setAttribute('aria-selected', 'false');
    t.setAttribute('tabindex', '-1');
  });
  panels.forEach(panel => {
    panel.classList.remove('active');
    panel.setAttribute('hidden', '');
  });

  tabElement.classList.add('active');
  tabElement.setAttribute('aria-selected', 'true');
  tabElement.setAttribute('tabindex', '0');

  const panel = document.getElementById(tabName);
  panel.classList.add('active');
  panel.removeAttribute('hidden');
  panel.focus({ preventScroll: true });
}

const databaseSource = document.getElementById('database-source');
const fileSource = document.getElementById('file-source');
const fileSourceInput = document.querySelector('.file-source-input');
const databaseSourceInput = document.querySelector('.database-source-input');

// Event listener for radio buttons
databaseSource.addEventListener('change', function () {
  if (this.checked) {
    fileSourceInput.style.display = 'none';
    databaseSourceInput.style.display = 'block';
    removeGeojsonFile();
    masterMapSelect.innerHTML = '<option disabled>Select Database...</option>';
  }
});

fileSource.addEventListener('change', function () {
  if (this.checked) {
    databaseSourceInput.style.display = 'none';
    fileSourceInput.style.display = 'block';
    removeGeojsonFile();
    masterMapSelect.innerHTML = '<option disabled>Select Database...</option>';
  }
});

databaseSourceInput.style.display = 'none';
fileSourceInput.style.display = 'block';

document.querySelectorAll('input[name="masterMapType"]').forEach(function (el) {
  el.addEventListener('change', function (event) {
    removeGeojsonFile();
    masterMapSelect.innerHTML = '<option disabled>Select Database...</option>';
  });
});



