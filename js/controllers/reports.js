controllers.controller('ReportsController', ['$location', 'Tracker', function ReportsControllerFactory($location, Tracker) {

    var reports = this;

    reports.getTotalExpenses = Tracker.totalExpenses;
    reports.getTotalIncomes = Tracker.totalIncomes;

}]);
