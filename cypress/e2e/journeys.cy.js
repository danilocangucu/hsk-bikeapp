describe("Journeys details should match the API data", () => {
  it("Checks that each journey's details are correct", function () {
    cy.visit("http://localhost:8080");

    cy.request(`http://localhost:8080/journeys`).its("body").as("journeysData");

    cy.get("table tr").each(($row, index) => {
      // Skip the first row, which contains the column headers
      if (index === 0) return;

      // Find the <td> elements within the row
      const $tds = $row.find("td");

      // Extract the station names and timestamps from the <td> elements
      const departureTimestamp = $tds.eq(0).text();
      const departureStationName = $tds.eq(1).text();
      const returnStationName = $tds.eq(2).text();
      const distanceInKm = $tds.eq(3).text().replace(/\.0+/, "");

      // Convert the date in the table to the same format as the API response
      const dateParts = departureTimestamp.split(/\s+/)[0].split(".");
      const timeParts = departureTimestamp.split(/\s+/)[1].split(".");
      const isoDateString = `${dateParts[2]}-${dateParts[1]}-${dateParts[0]}T${timeParts[0]}:${timeParts[1]}:${timeParts[2]}`;

      // Use the station names and timestamps to verify the row details against the API data
      expect(isoDateString).to.equal(this.journeysData[index - 1].Departure);
      expect(departureStationName).to.equal(
        this.journeysData[index - 1].DepartureStationName
      );
      expect(returnStationName).to.equal(
        this.journeysData[index - 1].ReturnStationName
      );
      const apiDistanceInKm =
        Math.round(this.journeysData[index - 1].CoveredDistanceM / 100) / 10;
      expect(distanceInKm).to.equal(apiDistanceInKm + "km");
    });
  });
});
