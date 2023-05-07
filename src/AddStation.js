import { getStations } from './Stations.js'

let addStationResponse = document.getElementById('add-station-response')

export const addStation = () => {
  const form = document.getElementById("station-form");
  form.addEventListener("submit", async function (event) {
    event.preventDefault();
    addStationResponse.textContent = String.fromCharCode(160);

    const data = new FormData(form);
    const newStation = {
      ID: parseInt(data.get("ID")),
      Nimi: data.get("Nimi"),
      Namn: data.get("Namn"),
      Name: data.get("Name"),
      Osoite: data.get("Osoite"),
      Adress: data.get("Adress"),
      Operaattor: data.get("Operaattor"),
      Kapasiteet: parseInt(data.get("Kapasiteet")),
      JourneysFrom: parseInt(data.get("JourneysFrom")),
      JourneysTo: parseInt(data.get("JourneysTo")),
    };

    const result = await validateNewStation(newStation);

    addStationResponse.innerHTML = ""
    let responseDiv = document.createElement("div")
    addStationResponse.appendChild(responseDiv)

    if (result.isValid) {
      fetch("/addstation", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          ...newStation,
          Kaupunki: result.Kaupunki,
          Stad: result.Stad,
          Latitude: result.Latitude,
          Longitude: result.Longitude,
        }),
      })
        .then((response) => response.json())
        .then((data) => {
          if (data) {
            if (data.message) {
              // validation succeed
              responseDiv.innerHTML = data.message;
              form.reset();
              getStations();
            } else {
              const keys = Object.keys(data);
              for (const key of keys) {
                if (data[key] != null) {
                  let responseHTML = "";
                  for (let i = 0; i < data[key].length; i++) {
                    responseHTML += `${data[key][i]}<br>`;
                  }
                  responseDiv.innerHTML = responseHTML;
                  break;
                }
              }
            }
          }
        })
        .catch((error) => {
          responseDiv.textContent =
            "Sorry! An error occurred while adding the new station.";
          console.error(error);
        });
    } else {
      // validation failed
      result.errors.forEach(
        errorArray => {
          if (errorArray.length > 1) {
            errorArray.forEach(
              error => responseDiv.innerHTML += `${error}<br>`
            );
          } else {
            responseDiv.innerHTML += `${errorArray}<br>`;
          }
        }
      );
    }
  });
};

const validateNewStation = async (newStation) => {
  const addressesValidationResult = await validateAddresses(
    newStation.Adress,
    newStation.Osoite
  );
  const namesValidationResult = validateNames(
    newStation.Nimi,
    newStation.Namn,
    newStation.Name
  );
  const operatorValidationResult = validateOperator(newStation.Operaattor);
  const capacityValidationResult = validateCapacity(newStation.Kapasiteet);
  const validateJourneysFromCount = validateJourneysCount(newStation.JourneysFrom);
  const validateJourneysToCount = validateJourneysCount(newStation.JourneysTo);

  if (
    addressesValidationResult.isValid &&
    namesValidationResult.isValid &&
    operatorValidationResult.isValid &&
    capacityValidationResult.isValid &&
    validateJourneysFromCount.isValid &&
    validateJourneysToCount.isValid
  ) {
    return {
      isValid: true,
      Kaupunki: addressesValidationResult.Kaupunki,
      Stad: addressesValidationResult.Stad,
      Latitude: addressesValidationResult.Latitude,
      Longitude: addressesValidationResult.Longitude,
    };
  } else {
    const errors = [];

    const validationResults = [
      addressesValidationResult,
      namesValidationResult,
      operatorValidationResult,
      capacityValidationResult,
      validateJourneysFromCount,
      validateJourneysToCount,
    ];

    validationResults.forEach((result) => {
      if (!result.isValid) {
        errors.push(result.error);
      }
    });

    return { isValid: false, errors };
  }
}

