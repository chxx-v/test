test one md file contains two gherkin snippit
## First gherkin snippit
```gherkin
Feature: earn godogs md
  In order to live
  As a godog lover
  I need to be able to earn godogs

  Scenario: Have 12, earn 6 more md
    Given there are 12 godogs
    When I earn 6
    Then there should be 18 total
```
## Second gherkin snippit
```gherkin
Feature: eat godogs md2
  In order to be happy
  As a hungry gopher
  I need to be able to eat godogs

  Scenario: Eat 8 out of 12 md
    Given there are 12 godogs
    When I eat 8
    Then there should be 4 remaining
```