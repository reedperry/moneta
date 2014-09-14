controllers.controller('AppController', ['$location', function AppControllerFactory($location) {

    var app = this;

    app.goHome = function() {
        $location.path('/');
    }

    app.goToTracking = function() {
        $location.path('/tracking');
    }

    app.goToReports = function() {
        $location.path('/reports');
    }

    app.goToPreferences = function() {
        $location.path('/prefs');
    }

}]);
