describe('Stations test', () => {
    it('Clicks on 10 stations', () => {
      cy.visit('http://localhost:8080')
  
      cy.get('#stations-list > div').each(($div, index) => {
        if (index >= 10) return false; // exit the loop after the 10th iteration
  
        // Get the station string from the div's inner HTML
        const stationName = $div.html().split('<br>').join(', ').toLowerCase().trim();
  
        // Click on the station div and wait for the station details to be loaded
        cy.wrap($div).click().then(() => {
          cy.wait(2000);
  
          // Get the station details text and remove newline characters
          cy.get('#station-details-text h2').invoke('text').then((text) => {
            const lowercaseText = text.toLowerCase().replace(/\n/g, '');
            const lowercaseTextTrimmed = lowercaseText.replace(/\s/g, '');
            const stationNameTrimmed = stationName.replace(/\s/g, '');
  
            // Assert that the station name is in the station details text
            expect(lowercaseTextTrimmed).to.equal(stationNameTrimmed);
          });
        });
      });    
    })
  });
  