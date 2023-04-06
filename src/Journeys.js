let journeysList = document.getElementById('journeys-list')
let journeyBatch = 0
let currentPage = 0
let fetchingPreviousBatch = false

const addZero = (num) => num < 10 ? `0${num}` : num

export const getJourneys = (e) => {
    e.preventDefault()
    fetch(`/journeys?lastJourneyId=${journeyBatch}`)
    .then((response) => response.json())
    .then((journeys) => {
        if (fetchingPreviousBatch){
            fetchingPreviousBatch = false
            showJourneysBatch(journeys, 29)
        } else {
            showJourneysBatch(journeys, 0)
        }
    })
}

const showJourneysBatch = (journeys, page) => {
    const start = page * 100;
    const end = start + 100;
    const journeysSlice = journeys.slice(start, end)
    
    journeysList.innerHTML = "";
    let journeysTable = document.createElement("TABLE");
        journeysTable.innerHTML = `<tr>
        <td>Departure</td>
        <td>From</td>
        <td>To</td>
        <td>Distance</td>
        <td>Duration</td>
        </tr>`
        journeysSlice.forEach(journey => {
            let hours = 0
            let minutes = 0
            let seconds = 0
            const dateString = `${journey.Departure}Z`
            const date = new Date(dateString)
            const formattedDate = `${addZero(date.getUTCDate())}.${addZero(date.getUTCMonth() + 1)}.${date.getUTCFullYear()} ${addZero(date.getUTCHours())}.${addZero(date.getUTCMinutes())}.${addZero(date.getUTCSeconds())}`;

            if (journey.DurationSec > 3600){
                hours = Math.floor(journey.DurationSec / 3600)
                minutes = Math.floor((journey.DurationSec - hours * 3600)/60)
            } else {
                minutes = Math.floor(journey.DurationSec / 60)
                seconds = journey.DurationSec - minutes * 60;
            }
            const km = (journey.CoveredDistanceM / 1000).toFixed(1)
            journeysTable.innerHTML += `<tr>
            <td>${formattedDate}</td>
            <td>${journey.DepartureStationName}</td>
            <td>${journey.ReturnStationName}</td>
            <td>${km}km</td>
            <td>${hours > 0 ? `${hours}h` : ""}
            ${minutes > 0 ? `${minutes}min` : ""}
            ${seconds > 0 ? `${seconds}s` : ""}</td>
            </tr>`
        });
        journeysList.appendChild(journeysTable)

        currentPage = page
        const pageCount = Math.ceil(journeys.length / 100);
        const count = document.createElement('div')
        count.innerHTML = `Showing ${journeys[start].ID}–${journeys[end-1].ID} of XXX journeys`
        journeysList.appendChild(count)

        if (journeysSlice.length === 100) {
            const nextButton = document.createElement("button");
            nextButton.textContent = "Next";
            if (pageCount == currentPage+1){
                nextButton.addEventListener("click", (e) => {
                    journeyBatch += 3000
                    getJourneys(e)
                });
            } else {
                nextButton.addEventListener("click", () => {
                    showJourneysBatch(journeys, currentPage + 1)
                });
            }
            
            journeysList.appendChild(nextButton);
        }

        if (currentPage > 0 || journeyBatch > 0){
            const prevButton = document.createElement("button");
            prevButton.textContent = "Previous";
            if (journeyBatch > 0 && currentPage == 0){
                prevButton.addEventListener("click", (e) => {
                    journeyBatch -= 3000
                    fetchingPreviousBatch = true
                    getJourneys(e);
                });
            } else {
                prevButton.addEventListener("click", () => {
                    showJourneysBatch(journeys, currentPage - 1);
                });
            }
            journeysList.appendChild(prevButton)
        }
}
