/// <reference types="cypress" />

describe('simple', () => {
  it('should render report', () => {
    cy.visit('./test/browser/cypress/html/simple.html')
    cy.get('#report').should('include.text', "\"conforms\": false")
  })
})
