/// <reference types="cypress" />

describe('normalize', () => {
  it('should render normalized input', () => {
    cy.visit('cypress/html/normalize.html')
    cy.get('#normalize').should('include.text', "@ids")
  })
})
