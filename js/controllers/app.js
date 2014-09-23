controllers.controller('AppController', ['$location', App]);

function App($location) {
    var app = this;
    app.$location = $location;

    this.goHome = function() {
        this.$location.path('/');
    }

    this.goToTracking = function() {
        this.$location.path('/track');
    }

    this.goToReports = function() {
        this.$location.path('/reports');
    }

    this.goToPreferences = function() {
        this.$location.path('/prefs');
    }
}

