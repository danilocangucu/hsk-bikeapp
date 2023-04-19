describe("Stations details display test", () => {
  it("Clicks on a station and verifies that its details match the API data", () => {
    cy.visit("http://localhost:8080");

    cy.get("#stations-list > div").each(($div, index) => {
      if (index >= 10) return false; // exit the loop after the 10th iteration

      const stationId = $div.attr("id");

      // Click on the station div and wait for the station details to be loaded
      cy.wrap($div)
        .click()
        .then(() => {
          cy.wait(1000);

          cy.get("#station-details-text")
            .invoke("text")
            .then((detailsText) => {
              // Make a request to the detail page and assert that its data is the single station API
              cy.request(`http://localhost:8080/stations?id=${stationId}`)
                .its("body")
                .then((stationData) => {
                  expect(detailsText).to.contain(
                    stationData[0].Nimi
                  );
                  expect(detailsText).to.contain(
                    stationData[0].Namn
                  );
                  expect(detailsText).to.contain(
                    stationData[0].Osoite
                  );
                  expect(detailsText).to.contain(
                    stationData[0].Adress
                  );
                  expect(detailsText).to.contain(
                    stationData[0].JourneysFrom
                  );
                  expect(detailsText).to.contain(
                    stationData[0].JourneysTo
                  );

                });
            });
        });
    });
  });
});

describe("Stations list names test", () => {
  it("Scrolls and verifies 20 more station names per scroll", () => {
    cy.visit("http://localhost:8080");

    // Define the number of times to scroll
    const numScrolls = 5;

    // Get the initial number of stations displayed
    cy.get("#stations-list > div").then(($stations) => {
      let count = $stations.length;

      // Scroll down and verify the number of stations after each scroll
      for (let i = 0; i < numScrolls; i++) {
        cy.get("#stations-list").scrollTo("bottom");
        cy.wait(2000); // adjust the wait time as needed
        cy.get("#stations-list > div").should("have.length", count + 20);
        count += 20;
      }
    });
  });
});


