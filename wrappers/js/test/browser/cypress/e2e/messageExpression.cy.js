/// <reference types="cypress" />

describe('messageExpression', () => {
  it('should render report with custom message', () => {
    cy.visit('./test/browser/cypress/html/messageExpression.html')
    cy.get('#report').should('include.text', "Movie 'Disaster Movie' has a rating of 1.9 but it does not have at least 10 reviews (actual reviews: 5) to support that rating")
  })
})
