let journeysList = document.getElementById('journeys-list')
let journeyBatch = 0
let currentPage = 0
let fetchingPreviousBatch = false

const addZero = (num) => num < 10 ? `0${num}` : num

export const getJourneys = () => {
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

let journeysTable = document.createElement("TABLE");
journeysTable.cellSpacing = 0;

const tableHeader = document.createElement("thead");
const headerRow = document.createElement("tr");

headerRow.innerHTML = `
    <th><h4>Departure</h4></th>
    <th><h4>From</h4></th>
    <th><h4>To</h4></th>
    <th><h4>Distance</h4></th>
    <th><h4>Duration</h4></th>
`;

tableHeader.appendChild(headerRow);
journeysTable.appendChild(tableHeader);
journeysList.appendChild(journeysTable)

const showJourneysBatch = (journeys, page) => {
    const start = page * 50;
    const end = start + 50;
    const journeysSlice = journeys.slice(start, end)

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

        const tableRow = document.createElement("tr");
        tableRow.innerHTML = `
            <tr>
                <td>${formattedDate}</td>
                <td>${journey.DepartureStationName}</td>
                <td>${journey.ReturnStationName}</td>
                <td>${km}km</td>
                <td>
                ${hours > 0 ? `${hours}h` : ""}
                ${minutes > 0 ? `${minutes}min` : ""}
                ${seconds > 0 ? `${seconds}s` : ""}
                </td>
            </tr>
        `
        journeysTable.appendChild(tableRow)
    });

    currentPage = page;
    const pageCount = Math.ceil(journeys.length / 50);

    if (journeysSlice.length === 50) {
        let hasReachedEnd = false
        const handleScroll = () => {
            if (hasReachedEnd) {
                return;
            }
            const { scrollTop, scrollHeight, clientHeight } = journeysList;
            if (scrollTop + clientHeight >= scrollHeight/40) {
                hasReachedEnd = true;
                if (pageCount == currentPage+1) {
                    journeyBatch += 3000;
                    getJourneys();
                } else {
                    showJourneysBatch(journeys, currentPage + 1)
                }
                
                journeysList.removeEventListener('scroll', handleScroll);
            }
        };
        journeysList.addEventListener('scroll', handleScroll);
    }
};