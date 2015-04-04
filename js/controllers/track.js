controllers.controller('TrackingController', ['$location', 'DB', 'Tracker', function TrackingControllerFactory($location, DB, Tracker) {

    var track = this;
    var INCOME = 'income',
        EXPENSE = 'expense';

    track.getExpenses = Tracker.getExpenses;
    track.getTotalExpenses = Tracker.totalExpenses;

    track.addExpense = function() {
        var expense = {
            amount: track.amount,
            comment: track.comment,
            date: track.date,
            kind: EXPENSE
        };
        DB.storeExpense(expense);
        Tracker.addExpense(expense);
        track.clearEntry();
    };

    track.getIncomes = Tracker.getIncomes;
    track.getTotalIncomes = Tracker.totalIncomes;

    track.addIncome = function() {
        var income = {
            amount: track.amount,
            comment: track.comment,
            date: track.date,
            kind: INCOME
        };
        DB.storeIncome(income);
        Tracker.addIncome(income);
        track.clearEntry();
    };

    track.clearEntry = function() {
        track.amount = null;
        track.comment = null;
        track.date = new Date();
    };

    track.clearEntry();
}]);
