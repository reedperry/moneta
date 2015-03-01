controllers.controller('TrackingController', ['$location', 'Tracker', function TrackingControllerFactory($location, Tracker) {

    var track = this;
    var CREDIT = 'credit',
        EXPENSE = 'expense';

    track.getExpenses = Tracker.getExpenses;
    track.getTotalExpenses = Tracker.totalExpenses;

    track.addExpense = function() {
        var expense = {
            amount: track.amount,
            comment: track.comment,
            date: track.date,
            type: EXPENSE
        };
        Tracker.addExpense(expense);
        track.clearEntry();
    };

    track.getCredits = Tracker.getCredits;
    track.getTotalCredits = Tracker.totalCredits;

    track.addCredit = function() {
        var credit = {
            amount: track.amount,
            comment: track.comment,
            date: track.date,
            type: CREDIT
        };
        Tracker.addCredit(credit);
        track.clearEntry();
    };
    
    track.clearEntry = function() {
        track.amount = null;
        track.comment = null;
        track.date = new Date();
    };

    track.clearEntry();
}]);
