let stationsList = document.getElementById('stations-list')
let stationDetails = document.getElementById('station-details-text')

export const getStations = () => {
    fetch("/stations")
    .then((response) => response.json())
    .then((stations) => {
        showAllStations(stations)
    })
}

const showAllStations = (stations) => {
  const itemsPerPage = 20;
  let currentPage = 1;

  const renderStations = () => {
    const startIndex = (currentPage - 1) * itemsPerPage;
    const endIndex = startIndex + itemsPerPage;
    const stationsToRender = stations.slice(startIndex, endIndex);

    const stationElements = stationsToRender.map(station => {
      const stationDiv = document.createElement("div");
      stationDiv.innerHTML = `${station.Nimi}<br>${station.Namn}`;
      stationDiv.id = station.ID;
      stationDiv.className = "station-names";
      stationDiv.addEventListener("click", () => {
        showSingleStation({detail: {id: station.ID}});
      });
      return stationDiv;
    });

    stationsList.append(...stationElements);

    if (currentPage === 1 && startIndex === 0) {
      showSingleStation({detail: {id: stations[0].ID}});
    }
  };

  renderStations();

  const handleScroll = () => {
    const scrollTop = stationsList.scrollTop;
    const scrollHeight = stationsList.scrollHeight;
    const clientHeight = stationsList.clientHeight;

    if (scrollTop + clientHeight >= scrollHeight - 50) {
      currentPage += 1;
      renderStations();
    }
  };

  stationsList.addEventListener("scroll", handleScroll);
};

const renderMap = async (stationData) => {
  try {
    const { Latitude: latitude, Longitude: longitude } = stationData;

    mapboxgl.accessToken =
      "INSERT_MAPBOX_API_KEY_HERE";
    const mapElement = document.getElementById("station-details-map");
    const map = new mapboxgl.Map({
      container: mapElement,
      style: "./src/style.json",
      center: [latitude, longitude],
      zoom: 16,
    });

    map.on('styleimagemissing', function(e) {
        if (e.id === 'parking-paid') {
          map.addImage('parking-paid', {
            "width": 32,
            "height": 32,
            "data": new Uint8Array(4 * 32 * 32).fill(255, 0, 4 * 32 * 32) // Set all pixels to white
          });
        }
      });

    const marker = new mapboxgl.Marker()
      .setLngLat([latitude, longitude])
      .addTo(map);

    const popup = new mapboxgl.Popup().setHTML(
      `<h3>${stationData["Name"]}</h3><p>${stationData["Osoite"]}</p>`
    );

    marker.setPopup(popup).togglePopup();
  } catch (error) {
    console.error("Error rendering map:", error);
  }
};

const getStationData = async (id) => {
  try {
    const response = await fetch(`/stations?id=${id}`);
    const stationData = await response.json();
    return stationData;
  } catch (error) {
    console.error("Error fetching station data:", error);
  }
};

const showSingleStation = async (event) => {
  try {
    const stationData = await getStationData(event.detail.id);
    stationDetails.innerHTML = "";

    stationDetails.innerHTML = `<h2>${stationData["Nimi"]},
    ${stationData["Namn"]}</h2><br>
    <i class="fas fa-map-marker-alt"></i> Bike station located at ${stationData["Osoite"]}, ${stationData["Adress"]}.<br>
    <i class="fas fa-arrow-circle-up"></i> ${stationData["JourneysFrom"]} journeys began here, while<br>
    <i class="fas fa-arrow-circle-down"></i> ${stationData["JourneysTo"]} journeys came to an end here.`;

    await renderMap(stationData);

    const stationElements = document.querySelectorAll(".station-names");
    stationElements.forEach((element) => {
      if (parseInt(element.id) === event.detail.id) {
        element.classList.add("selected");
      } else {
        element.classList.remove("selected");
      }
    });
  } catch (error) {
    console.error("Error showing single station:", error);
  }
};

const handleWheelEvent = (event) => {
  event.preventDefault();
  stationsList.scrollLeft += event.deltaY + event.deltaX;
}

const mediaQuery = window.matchMedia("(max-width: 590px)");

if (mediaQuery.matches) {
  stationsList.addEventListener("wheel", handleWheelEvent);
} else {
  stationsList.removeEventListener("wheel", handleWheelEvent);
}

mediaQuery.addEventListener("change", (event) => {
  if (event.matches) {
    stationsList.addEventListener("wheel", handleWheelEvent);
  } else {
    stationsList.removeEventListener("wheel", handleWheelEvent);
  }
});