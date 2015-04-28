controllers.controller('PrefsController', ['DB', PrefsControllerFactory]);

function PrefsControllerFactory(DB) {

  var prefs = this;

  prefs.startAddingCategory = startAddingCategory;
  prefs.endAddingCategory = endAddingCategory;

  prefs.daysOfWeek = [
    'Sunday',
    'Monday',
    'Tuesday',
    'Wednesday',
    'Thursday',
    'Friday',
    'Saturday'
  ];

  function startAddingCategory() {
    prefs.addingCategory = true;
  }

  function endAddingCategory() {
    prefs.addingCategory = false;
  }

}
