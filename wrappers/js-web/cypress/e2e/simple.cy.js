/// <reference types="cypress" />

describe('simple', () => {
  it('should render report', () => {
    cy.visit('cypress/html/simple.html')
    cy.get('#report').should('include.text', "\"conforms\": false")
  })
})
