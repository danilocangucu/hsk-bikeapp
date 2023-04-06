import { getStations } from './Stations.js'
import { getJourneys } from './Journeys.js'

const stations = document.getElementById('stations')
const journeys = document.getElementById('journeys')

stations.addEventListener('click', getStations)
journeys.addEventListener('click', getJourneys)