const validateAddresses = async (adress, osoite) => {
  if (adress === osoite) {
    return {
      isValid: false,
      error: ["Addresses in Finnish and Swedish must be different."],
    };
  }

  const addresses = [osoite, adress];
  const apiKey = "INSERT_GOOGLE_CLOUD_API_KEY";
  const apiUrl = "https://maps.googleapis.com/maps/api/geocode/json";

  const addressesPromises = addresses.map((address) => {
    const url = `${apiUrl}?address=${encodeURIComponent(
      address+", Helsinki, Finland"
    )}&key=${apiKey}&components=administrative_area:"Helsinki"&language=sv`;
    return fetch(url)
      .then((response) => response.json())
      .then((data) => {
        let latitude = undefined;
        let longitude = undefined;
        console.log(data.results)
        if (!(data.results[0].formatted_address == "Helsingfors, Finland")) {
          const addressComponents = data.results[0].address_components;
          const city = addressComponents.find((component) =>
            component.types.includes("locality")
          );
          if (
            city &&
            (city.long_name === "Helsingfors" || city.long_name === "Esbo")
          ) {
            latitude = data.results[0].geometry.location.lat;
            longitude = data.results[0].geometry.location.lng;
            return {
              isValid: true,
              Kaupunki: city.long_name === "Helsingfors" ? "Helsinki" : "Espo",
              Stad: city.long_name,
              Latitude: latitude,
              Longitude: longitude,
            };
          } else {
            return {
              isValid: false,
              error: [`The address ${address} is not in Helsinki or Espoo.`],
            };
          }
        } else {
          return {
            isValid: false,
            error: [`The address ${address} is not valid.`],
          };
        }
      })
      .catch((error) => {
        return {
          isValid: false,
          error: [`An error occurred while validating the address ${address}: ${error}`],
        };
      });
  });

  const promisesResults = await Promise.all(addressesPromises);

  const arePromisesResultsValid = promisesResults.every(
    (promiseResult) => promiseResult.isValid
  );

  if (arePromisesResultsValid) {
    const [address1, address2] = promisesResults;
    if (
      address1.Kaupunki === address2.Kaupunki &&
      address1.Stad === address2.Stad &&
      address1.Latitude === address2.Latitude &&
      address1.Longitude === address2.Longitude
    ) {
      // The addresses are valid and the same in different languages
      return address1;
    } else {
      return {
        isValid: false,
        error: ["Addresses are valid but not the same"],
      };
    }
  } else {
    const addressesErrors = [];
    for (const promiseResult of promisesResults) {
      if (!promiseResult.isValid) {
        addressesErrors.push(promiseResult.error);
      }
    }
    return { isValid: false, error: addressesErrors };
  }
};

const validateNames = (nimi, namn, name) => {
  const nameRegex = /^[a-zA-Z\s]*$/;
  const isNimiValid = nimi && nimi.trim().length > 0 && nameRegex.test(nimi);
  const isNamnValid = namn && namn.trim().length > 0 && nameRegex.test(namn);
  const isNameValid = name && name.trim().length > 0 && nameRegex.test(name);
  const areNamesDifferent = nimi.trim() !== namn.trim();

  if (isNimiValid && isNamnValid && isNameValid && areNamesDifferent) {
    return { isValid: true };
  } else {
    const namesErrors = [];
    if (!isNimiValid) {
      namesErrors.push("Finnish name is not valid.");
    }
    if (!isNamnValid) {
      namesErrors.push("Swedish name is not valid.");
    }
    if (!isNameValid) {
      namesErrors.push("English name is not valid.");
    }
    if (!areNamesDifferent) {
      namesErrors.push("Names in Finnish and Swedish must be different.");
    }
    return { isValid: false, error: namesErrors };
  }
};

const validateOperator = (operator) => {

  if (operator === ""){
    return { isValid: true }
  }
  
  const regex = /^[a-zA-Z\s]+$/;
  const isValid = regex.test(operator);

  if (!isValid) {
    return { isValid: false, error: ["Operator must include only letters."] };
  }

  return { isValid };
};

const validateCapacity = (capacity) => {
  if (isNaN(capacity) || capacity < 1 || capacity > 44) {
    return {
      isValid: false,
      error: ["Capacity must be a number greater than or equal to 1."],
    };
  }
  return { isValid: true };
};

const validateJourneysCount = (journeysCount) => {
  if (isNaN(journeysCount) || journeysCount < 0) {
    return {
      isValid: false,
      error: ["The counts of journeys must be a number greater than or equal to 0."],
    };
  }
  return { isValid: true };
};