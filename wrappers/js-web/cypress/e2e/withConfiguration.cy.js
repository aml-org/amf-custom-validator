/// <reference types="cypress" />

describe('withConfiguration', () => {
  it('should render report', () => {
    cy.visit('cypress/html/withConfiguration.html')
    cy.get('#report').should('have.text', 'does not have date')
  })
})
