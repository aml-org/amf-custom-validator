/// <reference types="cypress" />

describe('double', () => {
  it('should render report', () => {
    cy.visit('cypress/html/double.html')
    cy.get('#report').should('include.text', "\"conforms\": false")
  })
})
