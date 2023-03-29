let currentPage = 0;
let stationsSection = document.getElementById('stations-list')

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

    stationsSection.innerHTML = ""; // clear the previous results

    stationsSlice.forEach((station) => {
        let stationDiv = document.createElement('div')
        stationDiv.innerHTML = `${station.Nimi} / ${station.Namn}`
        stationDiv.id = station.ID
        stationDiv.addEventListener('click', showSingleStation)
        stationsSection.appendChild(stationDiv)
    });

    currentPage = page;

    const pageCount = Math.ceil(stations.length / 10);

    const count = document.createElement("div")
    count.innerHTML = `Page ${page + 1} of ${pageCount}`
    stationsSection.appendChild(count)

    if (stationsSlice.length === 10) {
        const nextButton = document.createElement("button");
        nextButton.textContent = "Next";
        nextButton.addEventListener("click", () => {
            showAllStations(stations, currentPage + 1);
        });
        stationsSection.appendChild(nextButton);
    }

    if (currentPage > 0){
        const prevButton = document.createElement("button");
        prevButton.textContent = "Previous";
        prevButton.addEventListener("click", () => {
            showAllStations(stations, currentPage - 1);
        });
        stationsSection.appendChild(prevButton)
    }
    
}

const showSingleStation = (event) => {
    fetch("/stations?id=" + event.target.id)
    .then((response) => response.json())
    .then((station) => {
        stationsSection.innerHTML = ""
        for (const key in station[0]){
            if (key == 'JourneysFrom' || key == 'JourneysTo'){
                stationsSection.innerHTML += station[0][key]["String"] + "<br>"
            } else {
                stationsSection.innerHTML += station[0][key] + "<br>"
            }
        }
    })
}