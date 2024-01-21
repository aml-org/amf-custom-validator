/// <reference types="cypress" />

describe('double', () => {
  it('should render report', () => {
    cy.visit('./test/browser/cypress/html/double.html')
    cy.get('#report').should('include.text', "\"conforms\": false")
  })
})
