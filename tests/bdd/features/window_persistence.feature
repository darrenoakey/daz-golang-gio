Feature: Window Position Persistence

  Scenario: App launches with default window
    Given the app is running
    When I browse to the example window
    Then the window should be visible
    And the window should have a non-zero size

  Scenario: Fresh app creates no state file
    Given no saved state exists for "bddtest"
    When I launch a persist window named "bddtest"
    Then the window should be created successfully

  Scenario: Window state is saved on close
    Given a persist window "bddtest" is open
    When I close the window
    Then the saved state file for "bddtest" should exist
    And the saved width should be greater than zero
    And the saved height should be greater than zero

  Scenario: Window restores saved size on relaunch
    Given saved state for "bddtest" has width 900 and height 700
    When I launch a persist window named "bddtest"
    Then the window size should be approximately 900x700

  Scenario: State file uses correct location
    Given a persist window "bddtest" is open
    When I close the window
    Then the state file should be at "~/.config/daz-golang-gio/bddtest.json"

  Scenario: Multiple apps have independent state
    Given saved state for "app1" has width 800 and height 600
    And saved state for "app2" has width 1200 and height 900
    When I launch a persist window named "app1"
    Then the window size should be approximately 800x600

  Scenario: Window position is saved on macOS
    Given a persist window "bddtest" is open on macOS
    When I move the window to position 200 300
    And I close the window
    Then the saved x should be approximately 200
    And the saved y should be approximately 300

  Scenario: Window restores saved position on macOS
    Given saved state for "bddtest" has position 200 300 and size 900 700
    When I launch a persist window named "bddtest"
    Then the window position should be approximately 200 300
