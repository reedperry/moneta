controllers.controller('AppController', ['$location', App]);

function App($location) {
    var app = this;
    app.$location = $location;
}

App.prototype.goHome = function() {
    this.$location.path('/');
}

App.prototype.goToTracking = function() {
    this.$location.path('/track');
}

App.prototype.goToReports = function() {
    this.$location.path('/reports');
}

App.prototype.goToPreferences = function() {
    this.$location.path('/prefs');
}
