let stationsList = document.getElementById('stations-list')
let stationDetails = document.getElementById('station-details')
let currentPage = 0;

export const getStations = (event) => {
    event.preventDefault()
    fetch("/stations")
    .then((response) => response.json())
    .then((stations) => {
        showAllStations(stations, 0)
    })

}

const showAllStations = (stations, page) => {
    const start = page * 10;
    const end = start + 10;

    const stationsSlice = stations.slice(start, end);

    stationsList.innerHTML = ""; // clear the previous results

    stationsSlice.forEach((station) => {
        let stationDiv = document.createElement('div')
        stationDiv.innerHTML = `${station.Nimi} / ${station.Namn}`
        stationDiv.id = station.ID
        stationDiv.addEventListener('click', showSingleStation)
        stationsList.appendChild(stationDiv)
    });

    currentPage = page;

    const pageCount = Math.ceil(stations.length / 10);

    const count = document.createElement("div")
    count.innerHTML = `Page ${page + 1} of ${pageCount}`
    stationsList.appendChild(count)

    if (stationsSlice.length === 10) {
        const nextButton = document.createElement("button");
        nextButton.textContent = "Next";
        nextButton.addEventListener("click", () => {
            showAllStations(stations, currentPage + 1);
        });
        stationsList.appendChild(nextButton);
    }

    if (currentPage > 0){
        const prevButton = document.createElement("button");
        prevButton.textContent = "Previous";
        prevButton.addEventListener("click", () => {
            showAllStations(stations, currentPage - 1);
        });
        stationsList.appendChild(prevButton)
    }
    
}

// const isJourneyCount = (key) => key === 'JourneysFrom' || key === 'JourneysTo'

const showSingleStation = (event) => {
    fetch(`/stations?id=${event.target.id}`)
        .then((response) => response.json())
        .then((station) => {
        const [stationData] = station;
        stationDetails.innerHTML = "";
        for (const key in stationData) {
            stationDetails.innerHTML += `${stationData[key]}<br>`;
        }
        });
};